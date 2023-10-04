import bisect
import math
import os
import re
import sys
import warnings
from typing import Sequence

import lmdb
import numpy as np
import six
import tensorflow as tf
from natsort import natsorted
from PIL import Image
from tensorflow import keras
from tensorflow.keras import preprocessing

from preprocess_image import all_preprocessing


# taken from python 3.5 docs
def _accumulate(iterable, fn=lambda x, y: x + y):
    "Return running totals"
    # _accumulate([1,2,3,4,5]) --> 1 3 6 10 15
    # _accumulate([1,2,3,4,5], operator.mul) --> 1 2 6 24 120
    it = iter(iterable)
    try:
        total = next(it)
    except StopIteration:
        return
    yield total
    for element in it:
        total = fn(total, element)
        yield total


def hierarchical_dataset(root, opt, select_data="/"):
    """select_data='/' contains all sub-directory of root directory"""
    dataset_list = []
    dataset_log = f"dataset_root:    {root}\t dataset: {select_data[0]}"
    print(dataset_log)
    dataset_log += "\\n"
    for dirpath, dirnames, filenames in os.walk(root + "/"):
        if not dirnames:
            select_flag = False
            for selected_d in select_data:
                if selected_d in dirpath:
                    select_flag = True
                    break

            if select_flag:
                dataset = LmdbDataset(dirpath, opt)
                sub_dataset_log = f"sub-directory:\t/{os.path.relpath(dirpath, root)}\t num samples: {len(dataset)}"
                print(sub_dataset_log)
                dataset_log += f"{sub_dataset_log}\n"
                dataset_list.append(dataset)

    concatenated_dataset = ConcatDataset(dataset_list)
    return concatenated_dataset, dataset_log


def tensorflow_dataloader(
    dataset,
    batch_size=2,
    shuffle=False,
    num_workers=0,
    collate_fn=None,
    prefetch_factor=2,
    imgH=32,
    imgW=100
):
    data = tf.data.Dataset.from_generator(
        dataset,
        output_signature=(
            tf.TensorSpec(shape=(imgH, imgW, 1), dtype=tf.float64),
            tf.TensorSpec(shape=(1), dtype=tf.string),
        ),
    )
    data = data.batch(batch_size)
    data = data.prefetch(prefetch_factor)
    return data


def tensor2im(image_tensor, imtype=np.uint8):
    image_numpy = tf.cast(image_tensor, dtype=tf.float32).numpy()
    if image_numpy.shape[0] == 1:
        image_numpy = np.tile(image_numpy, (3, 1, 1))
    image_numpy = (np.transpose(image_numpy, (1, 2, 0)) + 1) / 2.0 * 255.0
    return image_numpy.astype(imtype)


def save_image(image_numpy, image_path):
    image_pil = Image.fromarray(image_numpy)
    image_pil.save(image_path)


# need to be fix
class ApplyCollate(keras.utils.Sequence):
    def __init__(self, dataset, collate_fn=None):
        self.dataset = dataset
        self.collate_fn = collate_fn
        self.indexs = tf.random.shuffle(tf.range(len(self.dataset)))

    def __getitem__(self, idx):
        # print(type(self.dataset[idx][1]))
        return self.collate_fn([self.dataset[idx]])

    def __len__(self):
        return len(self.dataset)

    def on_epoch_end(self):
        self.indexs = tf.random.shuffle(self.indexs)

    def __call__(self):
        self.indexs = tf.random.shuffle(self.indexs)
        yield self.__getitem__(self.indexs[0])


# rawan error
class Subset(keras.utils.Sequence):
    def __init__(self, dataset, indices: Sequence, batch_size=2, collate_fn=None):
        self.dataset = dataset
        self.indices = indices
        self.batch_size = batch_size
        self.collate_fn = collate_fn
        self.indexs = tf.random.shuffle(tf.range(len(self)))

    def __getitem__(self, idx):
        if isinstance(idx, list):
            return self.collate_fn([self.dataset[[self.indices[i] for i in idx]]])
        else:
            return self.collate_fn([self.dataset[self.indices[idx]]])

    def __len__(self):
        return len(self.indices)

    def on_epoch_end(self):
        self.indexs = tf.random.shuffle(self.indexs)

    def __call__(self):
        for i in range(len(self)):
            self.indexs = tf.random.shuffle(self.indexs)
            yield self.__getitem__(self.indexs[i])


# rawan error
class ConcatDataset(keras.utils.Sequence):
    @staticmethod
    def cumsum(sequence):
        r, s = [], 0
        for e in sequence:
            l = len(e)
            r.append(l + s)
            s += l
        return r

    def __init__(self, datasets: keras.utils.Sequence):
        super(ConcatDataset, self).__init__()
        self.datasets = list(datasets)
        assert len(self.datasets) > 0, "datasets should not be an empty iterable"
        self.cumulative_sizes = self.cumsum(self.datasets)

    def __len__(self):
        return self.cumulative_sizes[-1]

    def __getitem__(self, index):
        if index < 0:
            if -index > len(self):
                raise ValueError(
                    "absolute value of index should not exceed dataset length"
                )

            index = len(self) + index
        dataset_index = bisect.bisect_right(self.cumulative_sizes, index)
        if dataset_index == 0:
            sample_index = index
        else:
            sample_index = index - self.cumulative_sizes[dataset_index - 1]
        return self.datasets[dataset_index][sample_index]

    @property
    def cummulative_size(self):
        warnings.warn(
            "cummulative_sizes attribute is renamed to " "cumulative_sizes",
            DeprecationWarning,
            stacklevel=2,
        )
        return self.cumulative_sizes

    def __call__(self):
        for i in range(len(self)):
            yield self.__getitem__(i)


class LmdbDataset(keras.utils.Sequence):
    def __init__(self, root, opt):
        self.root = root
        self.opt = opt
        self.env = lmdb.open(
            root,
            max_readers=32,
            readonly=True,
            lock=False,
            readahead=False,
            meminit=False,
        )
        if not self.env:
            print("cannot create lmdb from %s" % (root))
            sys.exit(0)

        with self.env.begin(write=False) as txn:
            nSamples = int(txn.get("num-samples".encode()))
            self.nSamples = nSamples

            if self.opt.data_filtering_off:
                # for fast check or benchmark evaluation with no filtering
                self.filtered_index_list = [index + 1 for index in range(self.nSamples)]
            else:
                """Filtering part
                If you want to evaluate IC15-2077 & CUTE datasets which have special character labels,
                use --data_filtering_off and only evaluate on alphabets and digits.
                see https://github.com/clovaai/deep-text-recognition-benchmark/blob/6593928855fb7abb999a99f428b3e4477d4ae356/dataset.py#L190-L192
                And if you want to evaluate them with the model trained with --sensitive option,
                use --sensitive and --data_filtering_off,
                see https://github.com/clovaai/deep-text-recognition-benchmark/blob/dff844874dbe9e0ec8c5a52a7bd08c7f20afe704/test.py#L137-L144
                """
                self.filtered_index_list = []
                for index in range(self.nSamples):
                    index += 1  # lmdb starts with 1
                    label_key = "label-%09d".encode() % index
                    label = txn.get(label_key).decode("utf-8")

                    if len(label) > self.opt.batch_max_length:
                        # print(f'The length of the label is longer than max_length: length
                        # {len(label)}, {label} in dataset {self.root}')
                        continue

                    # By default, images containing characters which are not in opt.character are filtered.
                    # You can add [UNK] token to `opt.character` in utils.py instead of this filtering.
                    out_of_char = f"[^{self.opt.character}]"
                    if re.search(out_of_char, label.lower()):
                        continue

                    self.filtered_index_list.append(index)

                self.nSamples = len(self.filtered_index_list)

    def __len__(self):
        return self.nSamples

    def __getitem__(self, index):
        assert index <= len(self), "index range error"
        index = self.filtered_index_list[index]

        with self.env.begin(write=False) as txn:
            label_key = "label-%09d".encode() % index
            label = txn.get(label_key).decode("utf-8")
            img_key = "image-%09d".encode() % index
            imgbuf = txn.get(img_key)

            buf = six.BytesIO()
            buf.write(imgbuf)
            buf.seek(0)
            try:
                if self.opt.rgb:
                    img = Image.open(buf).convert("RGB")  # for color image
                else:
                    img = Image.open(buf).convert("L")

            except IOError:
                print(f"Corrupted image for {index}")
                # make dummy image and dummy label for corrupted image.
                if self.opt.rgb:
                    img = Image.new("RGB", (self.opt.imgW, self.opt.imgH))
                else:
                    img = Image.new("L", (self.opt.imgW, self.opt.imgH))
                label = "[dummy_label]"

            if not self.opt.sensitive:
                label = label.lower()

            # We only train and evaluate on alphanumerics (or pre-defined character set in train.py)
            out_of_char = f"[^{self.opt.character}]"
            label = re.sub(out_of_char, "", label)

        return (img, label)

    def __call__(self):
        for i in range(len(self)):
            yield self.__getitem__(i)


class SingleDataset(keras.utils.Sequence):
    def __init__(self, image, left, top, right, bottom, collate_fn):
        self.image = image
        self.left = left
        self.top = top
        self.right = right
        self.bottom = bottom
        self.collate_fn = collate_fn

    def __len__(self):
        return 1

    def __getitem__(self, index: int):
        image_preprocessed = all_preprocessing(
            self.image, self.left, self.top, self.right, self.bottom
        )
        image_preprocessed = self.collate_fn([(image_preprocessed, "Prediction")])

        return image_preprocessed

    def __call__(self):
        for i in range(len(self)):
            yield self.__getitem__(i)


class RawDataset(keras.utils.Sequence):
    def __init__(self, root, opt):
        self.opt = opt
        self.image_path_list = []

        for dirpath, dirnames, filenames in os.walk(root):
            for name in filenames:
                _, ext = os.path.splitext(name)
                ext = ext.lower()
                if ext == ".jpg" or ext == ".jpeg" or ext == ".png":
                    self.image_path_list.append(os.path.join(dirpath, name))

        self.image_path_list = natsorted(self.image_path_list)
        self.nSamples = len(self.image_path_list)

    def __len__(self):
        return self.nSamples

    def __getitem__(self, index):
        try:
            if self.opt.rgb:
                img = Image.open(self.image_path_list[index]).convert(
                    "RGB"
                )  # for color image
            else:
                img = Image.open(self.image_path_list[index]).convert("L")

        except IOError:
            print(f"Corrupted image for {index}")
            # make dummy image and dummy label for corrupted image.
            if self.opt.rgb:
                img = Image.new("RGB", (self.opt.imgW, self.opt.imgH))
            else:
                img = Image.new("L", (self.opt.imgW, self.opt.imgH))

        return (img, self.image_path_list[index])

    def __call__(self):
        for i in range(len(self)):
            yield self.__getitem__(i)


class ResizeNormalize(object):
    def __init__(self, size, interpolation=Image.BICUBIC) -> None:
        self.size = size
        self.interpolation = interpolation
        self.toTensor = preprocessing.image.img_to_array

    def __call__(self, img) -> tf.Tensor:
        img = img.resize(self.size, self.interpolation)
        img = self.toTensor(img)
        img = tf.math.divide(img, 255.0)
        img = tf.math.multiply(img, 2)
        img = tf.math.subtract(img, 1)
        return img


class NormalizePAD(object):
    def __init__(self, max_size, PAD_type="right") -> None:
        self.toTensor = preprocessing.image.img_to_array
        self.max_size = max_size
        self.max_width_half = tf.floor(max_size[2] / 2)
        self.PAD_type = PAD_type

    def __call__(self, image: Image) -> tf.Tensor:
        img = self.toTensor(image)
        c, h, w = img.shape
        Pad_img = tf.zeros(shape=self.max_size)
        Pad_img[:, :, :w] = img
        if self.max_size[2] != w:
            Pad_img[:, :, w:] = tf.broadcast_to(
                tf.expand_dims(img, axis=2), shape=(c, h, self.max_size[2] - w)
            )

        return Pad_img


class AlignCollate(object):
    def __init__(self, imgH=32, imgW=100, keep_ratio_with_pad=False) -> None:
        self.imgH = imgH
        self.imgW = imgW
        self.keep_ratio_with_pad = keep_ratio_with_pad

    def __call__(self, batch):
        batch = filter(lambda x: x is not None, batch)
        images, labels = zip(*batch)
        if self.keep_ratio_with_pad:
            print("masuk padding")
            resized_max_w = self.imgW
            input_channel = 3 if images[0].mode == "RGB" else 1
            transform = NormalizePAD((input_channel, self.imgH, resized_max_w))

            resized_images = []
            for image in images:
                w, h = image.shape

                ratio = w / float(h)

                if math.ceil(self.imgH * ratio) > self.imgW:
                    resized_w = self.imgW
                else:
                    resized_w = math.ceil(self.imgH * ratio)

                resized_image = image.resize((resized_w, self.imgH), Image.BICUBIC)
                resized_images.append(transform(resized_image))
                # resized_image.save('./image_test/%d_test.jpg' % w)

            image_tensors = tf.concat(resized_images, 0)

        else:
            transform = ResizeNormalize((self.imgW, self.imgH))
            image_tensors = [transform(image) for image in images]
            image_tensors = tf.concat(image_tensors, 0)

        return image_tensors, labels


class Batch_Balanced_Dataset(object):
    def __init__(self, opt) -> None:
        """
        Modulate the data ratio in the batch.
        For example, when select_data is "MJ-ST" and batch_ratio is "0.5-0.5",
        the 50% of the batch is filled with MJ and the other 50% of the batch is filled with ST.
        """
        log = open(f"./saved_models/{opt.exp_name}/log_dataset.txt", "a")
        dashed_line = "-" * 80
        print(dashed_line)
        log.write(dashed_line + "\n")
        print(
            f"dataset_root: {opt.train_data}\nopt.select_data: {opt.select_data}\nopt.batch_ratio: {opt.batch_ratio}"
        )
        log.write(
            f"dataset_root: {opt.train_data}\nopt.select_data: {opt.select_data}\nopt.batch_ratio: {opt.batch_ratio}\n"
        )
        assert len(opt.select_data) == len(opt.batch_ratio)

        _AlignCollate = AlignCollate(
            imgH=opt.imgH, imgW=opt.imgW, keep_ratio_with_pad=opt.PAD
        )
        self.data_loader_list = []
        self.dataloader_iter_list = []
        batch_size_list = []
        Total_batch_size = 0
        for selected_d, batch_ratio_d in zip(opt.select_data, opt.batch_ratio):
            _batch_size = max(round(opt.batch_size * float(batch_ratio_d)), 1)
            print(dashed_line)
            log.write(dashed_line + "\\n")
            _dataset, _dataset_log = hierarchical_dataset(
                root=opt.train_data, opt=opt, select_data=[selected_d]
            )
            total_number_dataset = len(_dataset)
            log.write(_dataset_log)

            """
            The total number of data can be modified with opt.total_data_usage_ratio.
            ex) opt.total_data_usage_ratio = 1 indicates 100% usage, and 0.2 indicates 20% usage.
            See 4.2 section in our paper.
            """
            number_dataset = int(
                total_number_dataset * float(opt.total_data_usage_ratio)
            )
            dataset_split = [number_dataset, total_number_dataset - number_dataset]
            indices = range(total_number_dataset)
            _dataset, _ = [
                Subset(
                    _dataset,
                    indices[offset - length : offset],
                    collate_fn=_AlignCollate,
                )
                for offset, length in zip(_accumulate(dataset_split), dataset_split)
            ]
            selected_d_log = f"num total samples of {selected_d}: {total_number_dataset} x {opt.total_data_usage_ratio} (total_data_usage_ratio) = {len(_dataset)}\n"
            selected_d_log += f"num samples of {selected_d} per batch: {opt.batch_size} x {float(batch_ratio_d)} (batch_ratio) = {_batch_size}"
            print(selected_d_log)
            log.write(selected_d_log + "\\n")
            batch_size_list.append(str(_batch_size))
            Total_batch_size += _batch_size

            _data_loader = tensorflow_dataloader(
                _dataset,
                batch_size=_batch_size,
                shuffle=True,
                collate_fn=_AlignCollate,
                num_workers=int(opt.workers),
                imgH=opt.imgH,
                imgW=opt.imgW
            )
            self.data_loader_list.append(_data_loader)
            self.dataloader_iter_list.append(iter(_data_loader))

        Total_batch_size_log = f"{dashed_line}\n"
        batch_size_sum = "+".join(batch_size_list)
        Total_batch_size_log += (
            f"Total_batch_size: {batch_size_sum} = {Total_batch_size}\n"
        )
        Total_batch_size_log += f"{dashed_line}"
        opt.batch_size = Total_batch_size

        print(Total_batch_size_log)
        log.write(Total_batch_size_log + "\n")
        log.close()

    def get_batch(self):
        balanced_batch_images = []
        balanced_batch_texts = []

        for i, data_loader_iter in enumerate(self.data_loader_list):
            try:
                image, text = next(data_loader_iter.as_numpy_iterator())
                balanced_batch_images.append(image)
                balanced_batch_texts.append(text[0])
            except StopIteration:
                self.dataloader_iter_list[i] = iter(self.data_loader_list[i])
                image, text = self.dataloader_iter_list[i]
                balanced_batch_images.append(image)
                balanced_batch_texts.append(text[0])
            except ValueError:
                pass

        balanced_batch_images = tf.concat(balanced_batch_images, axis=0)

        return balanced_batch_images, balanced_batch_texts
