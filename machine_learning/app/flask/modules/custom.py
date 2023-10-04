import warnings

import tensorflow as tf
from tensorflow import keras


class CTCLossLayer(keras.layers.Layer):
    def call(self, inputs):
        labels = inputs[0]
        logits = inputs[1]
        label_len = inputs[2]
        logit_len = inputs[3]

        logits_trans = tf.transpose(logits, (1, 0, 2))
        label_len = tf.reshape(label_len, (-1,))
        logit_len = tf.reshape(logit_len, (-1,))
        loss = tf.reduce_mean(
            tf.nn.ctc_loss(labels, logits_trans, label_len, logit_len, blank_index=-1)
        )
        # define loss here instead of compile()
        self.add_loss(loss)

        # decode
        decoded, _ = tf.nn.ctc_greedy_decoder(logits_trans, logit_len)

        # Inaccuracy: label error rate
        ler = tf.reduce_mean(tf.edit_distance(tf.cast(decoded[0], tf.int32), labels))
        self.add_metric(ler, name="ler", aggregation="mean")

        return logits  # Pass-through layer


def custom_sparse_categorical_crossentropy(
    y_true: tf.Tensor,
    y_pred: tf.Tensor,
    from_logits: bool = False,
    ignore_index: int = -1,
    axis: int = -1,
):
    y_true = tf.convert_to_tensor(y_true)
    y_pred = tf.convert_to_tensor(y_pred)
    y_pred, from_logits = to_logits(y_pred, from_logits)

    y_true = squeeze(y_true)
    y_pred = squeeze(y_pred, dims=3)

    valid_mask = y_true != ignore_index
    indices = tf.where(valid_mask)

    ce = tf.losses.sparse_categorical_crossentropy(
        y_true[valid_mask], y_pred[valid_mask], from_logits, axis
    )
    ce = tf.scatter_nd(indices, ce, tf.cast(tf.shape(y_true), tf.int64))

    ce = tf.math.divide_no_nan(
        tf.reduce_sum(ce, axis=-1),
        tf.cast(tf.math.count_nonzero(valid_mask, axis=-1), ce.dtype),
    )

    return ce


def squeeze(y, dims: int = 2):
    if dims not in (2, 3):
        raise ValueError(
            f"Illegal value for parameter dims=`{dims}`. Can only squeeze "
            "positional signal, resulting in a tensor with rank 2 or 3."
        )
    shape = tf.shape(y)
    new_shape = [shape[0], -1]
    if dims == 3:  # keep channels.
        new_shape += [shape[-1]]
    return tf.reshape(y, new_shape)


def to_logits(output, from_logits: bool = False):
    if from_logits:
        return output, True

    if hasattr(output, "_keras_logits"):
        if from_logits:
            warnings.warn(
                '"`dig_logits_if_available` received `from_logits=True`, but '
                "the `output` argument was produced by a sigmoid or softmax "
                'activation and thus does not represent logits. Was this intended?"',
                stacklevel=2,
            )
        return output._keras_logits, True

    if (
        not isinstance(output, (tf.__internal__.EagerTensor, tf.Variable))
        and output.op.type in ("Softmax", "Sigmoid")
    ) and not hasattr(output, "_keras_history"):
        assert len(output.op.inputs) == 1
        return output.op.inputs[0], True

    return output, False
