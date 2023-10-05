package org.superosystem.traveller.data.model.profile

import com.google.gson.annotations.SerializedName

data class AccessEditProfileResponse(
    @field:SerializedName("status")
    val status: String? = null,

    @field:SerializedName("message")
    val message: String? = null,

    @field:SerializedName("data")
    val data: DataEditProfile? = null
)

data class DataEditProfile(
    @field:SerializedName("imageUri")
    val imageUri: String? = null
)