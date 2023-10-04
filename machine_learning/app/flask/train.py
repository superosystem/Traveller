import argparse
import os
import random
import string
import sys
import time
from test import validation

import numpy as np
import tensorflow as tf
from PIL import Image
from tensorflow import keras

from dataset import (
    AlignCollate,
    ApplyCollate,
    Batch_Balanced_Dataset,
    hierarchical_dataset,
    tensorflow_dataloader,
)
from model import Model
from modules.custom import custom_sparse_categorical_crossentropy
from utils import Averager, CTCLabelConverter, CTCLabelConverterForBaiduWarpctc


def ignore_index(tensor, ignored_index=0):
    return tf.cast(tf.not_equal(tensor, ignored_index), tf.float32)


def train(opt):
    """dataset preparation"""
    # show preferred channel order
    if not opt.data_filtering_off:
        print(
            "Filtering the images containing characters which are not in opt.character"
        )
        print("Filtering the images whose label is longer than opt.batch_max_length")

    opt.select_data = opt.select_data.split("-")
    opt.batch_ratio = opt.batch_ratio.split("-")

    train_dataset = Batch_Balanced_Dataset(opt)

    log = open(f"./saved_models/{opt.exp_name}/log_dataset.txt", "a")
    AlignCollate_valid = AlignCollate(
        imgH=opt.imgH, imgW=opt.imgW, keep_ratio_with_pad=opt.PAD
    )
    valid_dataset, valid_dataset_log = hierarchical_dataset(
        root=opt.valid_data, opt=opt
    )
    valid_dataset_dataset = ApplyCollate(valid_dataset, AlignCollate_valid)
    valid_loader = tensorflow_dataloader(
        valid_dataset_dataset,
        batch_size=opt.batch_size,
        shuffle=True,
        num_workers=int(opt.workers),
        collate_fn=AlignCollate_valid,
        imgH=opt.imgH,
        imgW=opt.imgW
    )
    log.write(valid_dataset_log)
    print("-" * 80)
    log.write("-" * 80 + "\n")
    log.close()

    """ model configuration """
    if opt.baiduCTC:
        converter = CTCLabelConverterForBaiduWarpctc(opt.character)
    else:
        converter = CTCLabelConverter(opt.character)

    opt.num_class = len(converter.character)

    if opt.rgb:
        opt.input_channel = 3
    model = Model(opt)

    opt.num_class = len(converter.character)
    print(f"opt num_class is {opt.num_class}")

    # weight initialization
    # TODO if model behave bad

    if opt.saved_model != "":
        print(f"loading pretrained model from {opt.saved_model}")
        model = keras.models.load_model(opt.saved_model)
    print("Model → ")

    if "CTC" in opt.Prediction:
        criterion = tf.nn.ctc_loss
    else:
        # kalo gak yang sparse ya yang biasa
        # selalu gunakan ignored_index=0 disini
        criterion = custom_sparse_categorical_crossentropy

    loss_avg = Averager()

    if opt.adam:
        optimizer = keras.optimizers.Adam(
            learning_rate=opt.lr, beta_1=opt.beta1, beta_2=0.999
        )
    else:
        optimizer = keras.optimizers.Adadelta(
            learning_rate=opt.lr, rho=opt.rho, epsilon=opt.eps
        )

    print("Optimizer:")
    print(optimizer)

    """ final options """
    # print(opt)
    with open(f"./saved_models/{opt.exp_name}/opt.txt", "a") as opt_file:
        opt_log = "------------ Options -------------\n"
        args = vars(opt)
        for k, v in args.items():
            opt_log += f"{str(k)}: {str(v)}\n"
        opt_log += "---------------------------------------\n"
        print(opt_log)
        opt_file.write(opt_log)
    """ start training """
    start_iter = 0
    if opt.saved_model != "":
        try:
            start_iter = int(opt.saved_model.split("_")[-1].split(".")[0])
            print(f"continue to train, start_iter: {start_iter}")
        except:
            pass

    start_time = time.time()
    best_accuracy = -1
    best_norm_ED = -1
    iteration = start_iter

    strategy = tf.distribute.MirroredStrategy()

    while True:
        image_tensors, labels = train_dataset.get_batch()
        image = image_tensors
        if opt.pretrained:
            image = tf.image.grayscale_to_rgb(image)
        # image = tf.transpose(image, perm=[0, 3, 1, 2])
        labels = labels[0].copy()
        labels[0] = str(labels[0], "utf-8")
        text, length = converter.encode(labels, batch_max_length=opt.batch_max_length)
        batch_size = image.shape[0]

        if "CTC" in opt.Prediction:
            with tf.GradientTape() as tape:
                preds = model(image, text, training=True)

                preds_size = tf.constant([preds.shape[1]] * batch_size)
                if opt.baiduCTC:
                    preds = tf.transpose(preds, perm=[1, 0, 2])
                    cost = (
                        criterion(
                            labels=preds,
                            logits=text,
                            label_length=preds_size,
                            logit_length=length,
                        )
                        / batch_size
                    )
                else:
                    # preds = tf.nn.log_softmax(preds, axis=2)
                    preds = tf.transpose(preds, perm=[1, 0, 2])

                    text = tf.cast(text, dtype=tf.int32)
                    cost = criterion(
                        logits=preds,
                        labels=text,
                        logit_length=preds_size,
                        label_length=length,
                        blank_index=0,
                    )
                    # zero_infinity
                    cost = tf.where(tf.math.is_inf(cost), tf.zeros_like(cost), cost)
                    # cost = tf.math.log(cost)
                    cost = tf.reduce_mean(cost, axis=-1)
        # this could be total mess
        gradients = tape.gradient(cost, model.trainable_variables)
        gradients, _ = tf.clip_by_global_norm(gradients, opt.grad_clip)
        optimizer.apply_gradients(zip(gradients, model.trainable_variables))

        loss_avg.add(cost)
        # validation part

        if (
            iteration + 1
        ) % opt.valInterval == 0 or iteration == 0:  # To see training progress, we also conduct validation when 'iteration == 0'
            elapsed_time = time.time() - start_time

            with open(f"./saved_models/{opt.exp_name}/log_train.txt", "a") as log:
                (
                    valid_loss,
                    current_accuracy,
                    current_norm_ED,
                    preds,
                    confidence_score,
                    labels,
                    infer_time,
                    length_of_data,
                ) = validation(model, criterion, valid_loader, converter, opt)
                loss_log = f"[{iteration+1}/{opt.num_iter}] Train loss: {loss_avg.val():0.5f}, Valid loss: {valid_loss:0.5f}, Elapsed_time: {elapsed_time:0.5f}"
                loss_avg.reset()

                current_model_log = f'{"Current_accuracy":17s}: {current_accuracy:0.3f}, {"Current_norm_ED":17s}: {current_norm_ED:0.2f}'

                if current_accuracy > best_accuracy:
                    best_accuracy = current_accuracy
                    model.save(f"./saved_models/{opt.exp_name}/best_accuracy")
                if current_norm_ED > best_norm_ED:
                    best_norm_ED = current_norm_ED
                    model.save(f"./saved_models/{opt.exp_name}/best_norm_ED")

                best_model_log = f'{"Best_accuracy":17s}: {best_accuracy:0.3f}, {"Best_norm_ED":17s}: {best_norm_ED:0.2f}'

                loss_model_log = f"{loss_log}\n{current_model_log}\n{best_model_log}"
                print(loss_model_log)
                log.write(loss_model_log + "\n")

                dashed_line = "-" * 80
                head = f'{"Ground Truth":25s} | {"Prediction":25s} | Confidence Score & T/F'
                predicted_result_log = f"{dashed_line}\n{head}\n{dashed_line}\n"
                for gt, pred, confidence in zip(
                    labels[:5], preds[:5], confidence_score[:5]
                ):
                    if "Attn" in opt.Prediction:
                        gt = gt[: gt.find("[s]")]
                        pred = pred[: pred.find("[s]")]

                    predicted_result_log += f"{gt:25s} | {pred:25s} | {confidence:0.4f}\t{str(pred == gt)}\n"
                predicted_result_log += f"{dashed_line}"
                print(predicted_result_log)
                log.write(predicted_result_log + "\n")

        # save model per 1e+5 iter.
        if (iteration + 1) % 1e5 == 0:
            model.save(f"./saved_models/{opt.exp_name}/iter_{iteration+1}")

        if (iteration + 1) == opt.num_iter:
            print("end the training")
            model.save(f"./saved_models/{opt.exp_name}/;ast_iter")
            sys.exit()
        iteration += 1


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--exp_name", help="Where to store logs and models")
    parser.add_argument("--train_data", required=True, help="path to training dataset")
    parser.add_argument(
        "--valid_data", required=True, help="path to validation dataset"
    )
    parser.add_argument(
        "--manualSeed", type=int, default=1111, help="for random seed setting"
    )
    parser.add_argument(
        "--workers", type=int, help="number of data loading workers", default=4
    )
    parser.add_argument("--batch_size", type=int, default=192, help="input batch size")
    parser.add_argument(
        "--num_iter", type=int, default=300000, help="number of iterations to train for"
    )
    parser.add_argument(
        "--valInterval", type=int, default=2000, help="Interval between each validation"
    )
    parser.add_argument(
        "--saved_model", default="", help="path to model to continue training"
    )
    parser.add_argument("--FT", action="store_true", help="whether to do fine-tuning")
    parser.add_argument(
        "--adam", action="store_true", help="Whether to use adam (default is Adadelta)"
    )
    parser.add_argument(
        "--lr", type=float, default=1, help="learning rate, default=1.0 for Adadelta"
    )
    parser.add_argument(
        "--beta1", type=float, default=0.9, help="beta1 for adam. default=0.9"
    )
    parser.add_argument(
        "--rho",
        type=float,
        default=0.95,
        help="decay rate rho for Adadelta. default=0.95",
    )
    parser.add_argument(
        "--eps", type=float, default=1e-8, help="eps for Adadelta. default=1e-8"
    )
    parser.add_argument(
        "--grad_clip", type=float, default=5, help="gradient clipping value. default=5"
    )
    parser.add_argument(
        "--baiduCTC", action="store_true", help="for data_filtering_off mode"
    )
    parser.add_argument(
        "--pretrained",
        action="store_true",
        help="Use pretrained feature extractor from Keras with architecture that choosen in feature extractor arguments",
    )
    """ Data processing """
    parser.add_argument(
        "--select_data",
        type=str,
        default="MJ-ST",
        help="select training data (default is MJ-ST, which means MJ and ST used as training data)",
    )
    parser.add_argument(
        "--batch_ratio",
        type=str,
        default="0.5-0.5",
        help="assign ratio for each selected data in the batch",
    )
    parser.add_argument(
        "--total_data_usage_ratio",
        type=str,
        default="1.0",
        help="total data usage ratio, this ratio is multiplied to total number of data.",
    )
    parser.add_argument(
        "--batch_max_length", type=int, default=25, help="maximum-label-length"
    )
    parser.add_argument(
        "--imgH", type=int, default=32, help="the height of the input image"
    )
    parser.add_argument(
        "--imgW", type=int, default=100, help="the width of the input image"
    )
    parser.add_argument("--rgb", action="store_true", help="use rgb input")
    parser.add_argument(
        "--character",
        type=str,
        default="0123456789abcdefghijklmnopqrstuvwxyz",
        help="character label",
    )
    parser.add_argument(
        "--sensitive", action="store_true", help="for sensitive character mode"
    )
    parser.add_argument(
        "--PAD",
        action="store_true",
        help="whether to keep ratio then pad for image resize",
    )
    parser.add_argument(
        "--data_filtering_off", action="store_true", help="for data_filtering_off mode"
    )
    """ Model Architecture """
    parser.add_argument(
        "--Transformation",
        type=str,
        required=True,
        help="Transformation stage. None|TPS",
    )
    parser.add_argument(
        "--FeatureExtraction",
        type=str,
        required=True,
        help="FeatureExtraction stage. VGG|RCNN|ResNet",
    )
    parser.add_argument(
        "--SequenceModeling",
        type=str,
        required=True,
        help="SequenceModeling stage. None|BiLSTM",
    )
    parser.add_argument(
        "--Prediction", type=str, required=True, help="Prediction stage. CTC|Attn"
    )
    parser.add_argument(
        "--num_fiducial",
        type=int,
        default=20,
        help="number of fiducial points of TPS-STN",
    )
    parser.add_argument(
        "--input_channel",
        type=int,
        default=1,
        help="the number of input channel of Feature extractor",
    )
    parser.add_argument(
        "--output_channel",
        type=int,
        default=512,
        help="the number of output channel of Feature extractor",
    )
    parser.add_argument(
        "--hidden_size", type=int, default=256, help="the size of the LSTM hidden state"
    )

    opt = parser.parse_args()

    if not opt.exp_name:
        opt.exp_name = f"{opt.Transformation}-{opt.FeatureExtraction}-{opt.SequenceModeling}-{opt.Prediction}"
        opt.exp_name += f"-Seed{opt.manualSeed}"
        # print(opt.exp_name)

    os.makedirs(f"./saved_models/{opt.exp_name}", exist_ok=True)

    """ vocab / character number configuration """
    if opt.sensitive:
        # opt.character += 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
        opt.character = string.printable[:-6]  # same with ASTER setting (use 94 char).

    """ Seed and GPU setting """
    # print("Random Seed: ", opt.manualSeed)
    random.seed(opt.manualSeed)
    np.random.seed(opt.manualSeed)
    tf.random.set_seed(opt.manualSeed)

    opt.num_gpu = len(tf.config.list_physical_devices("GPU"))
    # print('device count', opt.num_gpu)
    if opt.num_gpu > 1:
        print("------ Use multi-GPU setting ------")
        print(
            "if you stuck too long time with multi-GPU setting, try to set --workers 0"
        )
        # check multi-GPU issue https://github.com/clovaai/deep-text-recognition-benchmark/issues/1
        opt.workers = opt.workers * opt.num_gpu
        opt.batch_size = opt.batch_size * opt.num_gpu

        """ previous version
        print('To equlize batch stats to 1-GPU setting, the batch_size is multiplied with num_gpu and multiplied batch_size is ', opt.batch_size)
        opt.batch_size = opt.batch_size * opt.num_gpu
        print('To equalize the number of epochs to 1-GPU setting, num_iter is divided with num_gpu by default.')
        If you dont care about it, just commnet out these line.)
        opt.num_iter = int(opt.num_iter / opt.num_gpu)
        """

    train(opt)
