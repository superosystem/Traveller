package org.superosystem.traveller.viewmodel

import androidx.lifecycle.ViewModel
import okhttp3.MultipartBody
import okhttp3.RequestBody
import org.superosystem.traveller.data.repository.AccessProfileRepository

class AccessProfileViewModel(
    private val repo: AccessProfileRepository
) : ViewModel() {

    fun profileUser(data: String) = repo.profileUser(data)

    fun updateUser(
        accessToken: String,
        dataUsername: RequestBody?,
        dataEmail: RequestBody?,
        imageMultipart: MultipartBody.Part?
    ) = repo.updateUser(accessToken, dataUsername, dataEmail, imageMultipart)
}