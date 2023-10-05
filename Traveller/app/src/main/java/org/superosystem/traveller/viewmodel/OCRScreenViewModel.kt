package org.superosystem.traveller.viewmodel

import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import okhttp3.MultipartBody
import okhttp3.RequestBody
import org.superosystem.traveller.data.repository.OCRRepository
import org.superosystem.traveller.ui.ocr.OCRScreenActivity.Companion.DESIRED_HEIGHT_CROP_PERCENT
import org.superosystem.traveller.ui.ocr.OCRScreenActivity.Companion.DESIRED_WIDTH_CROP_PERCENT

class OCRScreenViewModel(
    private val repo: OCRRepository
) : ViewModel() {

    // We set desired crop percentages to avoid having to analyze the whole image from the live
    // camera feed. However, we are not guaranteed what aspect ratio we will get from the camera, so
    // we use the first frame we get back from the camera to update these crop percentages based on
    // the actual aspect ratio of images.
    val imageCropPercentages = MutableLiveData<Pair<Int, Int>>()
        .apply { value = Pair(DESIRED_HEIGHT_CROP_PERCENT, DESIRED_WIDTH_CROP_PERCENT) }

    val setLoadingOCRScreenDialog = MutableLiveData(false)

    // Post Scan ID Card
    fun scanIDCard(accessToken: String, file: MultipartBody.Part, data: RequestBody) =
        repo.scanIDCard(accessToken, file, data)

    // Retrieve ID Card Result
    fun retrieveIDCardResult(accessToken: String) =
        repo.retrieveIDCardResult(accessToken)

}