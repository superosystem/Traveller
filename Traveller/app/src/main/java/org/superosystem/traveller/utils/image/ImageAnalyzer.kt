package org.superosystem.traveller.utils.image

import android.content.Context
import android.graphics.Bitmap
import android.graphics.Rect
import androidx.camera.core.ImageAnalysis
import androidx.camera.core.ImageProxy
import androidx.lifecycle.MutableLiveData
import org.superosystem.traveller.data.model.ocr.IdentityCardImageCoordinate
import org.superosystem.traveller.data.model.ocr.JsonMemberClass
import org.superosystem.traveller.data.model.ocr.MarriedBoundingBoxCoordinate
import org.superosystem.traveller.data.model.ocr.NIKBoundingBoxCoordinate
import org.superosystem.traveller.data.model.ocr.NameBoundingBoxCoordinate
import org.superosystem.traveller.data.model.ocr.NationalityBoundingBoxCoordinate
import org.superosystem.traveller.data.model.ocr.SexBoundingBoxCoordinate
import org.tensorflow.lite.DataType
import org.tensorflow.lite.support.common.ops.NormalizeOp
import org.tensorflow.lite.support.image.ImageProcessor
import org.tensorflow.lite.support.image.TensorImage
import org.tensorflow.lite.support.image.ops.ResizeOp
import org.tensorflow.lite.task.core.BaseOptions
import org.tensorflow.lite.task.vision.detector.ObjectDetector
import java.io.File
import kotlin.math.floor

typealias ObjectDetectorCallback = (detectionResult: IdentityCardImageCoordinate?) -> Unit

class ImageAnalyzer(
    private val context: Context,
    private val imageFile: File,
    private val imageCropPercentages: MutableLiveData<Pair<Int, Int>>,
    private val setLoadingOCRScreenDialog: MutableLiveData<Boolean>,
    private val listener: ObjectDetectorCallback
) : ImageAnalysis.Analyzer {

    private var nikBoundingBoxCoordinate: NIKBoundingBoxCoordinate? = null
    private var nameBoundingBoxCoordinate: NameBoundingBoxCoordinate? = null
    private var sexBoundingBoxCoordinate: SexBoundingBoxCoordinate? = null
    private var marriedBoundingBoxCoordinate: MarriedBoundingBoxCoordinate? = null
    private var nationalityBoundingBoxCoordinate: NationalityBoundingBoxCoordinate? = null

    private lateinit var identityCardImageCoordinate: IdentityCardImageCoordinate

    @androidx.camera.core.ExperimentalGetImage
    override fun analyze(imageProxy: ImageProxy) {

        val mediaImage = imageProxy.image ?: return

        val rotationDegrees = imageProxy.imageInfo.rotationDegrees

        // We requested a setTargetAspectRatio, but it's not guaranteed that's what the camera
        // stack is able to support, so we calculate the actual ratio from the first frame to
        // know how to appropriately crop the image we want to analyze.
        val imageWidth = mediaImage.width
        val imageHeight = mediaImage.height

        val actualAspectRatio = imageWidth / imageHeight

        val convertImageToBitmap = convertYuv420888ImageToBitmap(mediaImage)
        val cropRect = Rect(0, 0, imageWidth, imageHeight)

        // If the image has a way wider aspect ratio than expected, crop less of the height so we
        // don't end up cropping too much of the image. If the image has a way taller aspect ratio
        // than expected, we don't have to make any changes to our cropping so we don't handle it
        // here.
        val currentCropPercentages = imageCropPercentages.value ?: return
        if (actualAspectRatio > 3) {
            val originalHeightCropPercentage = currentCropPercentages.first
            val originalWidthCropPercentage = currentCropPercentages.second
            imageCropPercentages.value =
                Pair(originalHeightCropPercentage / 2, originalWidthCropPercentage)
        }

        // If the image is rotated by 90 (or 270) degrees, swap height and width when calculating
        // the crop.
        val cropPercentages = imageCropPercentages.value ?: return
        val heightCropPercent = cropPercentages.first
        val widthCropPercent = cropPercentages.second
        val (widthCrop, heightCrop) = when (rotationDegrees) {
            90, 270 -> Pair(heightCropPercent / 100f, widthCropPercent / 100f)
            else -> Pair(widthCropPercent / 100f, heightCropPercent / 100f)
        }

        cropRect.inset(
            (imageWidth * widthCrop / 2).toInt(),
            (imageHeight * heightCrop / 2).toInt()
        )

        val croppedBitmap =
            rotateAndCrop(convertImageToBitmap, rotationDegrees, cropRect)

        runObjectDetection(croppedBitmap, croppedBitmap.width, croppedBitmap.height)

        imageProxy.close()

    }

    /**
     * runObjectDetection(bitmap: Bitmap)
     *      TFLite Object Detection function
     */
    private fun runObjectDetection(bitmap: Bitmap, bitmapWidth: Int, bitmapHeight: Int) {

        // Step 1: Create TFLite's TensorImage object

        val imageProcessor = ImageProcessor.Builder()
            .add(ResizeOp(300, 300, ResizeOp.ResizeMethod.BILINEAR))
            .add(NormalizeOp(0f, 1f))
            .build()

        var tensorImage = TensorImage(DataType.FLOAT32)
        tensorImage.load(bitmap)

        tensorImage = imageProcessor.process(tensorImage)

        // Step 2: Initialize the detector object

        val baseOptions = BaseOptions.builder().useNnapi().build()

        val options = ObjectDetector.ObjectDetectorOptions.builder()
            .setBaseOptions(baseOptions)
            .setMaxResults(5)
            .setScoreThreshold(0.7f)
            .build()

        val detector = ObjectDetector.createFromFileAndOptions(
            context,
            "final_model_metadata_v4.tflite",
            options
        )

        // Step 3: Feed given image to the detector
        val results = detector.detect(tensorImage)

        if (results.size == 5) {

            setLoadingOCRScreenDialog.postValue(true)

            results.forEach {

                when (it.categories.first().index) {
                    0 -> sexBoundingBoxCoordinate = SexBoundingBoxCoordinate(
                        xmin = floor((bitmapWidth * it.boundingBox.left) / 300).toInt() - 5,
                        ymin = floor((bitmapHeight * it.boundingBox.top) / 300).toInt() - 2,
                        xmax = floor((bitmapWidth * it.boundingBox.right) / 300).toInt() + 5,
                        ymax = floor((bitmapHeight * it.boundingBox.bottom) / 300).toInt() + 2
                    )

                    1 -> nationalityBoundingBoxCoordinate = NationalityBoundingBoxCoordinate(
                        xmin = floor((bitmapWidth * it.boundingBox.left) / 300).toInt() - 5,
                        ymin = floor((bitmapHeight * it.boundingBox.top) / 300).toInt() - 2,
                        xmax = floor((bitmapWidth * it.boundingBox.right) / 300).toInt() + 5,
                        ymax = floor((bitmapHeight * it.boundingBox.bottom) / 300).toInt() + 2
                    )

                    2 -> nameBoundingBoxCoordinate = NameBoundingBoxCoordinate(
                        xmin = floor((bitmapWidth * it.boundingBox.left) / 300).toInt() - 5,
                        ymin = floor((bitmapHeight * it.boundingBox.top) / 300).toInt() - 2,
                        xmax = floor((bitmapWidth * it.boundingBox.right) / 300).toInt() + 5,
                        ymax = floor((bitmapHeight * it.boundingBox.bottom) / 300).toInt() + 2
                    )

                    3 -> nikBoundingBoxCoordinate = NIKBoundingBoxCoordinate(
                        xmin = floor((bitmapWidth * it.boundingBox.left) / 300).toInt() - 5,
                        ymin = floor((bitmapHeight * it.boundingBox.top) / 300).toInt() - 2,
                        xmax = floor((bitmapWidth * it.boundingBox.right) / 300).toInt() + 5,
                        ymax = floor((bitmapHeight * it.boundingBox.bottom) / 300).toInt() + 2
                    )

                    4 -> marriedBoundingBoxCoordinate = MarriedBoundingBoxCoordinate(
                        xmin = floor((bitmapWidth * it.boundingBox.left) / 300).toInt() - 5,
                        ymin = floor((bitmapHeight * it.boundingBox.top) / 300).toInt() - 2,
                        xmax = floor((bitmapWidth * it.boundingBox.right) / 300).toInt() + 5,
                        ymax = floor((bitmapHeight * it.boundingBox.bottom) / 300).toInt() + 2
                    )
                }

            }

            if (nikBoundingBoxCoordinate != null &&
                nameBoundingBoxCoordinate != null &&
                sexBoundingBoxCoordinate != null &&
                marriedBoundingBoxCoordinate != null &&
                nationalityBoundingBoxCoordinate != null
            ) {

                bitmapToFile(imageFile, bitmap)

                identityCardImageCoordinate = IdentityCardImageCoordinate(
                    JsonMemberClass(
                        nikBoundingBoxCoordinate!!,
                        nameBoundingBoxCoordinate!!,
                        sexBoundingBoxCoordinate!!,
                        marriedBoundingBoxCoordinate!!,
                        nationalityBoundingBoxCoordinate!!
                    )
                )

                listener(identityCardImageCoordinate)

            }

        }

    }

}