package org.superosystem.traveller.data.model.flight

import com.google.gson.annotations.SerializedName

data class HistoryResponse(
    @field:SerializedName("status")
    val status: String? = null,

    @field:SerializedName("message")
    val message: String? = null,

    @field:SerializedName("data")
    val data: UserDataHistory? = null
)

data class UserDataHistory(
    @field:SerializedName("bookings")
    val bookings: List<Bookings>? = null
)

data class Bookings(
    @field:SerializedName("id")
    val id: String? = null,

    @field:SerializedName("departure")
    val departure: String? = null,

    @field:SerializedName("destination")
    val destination: String? = null,

    @field:SerializedName("booking_code")
    val booking_code: Int? = null,

    @field:SerializedName("price")
    val price: Int? = null,

    @field:SerializedName("status")
    val status: String? = null
)