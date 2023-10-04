import tensorflow as tf
from tensorflow import keras
from tensorflow.keras import layers


class AttentionCell(keras.models.Model):
    def __init__(self, hidden_size, num_embeddings):
        super().__init__()
        self.i2h = layers.Dense(hidden_size, use_bias=False)
        self.h2h = layers.Dense(hidden_size)
        self.score = layers.Dense(1, use_bias=False)
        self.rnn = layers.LSTMCell()
        self.hidden_size = hidden_size

    def call(self, prev_hidden, batch_H, char_onehots):
        # [batch_size x num_encoder_step x num_channel] -> [batch_size x num_encoder_step x hidden_size]
        batch_H_proj = self.i2h(batch_H)
        prev_hidden_proj = tf.expand_dims(self.h2h(prev_hidden[0]), axis=1)
        e = self.score(
            tf.math.tanh(tf.add(batch_H_proj + prev_hidden_proj))
        )  # rawan error

        alpha = tf.math.softmax(e, axis=1)
        context = tf.squeeze(tf.matmul(tf.transpose(alpha, [0, 2, 1]), batch_H), axis=1)
        concat_context = tf.concat([context, char_onehots], axis=1)
        cur_hidden = self.rnn(concat_context, prev_hidden)  # rawan error
        return cur_hidden, alpha


class Attention(keras.models.Model):
    def __init__(self, hidden_size, num_classes):
        super().__init__()
        self.attention_cell = AttentionCell()
        self.hidden_size = hidden_size
        self.num_classes = num_classes
        self.generator = layers.Dense(num_classes)

    @staticmethod
    def _char_to_onehot(input_char, onehot_dim=38):
        input_char = tf.expand_dims(input_char, axis=1)
        one_hot = tf.one_hot(input_char, onehot_dim)
        return one_hot  # rawan error

    def call(self, batch_H, text, training=True, batch_max_length=25):
        """
        input:
            batch_H : contextual_feature H = hidden state of encoder. [batch_size x num_steps x contextual_feature_channels]
            text : the text-index of each image. [batch_size x (max_length+1)]. +1 for [GO] token. text[:, 0] = [GO].
        output: probability distribution at each step [batch_size x num_steps x num_classes]
        """
        batch_size = batch_H.shape[0]
        num_steps = batch_max_length + 1

        output_hiddens = tf.zeros(shape=[batch_size, num_steps, self.hidden_size])
        hidden = (
            tf.zeros(shapes=[batch_size, self.hidden_size]),
            tf.zeros(shapes=[batch_size, self.hidden_size]),
        )

        if training:
            for i in range(num_steps):
                # one-hot vectors for a i-th char. in a batch
                char_onehots = self._char_to_onehot(
                    text[:, i], onehot_dim=self.num_classes
                )
                # hidden : decoder's hidden s_{t-1}, batch_H : encoder's hidden H, char_onehots : one-hot(y_{t-1})
                hidden, alpha = self.attention_cell(hidden, batch_H, char_onehots)
                output_hiddens[:, i, :] = hidden[0]

            probs = self.generator(output_hiddens)

        else:
            targets = tf.zeros(shape=[batch_size], dtypes=tf.float64)
            probs = tf.zeros(shape=[batch_size, num_steps, self.num_classes])

            for i in range(num_steps):
                char_onehots = self._char_to_onehot(
                    targets, onehot_dim=self.num_classes
                )
                hidden, alpha = self.attention_cell(hidden, batch_H, char_onehots)
                probs_step = self.generator(hidden[0])
                probs[:, i, :] = probs_step
                _, next_input = tf.reduce_max(probs_step, axis=1)  # rawan error
                targets = next_input

        return probs
