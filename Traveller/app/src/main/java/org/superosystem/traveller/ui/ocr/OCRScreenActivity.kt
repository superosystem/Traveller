package org.superosystem.traveller.ui.ocr

import LoadingOCRScreenDialog
import android.Manifest
import android.content.Intent
import android.content.pm.PackageManager
import android.graphics.Color
import android.graphics.Paint
import android.graphics.PixelFormat
import android.graphics.PorterDuff
import android.graphics.PorterDuffXfermode
import android.graphics.RectF
import android.os.Bundle
import android.os.Handler
import android.os.Looper
import android.util.DisplayMetrics
import android.util.Log
import android.view.SurfaceHolder
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.camera.core.AspectRatio
import androidx.camera.core.CameraControl
import androidx.camera.core.CameraSelector
import androidx.camera.core.ImageAnalysis
import androidx.camera.core.Preview
import androidx.camera.lifecycle.ProcessCameraProvider
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import androidx.lifecycle.LifecycleOwner
import androidx.lifecycle.ViewModelProvider
import com.google.gson.Gson
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.MultipartBody
import okhttp3.RequestBody
import okhttp3.RequestBody.Companion.asRequestBody
import okhttp3.RequestBody.Companion.toRequestBody
import org.superosystem.traveller.R
import org.superosystem.traveller.data.preference.SavedPreference
import org.superosystem.traveller.data.repository.AuthRepository
import org.superosystem.traveller.data.repository.OCRRepository
import org.superosystem.traveller.databinding.ActivityOcrscreenBinding
import org.superosystem.traveller.utils.Constants
import org.superosystem.traveller.utils.Resources
import org.superosystem.traveller.utils.image.ImageAnalyzer
import org.superosystem.traveller.utils.image.createTempFile
import org.superosystem.traveller.viewmodel.AuthViewModel
import org.superosystem.traveller.viewmodel.OCRScreenViewModel
import org.superosystem.traveller.viewmodel.factory.AuthViewModelFactory
import org.superosystem.traveller.viewmodel.factory.OCRScreenViewModelFactory
import java.io.File
import java.util.concurrent.ExecutorService
import java.util.concurrent.Executors
import kotlin.math.abs
import kotlin.math.ln
import kotlin.math.max
import kotlin.math.min


class OCRScreenActivity : AppCompatActivity() {
    private lateinit var binding: ActivityOcrscreenBinding

    // Session
    private lateinit var savedPref: SavedPreference
    private lateinit var accessToken: String
    private lateinit var bookingID: String
    private lateinit var data: RequestBody

    // image controller stuff
    private lateinit var cameraController: CameraControl
    private var toggleTorch: Boolean = false
    private var lensFacing: Int = CameraSelector.LENS_FACING_BACK

    // photo file stuff
    private lateinit var photoFile: File
    private lateinit var currentPhotoFilePath: String

    // image analyzer stuff
    private var imageAnalyzer: ImageAnalysis? = null
    private lateinit var cameraAnalysisExecutor: ExecutorService

    // dialog stuff
    private lateinit var loadingDialog: LoadingOCRScreenDialog

    // ViewModel
    private lateinit var viewModel: OCRScreenViewModel
    private lateinit var authViewModel: AuthViewModel

    companion object {
        // We only need to analyze and capture a section, so we set crop percentages
        // to avoid analyze the entire image from the live camera feed.
        const val DESIRED_WIDTH_CROP_PERCENT = 8
        const val DESIRED_HEIGHT_CROP_PERCENT = 60

        private const val RATIO_4_3_VALUE = 4.0 / 3.0
        private const val RATIO_16_9_VALUE = 16.0 / 9.0

        private val REQUIRED_PERMISSIONS = arrayOf(Manifest.permission.CAMERA)
        private const val REQUEST_CODE_PERMISSIONS = 10

        private const val TAG = "OCRScreenActivity"

    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityOcrscreenBinding.inflate(layoutInflater)
        setContentView(binding.root)

        supportActionBar?.hide()

        loadingDialog = LoadingOCRScreenDialog(this)

        setupBookingIDAndToken()

        // instantiate View Model
        val authFactory = AuthViewModelFactory(AuthRepository())
        authViewModel = ViewModelProvider(this, authFactory)[AuthViewModel::class.java]

        val factory = OCRScreenViewModelFactory(OCRRepository())
        viewModel = ViewModelProvider(this, factory)[OCRScreenViewModel::class.java]

        viewModel.setLoadingOCRScreenDialog.observe(this) {
            if (it == true) {
                loadingDialog.startLoadingDialog()
            } else if (it == false) {
                loadingDialog.dismissLoadingDialog()
            }
        }

        cameraAnalysisExecutor = Executors.newSingleThreadExecutor()

        photoFile = createTempFile(this).also {
            currentPhotoFilePath = it.absolutePath
        }

        binding.tvBackButton.setOnClickListener {
            finish()
        }

    }

    override fun onResume() {

        super.onResume()

        if (!allPermissionGranted()) {
            ActivityCompat.requestPermissions(
                this,
                REQUIRED_PERMISSIONS,
                REQUEST_CODE_PERMISSIONS
            )
        } else {

            binding.viewFinder.post {

                binding.overlay.apply {

                    setZOrderOnTop(true)
                    holder.setFormat(PixelFormat.TRANSPARENT)
                    holder.addCallback(object : SurfaceHolder.Callback {

                        override fun surfaceCreated(p0: SurfaceHolder) {
                            holder?.let {
                                drawOverlay(
                                    it,
                                    DESIRED_HEIGHT_CROP_PERCENT,
                                    DESIRED_WIDTH_CROP_PERCENT
                                )
                            }
                        }

                        override fun surfaceChanged(p0: SurfaceHolder, p1: Int, p2: Int, p3: Int) {

                        }

                        override fun surfaceDestroyed(p0: SurfaceHolder) {

                        }

                    })
                }

                setUpCameraAndCameraUtilityWhenCameraPermissionGranted()

            }

        }

    }

    override fun onBackPressed() {
        finish()
    }

    override fun onDestroy() {
        super.onDestroy()
        cameraAnalysisExecutor.shutdown()
    }

    private fun setUpCameraAndCameraUtilityWhenCameraPermissionGranted() {

        bindCameraUseCases()

        binding.switchTorch.setOnClickListener {

            if (toggleTorch) {

                toggleTorch = false
                binding.switchTorch.setImageResource(R.drawable.ic_flash_off)

            } else {

                toggleTorch = true
                binding.switchTorch.setImageResource(R.drawable.ic_flash_on)
            }

            cameraController.enableTorch(toggleTorch)

        }

        binding.switchCamera.setOnClickListener {

            lensFacing =
                if (lensFacing.equals(CameraSelector.LENS_FACING_BACK)) CameraSelector.LENS_FACING_FRONT
                else CameraSelector.LENS_FACING_BACK

            bindCameraUseCases()

        }

    }

    private fun setupBookingIDAndToken() {

        savedPref = SavedPreference(this)
        val tokenFromAPI = (savedPref.getData(Constants.ACCESS_TOKEN))!!
        accessToken = "Bearer $tokenFromAPI"

        bookingID = intent.extras?.getString("id").toString()

    }

    private fun bindCameraUseCases() {

        val cameraProviderFuture = ProcessCameraProvider.getInstance(this)
        cameraProviderFuture.addListener({

            // Camera provider is now guaranteed to be available
            val cameraProvider = cameraProviderFuture.get()

            // Get screen metrics used to setup camera for full screen resolution
            val metrics = DisplayMetrics().also { binding.viewFinder.display.getRealMetrics(it) }
            Log.d(TAG, "Screen metrics: ${metrics.widthPixels} x ${metrics.heightPixels}")

            val screenAspectRatio = aspectRatio(metrics.widthPixels, metrics.heightPixels)
            Log.d(TAG, "Preview aspect ratio: $screenAspectRatio")

            // Set up the view finder use case to display camera preview
            val preview = Preview.Builder()
                .setTargetAspectRatio(screenAspectRatio)
                .setTargetRotation(binding.viewFinder.display.rotation)
                .build()
                .also {
                    // Use the camera object to link our preview use case with the view
                    it.setSurfaceProvider(binding.viewFinder.surfaceProvider)
                }

            viewModel.imageCropPercentages.observe(this) {
                drawOverlay(
                    binding.overlay.holder,
                    it.first,
                    it.second
                )
            }

            // Create a new camera selector each time, enforcing lens facing
            val cameraSelector = CameraSelector.Builder().requireLensFacing(lensFacing).build()

            // Build the image analysis use case and instantiate our analyzer
            imageAnalyzer = ImageAnalysis.Builder()
                // We request aspect ratio but no resolution
                .setTargetAspectRatio(screenAspectRatio)
                .setTargetRotation(binding.viewFinder.display.rotation)
                .setBackpressureStrategy(ImageAnalysis.STRATEGY_KEEP_ONLY_LATEST)
                .build()
                .also {
                    it.setAnalyzer(
                        cameraAnalysisExecutor, ImageAnalyzer(
                            this,
                            photoFile,
                            viewModel.imageCropPercentages,
                            viewModel.setLoadingOCRScreenDialog
                        ) { detectedObjectList ->

                            it.clearAnalyzer()

                            data = Gson().toJson(detectedObjectList)
                                .toRequestBody("text/plain".toMediaType())

                            observerOCRScanning(accessToken, bookingID, photoFile, data)
                        }
                    )
                }

            try {
                // Apply declared configs to CameraX using the same lifecycle owner
                cameraProvider.unbindAll()

                val camera = cameraProvider.bindToLifecycle(
                    this as LifecycleOwner, cameraSelector, preview, imageAnalyzer
                )

                cameraController = camera.cameraControl
                cameraController.enableTorch(false)

            } catch (exc: Exception) {
                Toast.makeText(
                    this,
                    "Unable to show camera",
                    Toast.LENGTH_SHORT
                ).show()
            }

        }, ContextCompat.getMainExecutor(this))

    }

    private fun observerOCRScanning(
        accessToken: String,
        bookingID: String,
        photoFile: File,
        data: RequestBody
    ) {

        val requestImageFile = photoFile.asRequestBody("image/jpeg".toMediaTypeOrNull())
        val imageMultipart: MultipartBody.Part = MultipartBody.Part.createFormData(
            "file",
            photoFile.name,
            requestImageFile
        )

        Handler(Looper.getMainLooper()).post {
            viewModel.scanIDCard(accessToken, imageMultipart, data).observe(this) { response ->

                if (response is Resources.Error) {
                    viewModel.setLoadingOCRScreenDialog.value = false
                    Toast.makeText(this, response.error, Toast.LENGTH_SHORT).show()
                } else if (response is Resources.Success) {
                    val result = response.data
                    if (result != null) {

                        if (result.status == "Success") {

                            observerOCRResult(accessToken, bookingID)

                        } else {

                            val dataToken = hashMapOf(
                                "refreshToken" to savedPref.getData(Constants.REFRESH_TOKEN)
                            )

                            observeUpdateTokenForObserverOCRScanning(dataToken)
                        }

                    } else {
                        viewModel.setLoadingOCRScreenDialog.value = false
                        Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                    }
                }

            }
        }

    }

    private fun observerOCRResult(accessToken: String, bookingID: String) {

        viewModel.retrieveIDCardResult(accessToken).observe(this) { response ->

            if (response is Resources.Error) {
                viewModel.setLoadingOCRScreenDialog.value = false
                Toast.makeText(this, response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                val result = response.data
                if (result != null) {
                    if (result.status == "success") {
                        val intent = Intent(this, OCRResultActivity::class.java)
                        intent.putExtra("IDCardResult", result.data)
                        intent.putExtra("BookingID", bookingID)

                        viewModel.setLoadingOCRScreenDialog.value = false

                        startActivity(intent)
                        finish()

                    } else {
                        val dataToken = hashMapOf(
                            "refreshToken" to savedPref.getData(Constants.REFRESH_TOKEN)
                        )
                        observeUpdateTokenForOCRResult(dataToken)
                    }
                } else {
                    viewModel.setLoadingOCRScreenDialog.value = false
                    Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                }
            }

        }

    }

    private fun observeUpdateTokenForObserverOCRScanning(dataToken: HashMap<String, String?>) {
        authViewModel.updateToken(dataToken).observe(this) { response ->
            if (response is Resources.Error) {
                viewModel.setLoadingOCRScreenDialog.value = false
                Toast.makeText(this, response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                val result = response.data
                if (result != null) {
                    if (result.status.equals("success")) {

                        val newAccessToken = result.data?.accessToken.toString()
                        //save new token
                        savedPref.putData(Constants.ACCESS_TOKEN, newAccessToken)

                        //get new token
                        val tokenFromAPI = (savedPref.getData(Constants.ACCESS_TOKEN))
                        val accessToken = "Bearer $tokenFromAPI"


                        observerOCRScanning(accessToken, bookingID, photoFile, data)
                    } else {
                        Log.d("REGIS", result.status.toString())
                    }
                } else {
                    viewModel.setLoadingOCRScreenDialog.value = false
                    Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    private fun observeUpdateTokenForOCRResult(dataToken: HashMap<String, String?>) {
        authViewModel.updateToken(dataToken).observe(this) { response ->
            if (response is Resources.Error) {
                viewModel.setLoadingOCRScreenDialog.value = false
                Toast.makeText(this, response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                val result = response.data
                if (result != null) {
                    if (result.status.equals("success")) {

                        val newAccessToken = result.data?.accessToken.toString()
                        //save new token
                        savedPref.putData(Constants.ACCESS_TOKEN, newAccessToken)

                        //get new token
                        val tokenFromAPI = (savedPref.getData(Constants.ACCESS_TOKEN))
                        val accessToken = "Bearer $tokenFromAPI"


                        observerOCRResult(accessToken, bookingID)
                    } else {
                        Log.d("REGIS", result.status.toString())
                    }
                } else {
                    viewModel.setLoadingOCRScreenDialog.value = false
                    Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    private fun drawOverlay(
        holder: SurfaceHolder,
        heightCropPercent: Int,
        widthCropPercent: Int,
    ) {
        val canvas = holder.lockCanvas()

        val rectPaint = Paint()
        rectPaint.xfermode = PorterDuffXfermode(PorterDuff.Mode.CLEAR)
        rectPaint.style = Paint.Style.FILL
        rectPaint.color = Color.WHITE
        val outlinePaint = Paint()
        outlinePaint.style = Paint.Style.STROKE
        outlinePaint.color = Color.WHITE
        outlinePaint.strokeWidth = 4f
        val surfaceWidth = holder.surfaceFrame.width()
        val surfaceHeight = holder.surfaceFrame.height()

        val cornerRadius = 25f
        // Set rect centered in frame
        val rectTop = surfaceHeight * heightCropPercent / 2 / 100f
        val rectLeft = surfaceWidth * widthCropPercent / 2 / 100f
        val rectRight = surfaceWidth * (1 - widthCropPercent / 2 / 100f)
        val rectBottom = surfaceHeight * (1 - heightCropPercent / 2 / 100f)
        val rect = RectF(rectLeft, rectTop, rectRight, rectBottom)

        binding.viewFinderBackground.setViewFinderRect(rect)

        canvas.drawRoundRect(
            rect, cornerRadius, cornerRadius, rectPaint
        )
        canvas.drawRoundRect(
            rect, cornerRadius, cornerRadius, outlinePaint
        )

        holder.unlockCanvasAndPost(canvas)
    }

    private fun aspectRatio(width: Int, height: Int): Int {
        val previewRatio = ln(max(width, height).toDouble() / min(width, height))
        if (abs(previewRatio - ln(RATIO_4_3_VALUE))
            <= abs(previewRatio - ln(RATIO_16_9_VALUE))
        ) {
            return AspectRatio.RATIO_4_3
        }
        return AspectRatio.RATIO_16_9
    }

    override fun onRequestPermissionsResult(
        requestCode: Int,
        permissions: Array<out String>,
        grantResults: IntArray
    ) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults)
        if (requestCode == REQUEST_CODE_PERMISSIONS) {
            if (!allPermissionGranted()) {
                Toast.makeText(
                    this,
                    "Tidak mendapatkan permission.",
                    Toast.LENGTH_SHORT
                ).show()
                finish()
            } else {
                setUpCameraAndCameraUtilityWhenCameraPermissionGranted()
            }
        }
    }

    private fun allPermissionGranted() = REQUIRED_PERMISSIONS.all {
        ContextCompat.checkSelfPermission(baseContext, it) == PackageManager.PERMISSION_GRANTED
    }

}