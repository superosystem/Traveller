package org.superosystem.traveller.data.model.ocr

import com.google.gson.annotations.SerializedName

data class IdentityCardImageCoordinate(

    @field:SerializedName("class")
    val jsonMemberClass: JsonMemberClass

)

data class NIKBoundingBoxCoordinate(

    @field:SerializedName("Ymax")
    val ymax: Int,

    @field:SerializedName("Xmin")
    val xmin: Int,

    @field:SerializedName("Ymin")
    val ymin: Int,

    @field:SerializedName("Xmax")
    val xmax: Int

)

data class NameBoundingBoxCoordinate(

    @field:SerializedName("Ymax")
    val ymax: Int,

    @field:SerializedName("Xmin")
    val xmin: Int,

    @field:SerializedName("Ymin")
    val ymin: Int,

    @field:SerializedName("Xmax")
    val xmax: Int

)

data class SexBoundingBoxCoordinate(

    @field:SerializedName("Ymax")
    val ymax: Int,

    @field:SerializedName("Xmin")
    val xmin: Int,

    @field:SerializedName("Ymin")
    val ymin: Int,

    @field:SerializedName("Xmax")
    val xmax: Int

)

data class MarriedBoundingBoxCoordinate(

    @field:SerializedName("Ymax")
    val ymax: Int,

    @field:SerializedName("Xmin")
    val xmin: Int,

    @field:SerializedName("Ymin")
    val ymin: Int,

    @field:SerializedName("Xmax")
    val xmax: Int

)

data class NationalityBoundingBoxCoordinate(

    @field:SerializedName("Ymax")
    val ymax: Int,

    @field:SerializedName("Xmin")
    val xmin: Int,

    @field:SerializedName("Ymin")
    val ymin: Int,

    @field:SerializedName("Xmax")
    val xmax: Int

)

data class JsonMemberClass(

    @field:SerializedName("NIK")
    val nik: NIKBoundingBoxCoordinate,

    @field:SerializedName("name")
    val name: NameBoundingBoxCoordinate,

    @field:SerializedName("sex")
    val sex: SexBoundingBoxCoordinate,

    @field:SerializedName("married")
    val married: MarriedBoundingBoxCoordinate,

    @field:SerializedName("nationality")
    val nationality: NationalityBoundingBoxCoordinate

)