import tensorflow as tf
import tensorflow_addons as tfa
from tensorflow import keras
from tensorflow.keras import layers

from modules.feature_extraction import ResNet_FeatureExtractor, VGG_FeatureExtractor


class Model(keras.models.Model):
    def __init__(self, opt, **kwargs):
        super(Model, self).__init__(**kwargs)
        self.opt = opt

        print("No Transformation Module")

        if opt.FeatureExtraction == "VGG":
            if opt.pretrained:
                self.FeatureExtraction = keras.applications.vgg16.VGG16(
                    include_top=False, weights="imagenet", input_shape=(opt.imgH, opt.imgW, 3)
                )
                # for layer in self.FeatureExtraction.layers:
                #     layer.trainable = False
            else:
                self.FeatureExtraction = VGG_FeatureExtractor(opt.output_channel, input_shape=(opt.imgH, opt.imgW, 1))
        elif opt.FeatureExtraction == "ResNet":
            self.FeatureExtraction = ResNet_FeatureExtractor(opt.output_channel)
        self.FeatureExtraction_output = opt.output_channel
        # untuk sekarang
        if opt.FeatureExtraction == "VGG":
            if opt.pretrained:
                self.AdaptiveAvgPool = tfa.layers.AdaptiveAveragePooling2D(
                    output_size=(3, 1)
                )
            else:
                self.AdaptiveAvgPool = tfa.layers.AdaptiveAveragePooling2D(
                    output_size=(49, 1)
                )
        elif opt.FeatureExtraction == "ResNet":
            self.AdaptiveAvgPool = tfa.layers.AdaptiveAveragePooling2D(
                output_size=(23, 1)
            )
        print("No sequence modelling module specified")
        self.SequenModelling_output = self.FeatureExtraction_output
        self.flatten = layers.Flatten()
        self.Prediction = layers.Dense(opt.num_class)

    def call(self, X, text, training=None):

        """Feature Extraction Stage"""
        visual_feature = self.FeatureExtraction(X)
        visual_feature = self.AdaptiveAvgPool(
            tf.transpose(visual_feature, perm=[0, 2, 1, 3])
        )  # [b, w, h, c] -> [b, h, w, c]
        visual_feature = tf.squeeze(visual_feature, axis=2)

        # if self.opt.pretrained:
        #     contextual_feature = self.flatten(visual_feature)
        # else:
        contextual_feature = visual_feature
        prediction = self.Prediction(contextual_feature)
        return prediction
