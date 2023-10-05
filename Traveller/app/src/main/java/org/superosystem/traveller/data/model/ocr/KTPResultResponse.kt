package org.superosystem.traveller.data.model.ocr

import android.os.Parcelable
import com.google.gson.annotations.SerializedName
import kotlinx.parcelize.Parcelize

data class KTPResultResponse(
    @field:SerializedName("status")
    val status: String? = null,

    @field:SerializedName("data")
    val data: DataKTPResult? = null
)

@Parcelize
data class DataKTPResult(
    @field:SerializedName("nik")
    val nik: String? = null,

    @field:SerializedName("name")
    val name: String? = null,

    @field:SerializedName("sex")
    val sex: String? = null,

    @field:SerializedName("married")
    val married: String? = null,

    @field:SerializedName("nationality")
    val nationality: String? = null,

    @field:SerializedName("title")
    val title: String? = null
) : Parcelable