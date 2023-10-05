package org.superosystem.traveller.data.model.auth

import com.google.gson.annotations.SerializedName

data class UpdateTokenResponse(
    @field:SerializedName("status")
    val status: String? = null,

    @field:SerializedName("message")
    val message: String? = null,

    @field:SerializedName("data")
    val data: DataUpdateToken? = null
)

data class DataUpdateToken(
    @field:SerializedName("accessToken")
    val accessToken: String? = null,
)