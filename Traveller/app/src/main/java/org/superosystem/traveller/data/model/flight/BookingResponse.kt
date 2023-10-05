package org.superosystem.traveller.data.model.flight

import com.google.gson.annotations.SerializedName

data class BookingResponse(
    @field:SerializedName("status")
    val status: String,

    @field:SerializedName("message")
    val message: String? = null,

    @field:SerializedName("data")
    val data: DataBooking? = null
)

data class DataBooking(
    @field:SerializedName("bookingId")
    val bookingId: String? = null
)
