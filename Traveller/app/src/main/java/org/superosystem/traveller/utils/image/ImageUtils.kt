package org.superosystem.traveller.utils.image

import android.content.ContentResolver
import android.content.Context
import android.graphics.Bitmap
import android.graphics.ImageFormat
import android.graphics.Matrix
import android.graphics.Rect
import android.media.Image
import android.net.Uri
import android.os.Environment
import androidx.annotation.ColorInt
import java.io.File
import java.io.FileOutputStream
import java.io.InputStream
import java.io.OutputStream
import java.text.SimpleDateFormat
import java.util.*

private const val FILENAME_FORMAT = "s m k dd-MMM-yyyy"

val timeStamp: String = SimpleDateFormat(
    FILENAME_FORMAT,
    Locale.US
).format(System.currentTimeMillis())

fun createTempFile(context: Context): File {
    val storageDir: File? = context.getExternalFilesDir(Environment.DIRECTORY_PICTURES)
    storageDir?.deleteOnExit()
    return File.createTempFile(timeStamp, ".jpg", storageDir)
}

fun createCustomTempFile(context: Context): File {
    val storageDir: File? = context.getExternalFilesDir(Environment.DIRECTORY_PICTURES)
    return File.createTempFile(
        timeStamp,
        ".jpg",
        storageDir
    )
}

fun uriToFile(selectedImg: Uri, context: Context): File {
    val contentResolver: ContentResolver = context.contentResolver
    val myFile = createCustomTempFile(context)

    val inputStream = contentResolver.openInputStream(selectedImg) as InputStream
    val outputStream: OutputStream = FileOutputStream(myFile)
    val buf = ByteArray(1024)
    var len: Int
    while (inputStream.read(buf).also { len = it } > 0) outputStream.write(buf, 0, len)
    outputStream.close()
    inputStream.close()

    return myFile
}

fun bitmapToFile(file: File, pictureBitmap: Bitmap): File {

    val fOut = FileOutputStream(file)

    pictureBitmap.compress(
        Bitmap.CompressFormat.JPEG,
        100,
        fOut
    ) // saving the Bitmap to a file compressed as a JPEG with 85% compression rate

    fOut.flush() // Not really required

    fOut.close() // do not forget to close the stream

    return file

}

fun convertYuv420888ImageToBitmap(image: Image): Bitmap {
    require(image.format == ImageFormat.YUV_420_888) {
        "Unsupported image format $(image.format)"
    }

    val planes = image.planes

    // Because of the variable row stride it's not possible to know in
    // advance the actual necessary dimensions of the yuv planes.
    val yuvBytes = planes.map { plane ->
        val buffer = plane.buffer
        val yuvBytes = ByteArray(buffer.capacity())
        buffer[yuvBytes]
        buffer.rewind()  // Be kind…
        yuvBytes
    }

    val yRowStride = planes[0].rowStride
    val uvRowStride = planes[1].rowStride
    val uvPixelStride = planes[1].pixelStride
    val width = image.width
    val height = image.height
    @ColorInt val argb8888 = IntArray(width * height)
    var i = 0
    for (y in 0 until height) {
        val pY = yRowStride * y
        val uvRowStart = uvRowStride * (y shr 1)
        for (x in 0 until width) {
            val uvOffset = (x shr 1) * uvPixelStride
            argb8888[i++] =
                yuvToRgb(
                    yuvBytes[0][pY + x].toIntUnsigned(),
                    yuvBytes[1][uvRowStart + uvOffset].toIntUnsigned(),
                    yuvBytes[2][uvRowStart + uvOffset].toIntUnsigned()
                )
        }
    }
    val bitmap = Bitmap.createBitmap(width, height, Bitmap.Config.ARGB_8888)
    bitmap.setPixels(argb8888, 0, width, 0, 0, width, height)
    return bitmap
}

fun rotateAndCrop(
    bitmap: Bitmap,
    imageRotationDegrees: Int,
    cropRect: Rect
): Bitmap {
    val matrix = Matrix()
    matrix.preRotate(imageRotationDegrees.toFloat())
    return Bitmap.createBitmap(
        bitmap,
        cropRect.left,
        cropRect.top,
        cropRect.width(),
        cropRect.height(),
        matrix,
        true
    )
}

@ColorInt
private fun yuvToRgb(nY: Int, nU: Int, nV: Int): Int {

    val CHANNEL_RANGE = 0 until (1 shl 18)

    var nY = nY
    var nU = nU
    var nV = nV
    nY -= 16
    nU -= 128
    nV -= 128
    nY = nY.coerceAtLeast(0)

    // This is the floating point equivalent. We do the conversion in integer
    // because some Android devices do not have floating point in hardware.
    // nR = (int)(1.164 * nY + 2.018 * nU);
    // nG = (int)(1.164 * nY - 0.813 * nV - 0.391 * nU);
    // nB = (int)(1.164 * nY + 1.596 * nV);
    var nR = 1192 * nY + 1634 * nV
    var nG = 1192 * nY - 833 * nV - 400 * nU
    var nB = 1192 * nY + 2066 * nU

    // Clamp the values before normalizing them to 8 bits.
    nR = nR.coerceIn(CHANNEL_RANGE) shr 10 and 0xff
    nG = nG.coerceIn(CHANNEL_RANGE) shr 10 and 0xff
    nB = nB.coerceIn(CHANNEL_RANGE) shr 10 and 0xff
    return -0x1000000 or (nR shl 16) or (nG shl 8) or nB
}

private fun Byte.toIntUnsigned(): Int {
    return toInt() and 0xFF
}