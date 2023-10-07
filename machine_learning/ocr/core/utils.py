import json
from typing import Dict

import numpy as np
import tensorflow as tf
from PIL import Image


def read_json(path: str) -> Dict:
    """
    Reads a JSON file from the specified path and returns its content as a Python dictionary.
    Args:
       path (str): The file path to the JSON file.
    Returns:
       Dict: A dictionary containing the JSON data.
   """
    with open(path, "r") as openfile:
        file_json = json.load(openfile)

    return file_json


def show_normalized_image(image) -> None:
    """
    Displays a normalized image represented as a NumPy array using the PIL library.

    Args:
       image (numpy.ndarray): Normalized image data, typically in the range [-1, 1].

    Returns:
       None
    """
    # Convert the normalized image to the range [0, 255] and as unsigned 8-bit integers
    image_numpy = (((image + 1) / 2) * 255).astype(np.uint8)
    # Remove single-dimensional entries from the shape of the array
    image_numpy = np.squeeze(image_numpy[0], 2)
    # Create an image from the array and display it
    Image.fromarray(image_numpy).show()


class CTCLabelConverter(object):
    """
    This class converts between text-label and text-index for CTC (Connectionist Temporal Classification) loss.
    """

    def __init__(self, character):
        """
        Initializes the CTCLabelConverter object.

        Args:
            character (str): A string containing all the unique characters in the dataset.
        """
        dict_character = list(character)

        # Create a dictionary to map characters to unique indices starting from 1
        self.dict = {}
        for i, char in enumerate(dict_character):
            # 0 is reserved for "CTCBlank" token required by CTC loss
            self.dict[char] = i + 1

        # Character list including "[CTCBlank]" token at index 0
        self.character = ["[CTCBlank]"] + dict_character

    def encode(self, text, batch_max_length):
        """
        Convert text labels into text indices for CTCLoss calculation.
        Args:
            text (list): Text labels of each image. [batch_size]
            batch_max_length (int): Maximum length of text label in the batch. Default is 25.
        Returns:
            tuple: A tuple containing two elements:
                - batch_text (tf.Tensor): Text indices for CTCLoss. Shape: [batch_size, batch_max_length]
                - length (tf.Tensor): Length of each text. Shape: [batch_size]
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
            # Pad with 0s to match the batch_max_length
            text.extend([0] * (batch_max_length - len(text)))
            output_list.append(tf.constant(text, dtype=tf.float64))

        batch_text = tf.stack(output_list)
        return batch_text, tf.constant(length, dtype=tf.int32)

    def decode(self, text_index, length):
        """
        Convert text indices back into text labels.
        Args:
            text_index (tf.Tensor): Text indices for CTCLoss. Shape: [batch_size, batch_max_length]
            length (tf.Tensor): Length of each text. Shape: [batch_size]
        Returns:
            list: List of text labels for each entry in the batch.
        """
        texts = []
        for index, l in enumerate(length):
            t = text_index[index, :]
            char_list = []
            for i in range(l):
                # Ignore padding (index 0) and repeated characters
                if t[i] != 0 and (not (i > 0 and t[i - 1] == t[i])):
                    char_list.append(self.character[t[i]])
            text = "".join(char_list)
            texts.append(text)
        return texts


class CTCLabelConverterForBaiduWarpctc(object):
    """
    This class provides methods to convert text-labels to text-indices and vice versa for Baidu's Warp-CTC library.
    """

    def __init__(self, character):
        """
        Initializes the CTCLabelConverterForBaiduWarpctc object.
        Args:
            character (str): Set of possible characters.
        """
        dict_character = list(character)

        # Create a dictionary to map characters to unique indices starting from 1, with 0 reserved for 'CTCblank' token
        self.dict = {}
        for i, char in enumerate(dict_character):
            self.dict[char] = i + 1

        # Character list including dummy '[CTCblank]' token for CTCLoss (index 0)
        self.character = ["[CTCblank]"] + dict_character

    def encode(self, text, batch_max_length=25):
        """
        Converts text-labels into concatenated text indices for CTCLoss calculation.
        Args:
            text (list): Text labels of each image. [batch_size]
            batch_max_length (int): Maximum length of text label in the batch. Default is 25.
        Returns:
            tuple: A tuple containing two TensorFlow tensors:
                - tf.Tensor: Concatenated text indices for CTCLoss. Shape: [sum(text_lengths)]
                - tf.Tensor: Length of each text. Shape: [batch_size]
        """
        length = [len(s) for s in text]
        text = "".join(text)
        text = [self.dict[char] for char in text]

        return tf.constant(text, dtype=tf.int32), tf.constant(length, dtype=tf.int32)

    def decode(self, text_index, length):
        """
        Converts text indices back into text labels.
        Args:
            text_index (tf.Tensor): Text indices for CTCLoss. Shape: [batch_size, batch_max_length]
            length (tf.Tensor): Length of each text. Shape: [batch_size]
        Returns:
            list: List of text labels for each entry in the batch.
        """
        texts = []
        index = 0
        for l in length:
            t = text_index[index: index + l]

            char_list = []
            for i in range(l):
                # Ignore padding (index 0) and repeated characters while reconstructing the text labels
                if t[i] != 0 and not (i > 0 and t[i - 1] == t[i]):
                    char_list.append(self.character[t[i]])
            text = "".join(char_list)

            texts.append(text)
            index += l
        return texts


class Averager(object):
    """
    Averager class computes the average for TensorFlow constants, typically used for averaging losses.
    """

    def __init__(self):
        """
        Initializes the Averager object.
        """
        self.reset()

    def reset(self):
        """
        Resets the internal counters for computing the average.
        """
        self.n_count = 0
        self.sum = 0

    def add(self, v):
        """
        Adds the given TensorFlow constant to the running sum.
        Args:
            v (tf.Tensor): TensorFlow constant to be added to the sum.
        """
        count = tf.size(v)
        v = tf.reduce_sum(v)
        self.n_count += count
        self.sum += v

    def val(self):
        """
        Computes the average value.
        Returns:
            float: Average value, or 0 if no values have been added.
        """
        res = 0
        if self.n_count != 0:
            res = self.sum / float(self.n_count)
        return res
