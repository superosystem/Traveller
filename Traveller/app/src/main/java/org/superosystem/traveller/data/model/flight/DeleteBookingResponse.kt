package org.superosystem.traveller.data.model.flight

import com.google.gson.annotations.SerializedName

data class DeleteBookingResponse(

    @field:SerializedName("status")
    val status: String? = null,

    @field:SerializedName("message")
    val message: String? = null
)