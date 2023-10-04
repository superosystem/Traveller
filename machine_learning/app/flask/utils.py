import json
from typing import Dict

import numpy as np
import tensorflow as tf
from PIL import Image


def read_json(path: str) -> Dict:
    with open(path, "r") as openfile:
        file_json = json.load(openfile)

    return file_json


def show_normalized_image(image) -> None:
    image_numpy = (((image + 1) / 2) * 255).astype(np.uint8)
    image_numpy = np.squeeze(image_numpy[0], 2)
    Image.fromarray(image_numpy).show()


class CTCLabelConverter(object):
    """Convert between text-label and text-index"""

    def __init__(self, character):
        dict_character = list(character)

        self.dict = {}
        for i, char in enumerate(dict_character):
            # 0 is reserved for "CTCBlank" token required by CTCloss
            self.dict[char] = i + 1

        self.character = ["[CTCBlank]"] + dict_character

    def encode(self, text, batch_max_length):
        """convert text-label into text-index.
        input:
            text: text labels of each image. [batch_size]
            batch_max_length: max length of text label in the batch. 25 by default
        output:
            text: text index for CTCLoss. [batch_size, batch_max_length]
            length: length of each text. [batch_size]
        """
        length = [len(s) for s in text]
        # The index used for padding (=0) would not affect the CTC loss calculation.
        batch_text = tf.Variable(
            tf.zeros(shape=[len(text), batch_max_length], dtype=tf.float64)
        )

        output_list = []
        for i, t in enumerate(text):
            text = list(t)
            text = [self.dict.get(char, -1) for char in text]
            # print(len(text))
            text.extend([0] * (batch_max_length - len(text)))
            output_list.append(tf.constant(text, dtype=tf.float64))

        batch_text = tf.stack(output_list)
        return (batch_text, tf.constant(length, dtype=tf.int32))

    def decode(self, text_index, length):
        """convert text-index into text-label."""
        texts = []
        for index, l in enumerate(length):
            t = text_index[index, :]
            char_list = []
            for i in range(l):
                if t[i] != 0 and (not (i > 0 and t[i - 1] == t[i])):
                    char_list.append(self.character[t[i]])
            text = "".join(char_list)

            texts.append(text)
        return texts


class CTCLabelConverterForBaiduWarpctc(object):
    """Convert between text-label and text-index for baidu warpctc"""

    def __init__(self, character):
        # character (str): set of the possible characters.
        dict_character = list(character)

        self.dict = {}
        for i, char in enumerate(dict_character):
            # NOTE: 0 is reserved for 'CTCblank' token required by CTCLoss
            self.dict[char] = i + 1

        self.character = [
            "[CTCblank]"
        ] + dict_character  # dummy '[CTCblank]' token for CTCLoss (index 0)

    def encode(self, text, batch_max_length=25):
        """convert text-label into text-index.
        input:
            text: text labels of each image. [batch_size]
        output:
            text: concatenated text index for CTCLoss.
                    [sum(text_lengths)] = [text_index_0 + text_index_1 + ... + text_index_(n - 1)]
            length: length of each text. [batch_size]
        """
        length = [len(s) for s in text]
        text = "".join(text)
        text = [self.dict[char] for char in text]

        return (tf.constant(text, dtype=tf.int32), tf.constant(length, dtype=tf.int32))

    def decode(self, text_index, length):
        """convert text-index into text-label."""
        texts = []
        index = 0
        for l in length:
            t = text_index[index : index + l]

            char_list = []
            for i in range(l):
                if t[i] != 0 and (
                    not (i > 0 and t[i - 1] == t[i])
                ):  # removing repeated characters and blank.
                    char_list.append(self.character[t[i]])
            text = "".join(char_list)

            texts.append(text)
            index += l
        return texts


class Averager(object):
    """Compute average for tf.constant, used for loss average."""

    def __init__(self):
        self.reset()

    def reset(self):
        self.n_count = 0
        self.sum = 0

    def add(self, v):
        count = tf.size(v)
        v = tf.reduce_sum(v)
        self.n_count += count
        self.sum += v

    def val(self):
        res = 0
        if self.n_count != 0:
            res = self.sum / float(self.n_count)
        return res
