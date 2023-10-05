package org.superosystem.traveller.data.model.ocr

import com.google.gson.annotations.SerializedName

data class UpdatedKTPResponse(

    @field:SerializedName("data")
    val data: List<DataItem>,

    @field:SerializedName("status")
    val status: String
)

data class DataItem(

    @field:SerializedName("nik")
    val nik: String,

    @field:SerializedName("nationality")
    val nationality: String,

    @field:SerializedName("sex")
    val sex: String,

    @field:SerializedName("name")
    val name: String,

    @field:SerializedName("title")
    val title: String,

    @field:SerializedName("married")
    val married: String
)
