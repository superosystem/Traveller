import argparse

import numpy as np
import tensorflow as tf
import tensorflow_addons as tfa
from PIL import Image
from tensorflow import keras

import json

from dataset import AlignCollate, SingleDataset, tensorflow_dataloader
from utils import CTCLabelConverter, read_json, show_normalized_image


def getObject(file_json, label):
    objects = []
    objects.append(file_json["image"])
    objects.append(file_json["class"][label]["Xmin"])
    objects.append(file_json["class"][label]["Ymin"])
    objects.append(file_json["class"][label]["Xmax"])
    objects.append(file_json["class"][label]["Ymax"])
    return objects  # [image, Xmin, Ymin, Xmax, Ymax]


def predict_nik(saved_model, json_file):
    converter = CTCLabelConverter("0123456789")
    saved_model = saved_model + "/NIK_Model/best_accuracy"

    model = keras.models.load_model(
        saved_model, custom_objects={"AAP": tfa.layers.AdaptiveAveragePooling2D}
    )
    # prepare data. two demo images from https://github.com/bgshih/crnn#run-demo
    AlignCollate_demo = AlignCollate(imgH=32, imgW=100, keep_ratio_with_pad=False)
    obj = getObject(json_file, "NIK")
    demo_data = SingleDataset(
        image=Image.open(obj[0]),
        left=obj[1],
        top=obj[2],
        right=obj[3],
        bottom=obj[4],
        collate_fn=AlignCollate_demo,
    )
    demo_loader = tensorflow_dataloader(
        demo_data,
        batch_size=1,
        shuffle=True,
        num_workers=4,
        collate_fn=AlignCollate_demo,
    )
    image, text_for_pred = next(demo_loader.as_numpy_iterator())
    batch_size = image.shape[0]
    text_for_pred = tf.zeros(shape=(batch_size, 25), dtype=tf.float64)
    length_for_pred = tf.constant([25] * batch_size, dtype=tf.int32)

    preds = model(image, text_for_pred)

    # Select max probabilty (greedy decoding) then decode index to character
    preds_size = tf.constant([preds.shape[1]] * batch_size)
    preds_index = tf.math.argmax(preds, axis=2)
    # preds_index = preds_index.view(-1)
    preds_str = converter.decode(preds_index, preds_size)
    nik = preds_str[0]
    return nik  # String NIK


def predict_arial(saved_model, json_file):
    converter = CTCLabelConverter("abcdefghijklmnopqrstuvwxyz,. -")
    saved_model = saved_model + "/Arial_Model/last_iter"
    model = keras.models.load_model(
        saved_model, custom_objects={"AAP": tfa.layers.AdaptiveAveragePooling2D}
    )

    AlignCollate_demo = AlignCollate(imgH=64, imgW=200, keep_ratio_with_pad=False)

    labels = ["name", "sex", "married", "nationality"]
    objs = []
    for label in labels:
        objs.append(getObject(json_file, label))

    demo_datas = []
    for obj in objs:
        demo_data = SingleDataset(
            image=Image.open(obj[0]),
            left=obj[1],
            top=obj[2],
            right=obj[3],
            bottom=obj[4],
            collate_fn=AlignCollate_demo,
        )
        demo_datas.append(demo_data)

    arials = []
    for demo_data in demo_datas:
        demo_loader = tensorflow_dataloader(
            demo_data,
            batch_size=1,
            shuffle=True,
            num_workers=4,
            imgH=64,
            imgW=200,
            collate_fn=AlignCollate_demo,
        )
        image, text_for_pred = next(demo_loader.as_numpy_iterator())

        show_normalized_image(image)
        batch_size = image.shape[0]
        text_for_pred = tf.zeros(shape=(batch_size, 25), dtype=tf.float64)
        length_for_pred = tf.constant([25] * batch_size, dtype=tf.int32)

        preds = model(image, text_for_pred)

        # Select max probabilty (greedy decoding) then decode index to character
        preds_size = tf.constant([preds.shape[1]] * batch_size)
        preds_index = tf.math.argmax(preds, axis=2)
        # preds_index = preds_index.view(-1)
        preds_str = converter.decode(preds_index, preds_size)

        arials.append(preds_str[0].upper())

    return arials  # [name, sex, married, nationality]


def demo(opt):
    nik = predict_nik(opt.saved_model, read_json(opt.json))
    arials = predict_arial(opt.saved_model, read_json(opt.json))
    dict = {
        "nik": nik,
        "name": arials[0],
        "sex": arials[1],
        "married": arials[2],
        "nationality": arials[3],
    }
    print(dict)

    output_path = "./testing/output.json"
    with open(output_path, "w") as outfile:
        json.dump(dict, outfile)
    print(f"Successfully create JSON object at {output_path}")


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--batch_size", type=int, default=192, help="input batch size")
    parser.add_argument(
        "--saved_model",
        type=str,
        required=True,
        help="path to saved_model to evaluation",
    )
    parser.add_argument(
        "--json",
        type=str,
        required=True,
        help="JSON file contain image path and bounding box coordinate",
    )

    opt = parser.parse_args()

    opt.num_gpu = len(tf.config.list_physical_devices("GPU"))

    demo(opt)
