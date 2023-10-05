package org.superosystem.traveller.data.model.auth

import com.google.gson.annotations.SerializedName

data class RegisterResponse(
    @field:SerializedName("status")
    val status: String? = null,

    @field:SerializedName("message")
    val message: String? = null,

    @field:SerializedName("data")
    val data: UserDataRegister? = null
)

data class UserDataRegister(
    @field:SerializedName("user_id")
    val user_id: String? = null,
)