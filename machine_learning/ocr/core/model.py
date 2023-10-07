import tensorflow as tf
import tensorflow_addons as tfa
from tensorflow import keras
from tensorflow.keras import layers

from modules.feature_extraction import ResNet_FeatureExtractor, VGG_FeatureExtractor


class Model(keras.models.Model):
    """
    Custom model class for text recognition, allowing different feature extraction methods and sequence modeling.
    """

    def __init__(self, opt, **kwargs):
        """
        Initializes the Model object.
        Args:
            opt (object): An object containing various configuration options.
        """
        super(Model, self).__init__(**kwargs)
        self.opt = opt

        # Select feature extraction method based on the specified configuration
        if opt.FeatureExtraction == "VGG":
            if opt.pretrained:
                self.FeatureExtraction = keras.applications.vgg16.VGG16(
                    include_top=False, weights="imagenet", input_shape=(opt.imgH, opt.imgW, 3)
                )
            else:
                self.FeatureExtraction = VGG_FeatureExtractor(opt.output_channel, input_shape=(opt.imgH, opt.imgW, 1))
        elif opt.FeatureExtraction == "ResNet":
            self.FeatureExtraction = ResNet_FeatureExtractor(opt.output_channel)
        self.FeatureExtraction_output = opt.output_channel

        # Configure adaptive average pooling based on feature extraction method
        if opt.FeatureExtraction == "VGG":
            self.AdaptiveAvgPool = tfa.layers.AdaptiveAveragePooling2D(
                output_size=(3, 1) if opt.pretrained else (49, 1)
            )
        elif opt.FeatureExtraction == "ResNet":
            self.AdaptiveAvgPool = tfa.layers.AdaptiveAveragePooling2D(
                output_size=(23, 1)
            )

        # Configure sequence modeling components (not specified in the provided code)
        print("No sequence modelling module specified")
        self.SequenModelling_output = self.FeatureExtraction_output
        self.flatten = layers.Flatten()

        # Dense layer for final prediction
        self.Prediction = layers.Dense(opt.num_class)

    def call(self, X, text, training=None):
        """
        Defines the forward pass of the model.
        Args:
            X (tf.Tensor): Input images.
            text (tf.Tensor): Text data.
            training (bool): Indicates whether the model is in training mode.
        Returns:
            tf.Tensor: Predicted output.
        """
        # Feature Extraction Stage
        visual_feature = self.FeatureExtraction(X)
        visual_feature = self.AdaptiveAvgPool(
            tf.transpose(visual_feature, perm=[0, 2, 1, 3])
        )  # [b, w, h, c] -> [b, h, w, c]
        visual_feature = tf.squeeze(visual_feature, axis=2)

        # Sequence Modelling (not specified in the provided code)

        # Contextual Feature
        contextual_feature = visual_feature

        # Final Prediction
        prediction = self.Prediction(contextual_feature)
        return prediction
