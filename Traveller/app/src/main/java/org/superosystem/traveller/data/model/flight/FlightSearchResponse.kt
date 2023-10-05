package org.superosystem.traveller.data.model.flight

import com.google.gson.annotations.SerializedName

data class FlightSearchResponse(
    @field:SerializedName("status")
    val status: String? = null,

    @field:SerializedName("data")
    val data: DataFlight? = null
)

data class DataFlight(
    @field:SerializedName("flights")
    val flights: List<Flights>? = null
)

data class Flights(
    @field:SerializedName("id")
    val id: String? = null,

    @field:SerializedName("airline")
    val airline: String? = null,

    @field:SerializedName("icon")
    val icon: String? = null,

    @field:SerializedName("depart_time")
    val depart_time: String? = null,

    @field:SerializedName("arrival_time")
    val arrival_time: String? = null,

    @field:SerializedName("departure")
    val departure: String? = null,

    @field:SerializedName("destination")
    val destination: String? = null,

    @field:SerializedName("price")
    val price: Int? = null
)