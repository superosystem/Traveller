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
from PIL import Image
from natsort import natsorted
from tensorflow import keras
from tensorflow.keras import preprocessing

from command.preprocess_image import all_preprocessing


# taken from python 3.5 docs
def _accumulate(iterable, fn=lambda x, y: x + y):
    """
    Return an iterator that yields running totals.
    Args:
        iterable (iterable): An iterable sequence of numbers.
        fn (function): A function used for accumulation. Default is addition (lambda x, y: x + y).
    Returns:
        iterator: An iterator that yields running totals.
    """
    # Example usage:
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
    """
    Create a hierarchical dataset from the specified root directory and selected sub-directories.
    Args:
        root (str): Root directory containing sub-directories.
        opt (object): Configuration options.
        select_data (str): Sub-directory selector. Default is '/' which includes all sub-directories.
    Returns:
        dataset (torch.utils.data.ConcatDataset): Concatenated dataset containing samples from selected sub-directories.
        dataset_log (str): Log containing information about the selected sub-directories and their sample counts.
    """
    dataset_list = []  # List to store individual datasets from selected sub-directories
    dataset_log = f"dataset_root:    {root}\t dataset: {select_data[0]}"
    print(dataset_log)
    dataset_log += "\\n"

    # Traverse through the directory hierarchy
    for dirpath, dirnames, filenames in os.walk(root + "/"):
        if not dirnames:  # If the current directory has no sub-directories
            select_flag = False
            for selected_d in select_data:
                if selected_d in dirpath:
                    select_flag = True
                    break

            if select_flag:  # If the directory matches the selection criteria
                dataset = LmdbDataset(dirpath, opt)  # Create a dataset from the directory
                sub_dataset_log = f"sub-directory:\t/{os.path.relpath(dirpath, root)}\t num samples: {len(dataset)}"
                print(sub_dataset_log)
                dataset_log += f"{sub_dataset_log}\n"
                dataset_list.append(dataset)  # Add the dataset to the list

    concatenated_dataset = ConcatDataset(dataset_list)  # Concatenate individual datasets
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
    """
    Create a TensorFlow dataset loader from the given dataset.
    Args:
        dataset (callable): A generator function or any callable that yields tuples (image, label).
        batch_size (int): Number of samples in each batch. Default is 2.
        shuffle (bool): Whether to shuffle the dataset. Default is False.
        num_workers (int): Number of parallel workers for data loading. Default is 0.
        collate_fn (callable): Function to collate samples into batches. Default is None.
        prefetch_factor (int): Number of batches to prefetch asynchronously. Default is 2.
        imgH (int): Height of the input images. Default is 32.
        imgW (int): Width of the input images. Default is 100.
    Returns:
        tf.data.Dataset: TensorFlow dataset containing batches of images and corresponding labels.
    """
    data = tf.data.Dataset.from_generator(
        dataset,
        output_signature=(
            tf.TensorSpec(shape=(imgH, imgW, 1), dtype=tf.float64),
            tf.TensorSpec(shape=(1), dtype=tf.string),
        ),
    )
    data = data.batch(batch_size)
    if shuffle:
        data = data.shuffle(buffer_size=1000)  # Buffer size can be adjusted according to dataset size
    data = data.prefetch(prefetch_factor)
    return data


def tensor2im(image_tensor, imtype=np.uint8):
    """
    Convert a TensorFlow tensor to a NumPy image array.
    Args:
        image_tensor (tf.Tensor): Input image tensor.
        imtype (numpy.dtype): Data type for the output NumPy array. Default is np.uint8.
    Returns:
        numpy.ndarray: NumPy array representing the image.
    """
    image_numpy = tf.cast(image_tensor, dtype=tf.float32).numpy()
    if image_numpy.shape[0] == 1:
        image_numpy = np.tile(image_numpy, (3, 1, 1))
    image_numpy = (np.transpose(image_numpy, (1, 2, 0)) + 1) / 2.0 * 255.0
    return image_numpy.astype(imtype)


def save_image(image_numpy, image_path):
    """
    Save a NumPy image array as an image file.
    Args:
        image_numpy (numpy.ndarray): NumPy image array.
        image_path (str): File path to save the image.
    """
    image_pil = Image.fromarray(image_numpy)
    image_pil.save(image_path)


class ApplyCollate(keras.utils.Sequence):
    """
    Custom data generator class for collating samples using a specified collate function.
    Args:
        dataset (list): List of samples.
        collate_fn (callable): Function to collate samples. Default is None.
    """

    def __init__(self, dataset, collate_fn=None):
        self.dataset = dataset
        self.collate_fn = collate_fn
        self.indexes = tf.random.shuffle(tf.range(len(self.dataset)))

    def __getitem__(self, idx):
        """
        Get a batch of collated samples.
        Args:
            idx (int): Index of the batch.
        Returns:
            Any: Collated batch of samples.
        """
        return self.collate_fn([self.dataset[idx]])

    def __len__(self):
        """
        Get the number of batches in the dataset.
        Returns:
            int: Number of batches.
        """
        return len(self.dataset)

    def on_epoch_end(self):
        """
        Shuffle the dataset at the end of each epoch.
        """
        self.indexes = tf.random.shuffle(self.indexes)

    def __call__(self):
        """
        Generate batches by shuffling and yielding collated samples.
        Yields:
            Any: Collated batch of samples.
        """
        self.indexes = tf.random.shuffle(self.indexes)
        yield self.__getitem__(self.indexes[0])


class Subset(keras.utils.Sequence):
    """
    Custom data generator class for creating a subset of a dataset based on specified indices.
    Args:
        dataset (list): List of samples.
        indices (Sequence): List or sequence of indices representing the subset.
        batch_size (int): Number of samples in each batch. Default is 2.
        collate_fn (callable): Function to collate samples. Default is None.
    """

    def __init__(self, dataset, indices: Sequence, batch_size=2, collate_fn=None):
        self.dataset = dataset
        self.indices = indices
        self.batch_size = batch_size
        self.collate_fn = collate_fn
        self.indexes = tf.random.shuffle(tf.range(len(self)))

    def __getitem__(self, idx):
        """
        Get a batch of collated samples based on the specified indices.
        Args:
            idx (int or list): Index of the batch (or list of indices).
        Returns:
            Any: Collated batch of samples.
        """
        if isinstance(idx, list):
            return self.collate_fn([self.dataset[[self.indices[i] for i in idx]]])
        else:
            return self.collate_fn([self.dataset[self.indices[idx]]])

    def __len__(self):
        """
        Get the number of batches in the subset.
        Returns:
            int: Number of batches.
        """
        return len(self.indices)

    def on_epoch_end(self):
        """
        Shuffle the dataset at the end of each epoch.
        """
        self.indexes = tf.random.shuffle(self.indexes)

    def __call__(self):
        """
        Generate batches by shuffling and yielding collated samples.
        Yields:
            Any: Collated batch of samples.
        """
        for i in range(len(self)):
            self.indexes = tf.random.shuffle(self.indexes)
            yield self.__getitem__(self.indexes[i])


class ConcatDataset(keras.utils.Sequence):
    """
    Concatenates multiple datasets into a single dataset.
    Args:
        datasets (keras.utils.Sequence): List of datasets to concatenate.
    """

    @staticmethod
    def cumsum(sequence):
        """
        Computes the cumulative sum of the lengths of a sequence.
        Args:
            sequence (iterable): Input sequence.
        Returns:
            list: Cumulative sum of the lengths of the sequence.
        """
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
        """
        Get the total number of samples in the concatenated dataset.
        Returns:
            int: Total number of samples.
        """
        return self.cumulative_sizes[-1]

    def __getitem__(self, index):
        """
        Get a sample from the concatenated dataset.
        Args:
            index (int): Index of the sample.

        Returns:
            Any: Sample from the concatenated dataset.
        """
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

    def __call__(self):
        """
        Generate samples from the concatenated dataset.
        Yields:
            Any: Sample from the concatenated dataset.
        """
        for i in range(len(self)):
            yield self.__getitem__(i)

    @property
    def cummulative_size(self):
        """
        Warning: This attribute is deprecated. Use cumulative_sizes instead.
        """
        warnings.warn(
            "cummulative_sizes attribute is renamed to " "cumulative_sizes",
            DeprecationWarning,
            stacklevel=2,
        )
        return self.cumulative_sizes


class LmdbDataset(keras.utils.Sequence):
    """
    Custom data generator class for loading data from an LMDB database.
    Args:
        root (str): Root directory of the LMDB database.
        opt (object): Options object containing dataset configuration.
    """

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
                self.filtered_index_list = [index + 1 for index in range(self.nSamples)]
            else:
                # Filtering part: Remove samples based on specified criteria
                self.filtered_index_list = []
                for index in range(self.nSamples):
                    index += 1  # LMDB starts with 1
                    label_key = "label-%09d".encode() % index
                    label = txn.get(label_key).decode("utf-8")

                    if len(label) > self.opt.batch_max_length:
                        continue

                    out_of_char = f"[^{self.opt.character}]"
                    if re.search(out_of_char, label.lower()):
                        continue

                    self.filtered_index_list.append(index)

                self.nSamples = len(self.filtered_index_list)

    def __len__(self):
        """
        Get the number of samples in the dataset.
        Returns:
            int: Number of samples.
        """
        return self.nSamples

    def __getitem__(self, index):
        """
        Get a sample from the LMDB dataset.
        Args:
            index (int): Index of the sample.
        Returns:
            tuple: Tuple containing the image (PIL Image object) and the corresponding label (string).
        """
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
                # Make dummy image and dummy label for corrupted image
                if self.opt.rgb:
                    img = Image.new("RGB", (self.opt.imgW, self.opt.imgH))
                else:
                    img = Image.new("L", (self.opt.imgW, self.opt.imgH))
                label = "[dummy_label]"

            if not self.opt.sensitive:
                label = label.lower()

            out_of_char = f"[^{self.opt.character}]"
            label = re.sub(out_of_char, "", label)

        return img, label

    def __call__(self):
        """
        Generate samples from the LMDB dataset.
        Yields:
            tuple: Tuple containing the image (PIL Image object) and the corresponding label (string).
        """
        for i in range(len(self)):
            yield self.__getitem__(i)


class SingleDataset(keras.utils.Sequence):
    def __init__(self, image, left, top, right, bottom, collate_fn):
        """
        Initializes the SingleDataset object.
        Args:
            image: The image data.
            left, top, right, bottom: Coordinates for cropping the image.
            collate_fn: A function to collate the preprocessed image data.
        """
        self.image = image
        self.left = left
        self.top = top
        self.right = right
        self.bottom = bottom
        self.collate_fn = collate_fn

    def __len__(self):
        """
        Returns the number of batches in the dataset.
        Returns:
            int: Number of batches in the dataset (always 1 in this case as it handles a single image).
        """
        return 1

    def __getitem__(self, index: int):
        """
        Retrieves a batch of preprocessed image data and labels.
        Args:
            index (int): Batch index (ignored in this implementation as there is only one batch).
        Returns:
            tuple: A tuple containing preprocessed image data and labels.
        """
        # Preprocess the single image using the specified coordinates
        image_preprocessed = all_preprocessing(
            self.image, self.left, self.top, self.right, self.bottom
        )
        # Apply the collate function to prepare the batch
        image_preprocessed = self.collate_fn([(image_preprocessed, "Prediction")])

        return image_preprocessed

    def __call__(self):
        """
        Generator function to yield preprocessed image data and labels.
        This function allows iterating over the dataset using a loop.
        Yields:
            tuple: A tuple containing preprocessed image data and labels.
        """
        for i in range(len(self)):
            yield self.__getitem__(i)


class RawDataset(keras.utils.Sequence):
    def __init__(self, root, opt):
        """
        Initializes the RawDataset object.
        Args:
            root (str): The root directory containing raw images.
            opt: An object containing various options and configurations.
        """
        self.opt = opt
        self.image_path_list = []

        # Collect image file paths with specified extensions from the root directory
        for dirpath, dirnames, filenames in os.walk(root):
            for name in filenames:
                _, ext = os.path.splitext(name)
                ext = ext.lower()
                if ext == ".jpg" or ext == ".jpeg" or ext == ".png":
                    self.image_path_list.append(os.path.join(dirpath, name))

        # Sort the collected image paths naturally
        self.image_path_list = natsorted(self.image_path_list)
        self.nSamples = len(self.image_path_list)

    def __len__(self):
        """
        Returns the number of samples in the dataset.
        Returns:
            int: Number of samples in the dataset.
        """
        return self.nSamples

    def __getitem__(self, index):
        """
        Retrieves an image and its corresponding file path.
        Args:
            index (int): Index of the sample in the dataset.
        Returns:
            tuple: A tuple containing the loaded image and its file path.
        """
        try:
            # Open the image file and convert it to the specified color mode (RGB or grayscale)
            if self.opt.rgb:
                img = Image.open(self.image_path_list[index]).convert("RGB")
            else:
                img = Image.open(self.image_path_list[index]).convert("L")

        except IOError:
            print(f"Corrupted image for {index}")
            # Create a dummy image with specified dimensions for corrupted images
            if self.opt.rgb:
                img = Image.new("RGB", (self.opt.imgW, self.opt.imgH))
            else:
                img = Image.new("L", (self.opt.imgW, self.opt.imgH))

        # Return the loaded image and its file path
        return (img, self.image_path_list[index])

    def __call__(self):
        """
        Generator function to yield images and their file paths.
        This function allows iterating over the dataset using a loop.
        Yields:
            tuple: A tuple containing the loaded image and its file path.
        """
        for i in range(len(self)):
            yield self.__getitem__(i)


class ResizeNormalize(object):
    def __init__(self, size, interpolation=Image.BICUBIC):
        """
        Initializes the ResizeNormalize object.
        Args:
            size (tuple): A tuple containing the target size (width, height) of the image.
            interpolation (int, optional): Interpolation method for resizing. Defaults to Image.BICUBIC.
        """
        self.size = size
        self.interpolation = interpolation
        self.toTensor = preprocessing.image.img_to_array

    def __call__(self, img):
        """
        Resizes the input image and performs pixel normalization.
        Args:
            img (PIL.Image): The input image.
        Returns:
            tf.Tensor: A tensor representing the resized and normalized image.
        """
        # Resize the image using the specified interpolation method
        img = img.resize(self.size, self.interpolation)

        # Convert the resized image to a NumPy array
        img = self.toTensor(img)

        # Normalize pixel values to the range [-1, 1]
        img = tf.math.divide(img, 255.0)  # Divide by 255.0 to bring values to [0, 1]
        img = tf.math.multiply(img, 2)  # Multiply by 2 to bring values to [0, 2]
        img = tf.math.subtract(img, 1)  # Subtract 1 to bring values to [-1, 1]

        return img


class NormalizePAD(object):
    def __init__(self, max_size, PAD_type="right"):
        """
        Initializes the NormalizePAD object.
        Args:
            max_size (tuple): A tuple containing the maximum size (channels, height, width) of the image after padding.
            PAD_type (str, optional): Specifies the type of padding. Defaults to "right".
                                      Possible values: "right" (pad on the right side) or "center" (pad on both sides equally).
        """
        self.toTensor = preprocessing.image.img_to_array
        self.max_size = max_size
        self.max_width_half = tf.floor(max_size[2] / 2)
        self.PAD_type = PAD_type

    def __call__(self, image: Image) -> tf.Tensor:
        """
        Pads the input image to the specified size.
        Args:
            image (PIL.Image): The input image.
        Returns:
            tf.Tensor: A tensor representing the padded image.
        """
        img = self.toTensor(image)
        c, h, w = img.shape

        # Create a tensor filled with zeros of the specified max_size
        Pad_img = tf.zeros(shape=self.max_size)

        # Copy the original image to the padded tensor based on PAD_type
        if self.PAD_type == "right":
            Pad_img[:, :, :w] = img
        elif self.PAD_type == "center":
            start = int(self.max_width_half - tf.floor(w / 2))
            end = start + w
            Pad_img[:, :, start:end] = img

        return Pad_img


class AlignCollate(object):
    def __init__(self, imgH=32, imgW=100, keep_ratio_with_pad=False):
        """
        Initializes the AlignCollate object.
        Args:
            imgH (int, optional): The height of the output images after alignment. Defaults to 32.
            imgW (int, optional): The width of the output images after alignment. Defaults to 100.
            keep_ratio_with_pad (bool, optional): Whether to keep the aspect ratio and pad the images.
                                                  Defaults to False.
        """
        self.imgH = imgH
        self.imgW = imgW
        self.keep_ratio_with_pad = keep_ratio_with_pad

    def __call__(self, batch):
        """
        Aligns and collates a batch of images and labels.
        Args:
            batch (list): A list of tuples, where each tuple contains an image and its corresponding label.
        Returns:
            tuple: A tuple containing the aligned and collated image tensors and corresponding labels.
        """
        batch = filter(lambda x: x is not None, batch)
        images, labels = zip(*batch)

        if self.keep_ratio_with_pad:
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
                    indices[offset - length: offset],
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

        # Iterate through data loader iterators and create a balanced batch
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
        # Concatenate the images along the batch dimension
        balanced_batch_images = tf.concat(balanced_batch_images, axis=0)

        return balanced_batch_images, balanced_batch_texts
