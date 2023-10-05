package org.superosystem.traveller.data.model.ocr

import com.google.gson.annotations.SerializedName

data class ScanIDCardResponse(

    @field:SerializedName("data")
    val data: Data,

    @field:SerializedName("message")
    val message: String,

    @field:SerializedName("status")
    val status: String

)

data class Data(

    @field:SerializedName("imageId")
    val imageId: String

)