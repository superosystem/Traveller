package org.superosystem.traveller.data.model.auth

import com.google.gson.annotations.SerializedName

data class LogoutResponse(
    @field:SerializedName("status")
    val status: String? = null,

    @field:SerializedName("message")
    val message: String? = null
)
