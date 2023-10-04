import argparse
import os
import re
import string
import time

import tensorflow as tf
from nltk.metrics.distance import edit_distance
from tensorflow import keras

from dataset import AlignCollate, hierarchical_dataset, tensorflow_dataloader
from model import Model
from modules.custom import custom_sparse_categorical_crossentropy
from utils import Averager, CTCLabelConverter


def benchmark_all_eval(model, criterion, converter, opt, calculate_infer_time=False):
    return None


def validation(model, criterion, evaluation_loader, converter, opt):
    """validation or evaluation"""
    n_correct = 0
    norm_ED = 0
    length_of_data = 0
    infer_time = 0
    valid_loss_avg = Averager()

    for i, (image_tensors, labels) in enumerate(evaluation_loader):
        batch_size = image_tensors.shape[0]
        length_of_data = length_of_data + batch_size
        image = image_tensors
        if opt.pretrained:
            image = tf.image.grayscale_to_rgb(image)

        length_for_pred = tf.constant(
            [opt.batch_max_length] * batch_size, dtype=tf.int32
        )
        text_for_pred = tf.zeros(
            shape=[batch_size, opt.batch_max_length + 1], dtype=tf.float64
        )

        labels = labels[0].numpy()
        labels[0] = str(labels[0], "utf-8")

        text_for_loss, length_for_loss = converter.encode(
            labels, batch_max_length=opt.batch_max_length
        )
        start_time = time.time()

        if "CTC" in opt.Prediction:
            preds = model(image, text_for_pred)
            forward_time = time.time() - start_time
            preds_size = tf.constant([preds.shape[1]] * batch_size, dtype=tf.int32)
            if opt.baiduCTC:
                # preds = tf.transpose(preds, perm=[1, 0, 2])
                cost = (
                    criterion(
                        labels=preds,
                        logits=text_for_loss,
                        label_length=preds_size,
                        logit_length=length_for_loss,
                    )
                    / batch_size
                )
            else:
                # preds = tf.nn.log_softmax(preds, axis=2)
                preds = tf.transpose(preds, perm=[1, 0, 2])
                # preds = tf.math.log(preds)
                text_for_loss = tf.cast(text_for_loss, dtype=tf.int32)
                cost = criterion(
                    logits=preds,
                    labels=text_for_loss,
                    logit_length=preds_size,
                    label_length=length_for_loss,
                    blank_index=0,
                )
                cost = tf.where(tf.math.is_inf(cost), tf.zeros_like(cost), cost)
                cost = tf.reduce_mean(cost)

            # Select max probabilty (greedy decoding) then decode index to character
            if opt.baiduCTC:
                preds_index = tf.math.argmax(preds, axis=2)
                preds_index = tf.reshape(preds_index, shape=[-1])
            else:
                preds_index = tf.math.argmax(preds, axis=2)
            preds_str = converter.decode(preds_index.numpy().T, preds_size.numpy())
        else:
            preds = model(image, text_for_pred, training=False)
            forward_time = time.time() - start_time

            preds = preds[:, : text_for_loss.shape[1] - 1, :]
            target = text_for_loss[:, 1:]
            cost = criterion(
                tf.reshape(preds, shape=[-1, preds.shape[-1]]),
                tf.reshape(target, shape=[-1]),
            )

            preds_index = tf.math.argmax(preds, axis=2)
            preds_str = converter.decode(preds_index, length_for_pred)
            labels = converter.decode(text_for_loss[:, 1:], length_for_loss)

        infer_time += forward_time
        valid_loss_avg.add(cost)

        preds_probs = tf.nn.softmax(preds, axis=2)
        preds_max_prob = tf.math.reduce_max(preds_probs, axis=2)

        confidence_score_list = []
        for gt, pred, pred_max_prob in zip(labels, preds_str, preds_max_prob):
            if "Attn" in opt.Prediction:
                gt = gt[: gt.find("[s]")]
                pred_EOS = pred.find("[s]")
                pred = pred[:pred_EOS]  # prune after "end of sentence" token ([s])
                pred_max_prob = pred_max_prob[:pred_EOS]

            # To evaluate "case sensitive model" with alphanumeric and case insensitve setting.
            if opt.sensitive and opt.data_filtering_off:
                pred = pred.lower()
                gt = gt.lower()
                alphanumeric_case_insensitve = "0123456789abcdefghijklmnopqrstuvwxyz"
                out_of_alphanumeric_case_insensitve = (
                    f"[^{alphanumeric_case_insensitve}]"
                )
                pred = re.sub(out_of_alphanumeric_case_insensitve, "", pred)
                gt = re.sub(out_of_alphanumeric_case_insensitve, "", gt)

            if pred == gt:
                n_correct += 1

            """
            (old version) ICDAR2017 DOST Normalized Edit Distance https://rrc.cvc.uab.es/?ch=7&com=tasks
            "For each word we calculate the normalized edit distance to the length of the ground truth transcription."
            if len(gt) == 0:
                norm_ED += 1
            else:
                norm_ED += edit_distance(pred, gt) / len(gt)
            """

            # ICDAR2019 Normalized Edit Distance
            if len(gt) == 0 or len(pred) == 0:
                norm_ED += 0
            elif len(gt) > len(pred):
                norm_ED += 1 - edit_distance(pred, gt) / len(gt)
            else:
                norm_ED += 1 - edit_distance(pred, gt) / len(pred)

            # calculate confidence score (= multiply of pred_max_prob)
            try:
                confidence_score = tf.math.cumprod(pred_max_prob, axis=0)[-1]
            except:
                confidence_score = 0  # for empty pred case, when prune after "end of sentence" token ([s])
            confidence_score_list.append(confidence_score)
            # print(pred, gt, pred==gt, confidence_score)

    accuracy = n_correct / float(length_of_data) * 100
    norm_ED = norm_ED / float(length_of_data)  # ICDAR2019 Normalized Edit Distance

    return (
        valid_loss_avg.val(),
        accuracy,
        norm_ED,
        preds_str,
        confidence_score_list,
        labels,
        infer_time,
        length_of_data,
    )


def test(opt):
    """model configuration"""
    converter = CTCLabelConverter(opt.character)
    opt.num_class = len(converter.character)

    if opt.rgb:
        opt.input_channel = 3

    model = Model(opt)

    print("loading pretrained model from %s" % opt.saved_model)
    model = keras.models.load_model(opt.saved_model)
    opt.exp_name = "_".join(opt.saved_model.split("/")[1:])

    os.makedirs(f"./result/{opt.exp_name}", exist_ok=True)
    os.system(f"cp {opt.saved_model} ./result/{opt.exp_name}/")

    if "CTC" in opt.Prediction:
        criterion = tf.nn.ctc_loss
    else:
        # kalo gak yang sparse ya yang biasa
        # selalu gunakan ignored_index=0 disini
        criterion = custom_sparse_categorical_crossentropy

    log = open(f"./result/{opt.exp_name}/log_evaluation.txt", "a")
    AlignCollate_evaluation = AlignCollate(
        imgH=opt.imgH, imgW=opt.imgW, keep_ratio_with_pad=opt.PAD
    )
    eval_data, eval_data_log = hierarchical_dataset(root=opt.eval_data, opt=opt)
    evaluation_loader = tensorflow_dataloader(
        eval_data,
        batch_size=opt.batch_size,
        shuffle=False,
        num_workers=int(opt.workers),
        collate_fn=AlignCollate_evaluation,
    )
    _, accuracy_by_best_model, _, _, _, _, _, _ = validation(
        model, criterion, evaluation_loader, converter, opt
    )
    log.write(eval_data_log)
    print(f"{accuracy_by_best_model:0.3f}")
    log.write(f"{accuracy_by_best_model:0.3f}\n")
    log.close()


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--eval_data", required=True, help="path to evaluation dataset")
    parser.add_argument(
        "--benchmark_all_eval",
        action="store_true",
        help="evaluate 10 benchmark evaluation datasets",
    )
    parser.add_argument(
        "--workers", type=int, help="number of data loading workers", default=4
    )
    parser.add_argument("--batch_size", type=int, default=192, help="input batch size")
    parser.add_argument(
        "--saved_model", required=True, help="path to saved_model to evaluation"
    )
    """ Data processing """
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
    parser.add_argument(
        "--baiduCTC", action="store_true", help="for data_filtering_off mode"
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

    """ vocab / character number configuration """
    if opt.sensitive:
        opt.character = string.printable[:-6]  # same with ASTER setting (use 94 char).

    opt.num_gpu = len(tf.config.list_physical_devices("GPU"))

    test(opt)
