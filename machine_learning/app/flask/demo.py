import argparse

import numpy as np
import tensorflow as tf
import tensorflow_addons as tfa
from PIL import Image
from tensorflow import keras

from dataset import AlignCollate, SingleDataset, tensorflow_dataloader
from utils import CTCLabelConverter, read_json, show_normalized_image


def demo(opt):
    converter = CTCLabelConverter("0123456789")

    # read json
    file_json = read_json(opt.json)
    # load model
    model = keras.models.load_model(
        opt.saved_model, custom_objects={"AAP": tfa.layers.AdaptiveAveragePooling2D}
    )
    # prepare data. two demo images from https://github.com/bgshih/crnn#run-demo
    AlignCollate_demo = AlignCollate(imgH=32, imgW=100, keep_ratio_with_pad=False)
    image_path = file_json["image"]
    Ymin = file_json["class"]["NIK"]["Ymin"]
    Xmin = file_json["class"]["NIK"]["Xmin"]
    Xmax = file_json["class"]["NIK"]["Xmax"]
    Ymax = file_json["class"]["NIK"]["Ymax"]
    demo_data = SingleDataset(
        image=Image.open(image_path),
        opt=opt,
        left=Xmin,
        top=Ymin,
        right=Xmax,
        bottom=Ymax,
        collate_fn=AlignCollate_demo,
    )
    demo_loader = tensorflow_dataloader(
        demo_data,
        batch_size=opt.batch_size,
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

    print(f"{preds_str}")


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
