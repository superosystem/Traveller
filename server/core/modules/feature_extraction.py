import tensorflow as tf
from tensorflow import keras
from tensorflow.keras import layers  # type: ignore


class VGG_FeatureExtractor(keras.models.Model):
    """FeatureExtractor of CRNN (https://arxiv.org/pdf/1507.05717.pdf)
    Berarti Transfer Learning untuk pretrainednya VGG

    """

    def __init__(self, output_channel=512, input_shape=(32, 100, 1), **kwargs):
        super(VGG_FeatureExtractor, self).__init__(**kwargs)

        self.output_channel = [
            int(output_channel / 8),
            int(output_channel / 4),
            int(output_channel / 2),
            output_channel,
        ]
        self.ConvNet = tf.keras.Sequential(
            [
                layers.Conv2D(
                    self.output_channel[0],
                    kernel_size=3,
                    strides=1,
                    padding="SAME",
                    kernel_initializer="he_uniform",
                    input_shape=input_shape,
                ),
                layers.ReLU(),
                layers.MaxPool2D(pool_size=2),
                layers.Conv2D(
                    self.output_channel[1],
                    kernel_size=3,
                    kernel_initializer="he_uniform",
                    strides=1,
                    padding="SAME",
                ),
                layers.ReLU(),
                layers.MaxPool2D(pool_size=2),
                layers.Conv2D(
                    self.output_channel[2],
                    kernel_size=3,
                    kernel_initializer="he_uniform",
                    strides=1,
                    padding="SAME",
                ),
                layers.ReLU(),
                layers.Conv2D(
                    self.output_channel[2],
                    kernel_size=3,
                    kernel_initializer="he_uniform",
                    strides=1,
                    padding="SAME",
                ),
                layers.ReLU(),
                layers.MaxPool2D(pool_size=(2, 1), strides=(2, 1)),
                layers.Conv2D(
                    self.output_channel[3],
                    kernel_size=3,
                    strides=1,
                    padding="SAME",
                    kernel_initializer="he_uniform",
                    use_bias=False,
                ),
                layers.BatchNormalization(),
                layers.ReLU(),
                layers.Conv2D(
                    self.output_channel[3],
                    kernel_size=3,
                    strides=1,
                    padding="SAME",
                    kernel_initializer="he_uniform",
                    use_bias=False,
                ),
                layers.BatchNormalization(),
                layers.ReLU(),
                layers.MaxPool2D(pool_size=(2, 1), strides=(2, 1)),
                layers.Conv2D(
                    self.output_channel[3],
                    kernel_size=2,
                    kernel_initializer="he_uniform",
                    strides=1,
                    padding="VALID",
                ),
                layers.ReLU(),
            ]
        )

    def call(self, X, training=None):
        return self.ConvNet(X)


class GRCL_unit(keras.models.Model):
    def __init__(self):
        self.BN_gfu = layers.BatchNormalization()
        self.BN_gfx = layers.BatchNormalization()
        self.BN_fu = layers.BatchNormalization()
        self.BN_rx = layers.BatchNormalization()
        self.BN_Gx = layers.BatchNormalization()

    def call(self, wgf_u, wgr_x, wf_u, wr_x):
        G_first_term = self.BN_gfu(wgf_u)
        G_second_term = self.BN_grx(wgr_x)
        G = keras.activations.sigmoid(G_first_term + G_second_term)

        x_first_term = self.BN_fu(wf_u)
        x_second_term = self.BN_Gx(self.BN_rx(wr_x) * G)
        x = keras.activations.ReLU(x_first_term + x_second_term)

        return x


class GRCL(keras.models.Model):
    """Gated RCNN"""

    def __init__(self, output_channel, num_iteration, kernel_size, padding):
        super().__init__()
        self.wgf_u = layers.Conv2D(
            output_channel, kernel_size=1, strides=1, padding="VALID"
        )
        self.wgr_x = layers.Conv2D(
            output_channel, kernel_size=1, strides=1, padding="VALID"
        )
        self.wf_u = layers.Conv2D(
            output_channel,
            kernel_size=kernel_size,
            strides=1,
            padding=padding,
            use_bias=False,
        )
        self.wr_x = layers.Conv2D(
            output_channel,
            kernel_size=kernel_size,
            strides=1,
            padding=padding,
            use_bias=False,
        )

        self.BN_x_init = layers.BatchNormalization()

        self.num_iteration = num_iteration
        self.GRCL = [GRCL_unit(output_channel) for _ in range(num_iteration)]
        self.GRCL = keras.Sequential(self.GRCL)

    def call(self, X):
        """The input of GRCL is consistant over time t, which is denoted by u(0)
        thus wgf_u / wf_u is also consistant over time t.
        """
        wgf_u = self.wgf_u(X)
        wf_u = self.wf_u(X)

        x = keras.activations.ReLU(self.BN_x_init(wf_u))

        for i in range(self.num_iteration):
            x = self.GRCL.get_layer(index=i)(wgf_u, self.wgr_x(x), wf_u, self.wr_x(x))

        return x


class RCNN_FeatureExtractor(keras.models.Model):
    """FeatureExtractor of GRCNN (https://papers.nips.cc/paper/6637-gated-recurrent-convolution-neural-network-for-ocr.pdf)"""

    def __init__(self, output_channel):
        super().__init__()
        self.output_channel = [
            int(output_channel / 8),
            int(output_channel / 4),
            int(output_channel / 2),
            output_channel,
        ]
        self.ConvNet = keras.Sequential(
            [
                layers.Conv2D(
                    self.output_channel[0], kernel_size=3, strides=1, padding="SAME"
                ),
                layers.ReLU(),
                layers.MaxPool2D(pool_size=2),
                GRCL(
                    self.output_channel[0],
                    num_iteration=5,
                    kernel_size=3,
                    padding="SAME",
                ),
                layers.MaxPool2D(pool_size=2),
                GRCL(
                    self.output_channel[1],
                    num_iteration=5,
                    kernel_size=3,
                    padding="SAME",
                ),
                layers.Lambda(
                    lambda x: tf.concat(
                        [
                            x,
                            tf.cast(
                                tf.zeros(shape=tf.concat([x.shape[:-1], [1]], axis=-1)),
                                tf.float32,
                            ),
                        ],
                        axis=-1,
                    )
                ),  # rawan error
                layers.MaxPool2D(pool_size=2),
            ]
        )

    def call(self, X):
        return self.ConvNet(X)


class BasisBlock(keras.models.Model):
    expansion = 1

    def __init__(self, inplanes, planes, strides=1, downsample=None) -> None:
        super().__init__()
        self.conv1 = self._conv3x3(planes)
        self.bn1 = layers.BatchNormalization()
        self.conv2 = self._conv3x3(planes)
        self.bn2 = layers.BatchNormalization()
        self.relu = layers.ReLU()
        self.add = layers.Add()
        self.downsample = downsample
        self.strides = strides

    @staticmethod
    def _conv3x3(planes, strides=1):
        "3x3 Convolution with padding"
        return layers.Conv2D(
            filters=planes,
            strides=strides,
            kernel_size=3,
            padding="SAME",
            use_bias=False,
        )

    def call(self, X):
        # rawan gagal juga

        residual = X

        out = self.conv1(X)
        out = self.bn1(out)
        out = self.relu(out)

        out = self.conv2(out)
        out = self.bn2(out)

        if self.downsample is not None:
            residual = self.downsample(X)

        out = self.add([out, residual])
        out = self.relu(out)

        return out


class ResNet(keras.models.Model):
    def __init__(self, output_channel, block, layers_) -> None:
        super().__init__()

        self.output_channel_block = [
            output_channel // 4,
            output_channel // 2,
            output_channel,
            output_channel,
        ]

        self.inplanes = output_channel // 8
        self.conv0_1 = layers.Conv2D(
            output_channel // 16,
            kernel_size=3,
            strides=1,
            padding="SAME",
            use_bias=False,
        )
        self.bn0_1 = layers.BatchNormalization()

        self.conv0_2 = layers.Conv2D(
            self.inplanes, kernel_size=3, strides=1, padding="SAME", use_bias=False
        )
        self.bn0_2 = layers.BatchNormalization()
        self.relu = layers.ReLU()

        self.maxpool1 = layers.MaxPool2D(pool_size=2, strides=2, padding="VALID")
        self.layer1 = self._make_layer(block, self.output_channel_block[0], layers_[0])
        self.conv1 = layers.Conv2D(
            self.output_channel_block[0],
            kernel_size=3,
            strides=1,
            padding="SAME",
            use_bias=False,
        )
        self.bn1 = layers.BatchNormalization()

        self.maxpool2 = layers.MaxPool2D(pool_size=2, strides=2, padding="VALID")
        self.layer2 = self._make_layer(
            block, self.output_channel_block[1], layers_[1], strides=1
        )
        self.conv2 = layers.Conv2D(
            self.output_channel_block[1],
            kernel_size=3,
            strides=1,
            padding="SAME",
            use_bias=False,
        )
        self.bn2 = layers.BatchNormalization()
        self.lambda_pad = layers.Concatenate()
        # rawan error
        self.maxpool3 = layers.MaxPool2D(pool_size=2, strides=(2, 1), padding="VALID")
        self.layer3 = self._make_layer(
            block, self.output_channel_block[2], layers_[2], strides=1
        )
        self.conv3 = layers.Conv2D(
            self.output_channel_block[2],
            kernel_size=3,
            strides=1,
            padding="SAME",
            use_bias=False,
        )
        self.bn3 = layers.BatchNormalization()

        self.layer4 = self._make_layer(
            block, self.output_channel_block[3], layers_[3], strides=1
        )
        # self.lambda_pad
        self.conv4_1 = layers.Conv2D(
            self.output_channel_block[3],
            kernel_size=2,
            strides=(2, 1),
            padding="SAME",
            use_bias=False,
        )
        self.bn4_1 = layers.BatchNormalization()

        self.conv4_2 = layers.Conv2D(
            self.output_channel_block[3],
            kernel_size=2,
            strides=(2, 1),
            padding="VALID",
            use_bias=False,
        )
        self.bn4_2 = layers.BatchNormalization()

    def _make_layer(self, block, planes, blocks, strides=1):
        downsample = None
        if strides != 1 or self.inplanes != planes * block.expansion:
            downsample = keras.Sequential(
                [
                    layers.Conv2D(
                        planes * block.expansion,
                        kernel_size=1,
                        strides=strides,
                        use_bias=False,
                    ),
                    layers.BatchNormalization(),
                ]
            )

        layers_ = []
        layers_.append(block(self.inplanes, planes, strides, downsample))

        self.inplanes = planes * block.expansion
        for i in range(1, blocks):
            layers_.append(block(self.inplanes, planes))

        return keras.Sequential(layers_)

    def call(self, X):

        X = self.conv0_1(X)
        X = self.bn0_1(X)
        X = self.relu(X)
        X = self.conv0_2(X)
        X = self.bn0_2(X)
        X = self.relu(X)

        X = self.maxpool1(X)
        X = self.layer1(X)
        X = self.conv1(X)
        X = self.bn1(X)
        X = self.relu(X)

        X = self.maxpool2(X)
        X = self.layer2(X)
        X = self.conv2(X)
        X = self.bn2(X)
        X = self.relu(X)
        # # shape hardcoded
        # X = self.lambda_pad([X,
        #                 tf.zeros(shape=(1, 8, 25, 1))])
        X = self.maxpool3(X)
        X = self.layer3(X)
        X = self.conv3(X)
        X = self.bn3(X)
        X = self.relu(X)

        X = self.layer4(X)
        X = self.conv4_1(X)
        X = self.bn4_1(X)
        X = self.relu(X)
        X = self.conv4_2(X)
        X = self.bn4_2(X)
        X = self.relu(X)

        return X


class ResNet_FeatureExtractor(keras.models.Model):
    """FeatureExtractor of FAN (http://openaccess.thecvf.com/content_ICCV_2017/papers/Cheng_Focusing_Attention_Towards_ICCV_2017_paper.pdf)"""

    def __init__(self, output_channel=512) -> None:
        super().__init__()
        self.ConvNet = ResNet(output_channel, BasisBlock, [1, 2, 5, 3])

    def call(self, X) -> tf.Tensor:
        return self.ConvNet(X)
