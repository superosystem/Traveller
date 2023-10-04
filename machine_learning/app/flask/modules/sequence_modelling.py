import tensorflow as tf
from tensorflow import keras
from tensorflow.keras import layers


class BidirectionalLSTM(keras.models.Model):
    def __init__(self, hidden_size, output_size):
        super().__init__()
        # kalau di implementasinya time_major=False == batch_first=True
        self.rnn = layers.Bidirectional(layers.LSTM(hidden_size, time_major=False))
        self.linear = layers.Dense(output_size)

    def call(self, X):
        """
        input : visual feature [batch_size x T x input_size]
        output : contextual feature [batch_size x T x output_size]
        """

        # kalo di pytorch harusnya disini di flatten_parameters() dulu
        reccurent = self.rnn(X)
        # rawan error karna beda behavior antara LSTM Tensorflow dengan LSTM pytorch
        output = self.linear(reccurent)
        return output
