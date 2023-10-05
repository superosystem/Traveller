package org.superosystem.traveller.data.model.auth

import com.google.gson.annotations.SerializedName

data class LoginResponse(
    @field:SerializedName("status")
    val status: String? = null,

    @field:SerializedName("message")
    val message: String? = null,

    @field:SerializedName("data")
    val data: UserDataLogin? = null
)

data class UserDataLogin(
    @field:SerializedName("accessToken")
    val accessToken: String? = null,

    @field:SerializedName("refreshToken")
    val refreshToken: String? = null,

    var isLogin: Boolean = false
)