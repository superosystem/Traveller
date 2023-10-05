package org.superosystem.traveller.data.model.flight

import com.google.gson.annotations.SerializedName

data class DetailHistoryResponse(
    @field:SerializedName("status")
    val status: String? = null,

    @field:SerializedName("message")
    val message: String? = null,

    @field:SerializedName("data")
    val data: BookingDetail? = null
)

data class BookingDetail(
    @field:SerializedName("id")
    val id: String? = null,

    @field:SerializedName("departure")
    val departure: String? = null,

    @field:SerializedName("destination")
    val destination: String? = null,

    @field:SerializedName("status")
    val status: String? = null,

    @field:SerializedName("price")
    val price: Int? = null,

    @field:SerializedName("booking_code")
    val booking_code: Int? = null,

    @field:SerializedName("passenger_name")
    val passenger_name: String? = null,

    @field:SerializedName("passenger_title")
    val passenger_title: String? = null,

    @field:SerializedName("depart_time")
    val depart_time: String? = null,

    @field:SerializedName("arrival_time")
    val arrival_time: String? = null,

    @field:SerializedName("airline")
    val airline: String? = null,

    @field:SerializedName("icon")
    val icon: String? = null
)
