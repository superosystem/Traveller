package org.superosystem.traveller.data.repository

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.liveData
import com.google.gson.Gson
import okhttp3.MultipartBody
import okhttp3.RequestBody
import org.superosystem.traveller.data.api.RetrofitInstance
import org.superosystem.traveller.data.model.profile.AccessEditProfileResponse
import org.superosystem.traveller.data.model.profile.AccessProfileResponse
import org.superosystem.traveller.utils.Resources

class AccessProfileRepository {
    fun profileUser(
        data: String
    ): LiveData<Resources<AccessProfileResponse?>> = liveData {
        emit(Resources.Loading)

        val returnValue = MutableLiveData<Resources<AccessProfileResponse?>>()
        val response = RetrofitInstance.API_OBJECT.getProfile(data)

        if (response.isSuccessful) {
            returnValue.value = Resources.Success(response.body())
            emitSource(returnValue)
        } else {
            val error = Gson().fromJson(
                response.errorBody()?.stringSuspending(),
                AccessProfileResponse::class.java
            )
            response.errorBody()?.close()
            returnValue.value = Resources.Success(error)
            emitSource(returnValue)
        }
    }

    fun updateUser(
        accessToken: String,
        dataUsername: RequestBody?,
        dataEmail: RequestBody?,
        imageMultipart: MultipartBody.Part?
    ): LiveData<Resources<AccessEditProfileResponse?>> = liveData {
        emit(Resources.Loading)

        val returnValue = MutableLiveData<Resources<AccessEditProfileResponse?>>()
        val response = if (imageMultipart != null) {
            RetrofitInstance.API_OBJECT.updateProfileWithImage(
                accessToken,
                dataUsername!!,
                dataEmail!!,
                imageMultipart
            )
        } else {
            RetrofitInstance.API_OBJECT.updateProfile(accessToken, dataUsername!!, dataEmail!!)
        }

        if (response.isSuccessful) {
            returnValue.value = Resources.Success(response.body())
            emitSource(returnValue)
        } else {
            val error = Gson().fromJson(
                response.errorBody()?.stringSuspending(),
                AccessEditProfileResponse::class.java
            )
            response.errorBody()?.close()
            returnValue.value = Resources.Success(error)
            emitSource(returnValue)
        }
    }
}