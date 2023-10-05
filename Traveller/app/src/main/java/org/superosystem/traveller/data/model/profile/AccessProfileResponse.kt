package org.superosystem.traveller.data.model.profile

import com.google.gson.annotations.SerializedName

data class AccessProfileResponse(
    @field:SerializedName("status")
    val status: String? = null,

    @field:SerializedName("data")
    val data: AccessProfileData? = null
)

data class AccessProfileData(
    @field:SerializedName("user")
    val user: User? = null
)

data class User(
    @field:SerializedName("id")
    val id: String? = null,

    @field:SerializedName("name")
    val name: String? = null,

    @field:SerializedName("email")
    val email: String? = null,

    @field:SerializedName("password")
    val password: String? = null,

    @field:SerializedName("profile_picture")
    val foto_profil: String? = null
)