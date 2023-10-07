import tensorflow as tf
from tensorflow import keras
from tensorflow.keras import layers


class BidirectionalLSTM(keras.models.Model):
    def __init__(self, hidden_size, output_size):
        super().__init__()
        # Initialize Bidirectional LSTM layer and a Dense layer
        self.rnn = layers.Bidirectional(layers.LSTM(hidden_size, time_major=False))
        self.linear = layers.Dense(output_size)

    def call(self, X):
        """
        input : visual feature [batch_size x T x input_size]
        output : contextual feature [batch_size x T x output_size]
        """

        # Apply the Bidirectional LSTM layer to the input
        recurrent = self.rnn(X)
        # Apply the Dense layer to the output of the LSTM
        output = self.linear(recurrent)
        return output
