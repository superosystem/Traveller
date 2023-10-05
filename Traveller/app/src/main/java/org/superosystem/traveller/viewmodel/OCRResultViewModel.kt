package org.superosystem.traveller.viewmodel

import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import org.superosystem.traveller.data.repository.OCRRepository


class OCRResultViewModel(
    private val repo: OCRRepository
) : ViewModel() {

    val setLoadingOCRResultDialog = MutableLiveData(false)

    fun updateRetrievedDataToDatabase(
        accessToken: String,
        dataToBeSendToAPI: HashMap<String, String>
    ) = repo.updateRetrievedDataToDatabase(accessToken, dataToBeSendToAPI)

    fun updateBookingStatus(
        accessToken: String,
        dataBookingID: String,
        dataToBeSendToAPI1: HashMap<String, String>
    ) = repo.updateBookingStatus(accessToken, dataBookingID, dataToBeSendToAPI1)

}