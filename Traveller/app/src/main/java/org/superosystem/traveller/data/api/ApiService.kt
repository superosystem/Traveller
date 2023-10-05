package org.superosystem.traveller.data.api

import okhttp3.MultipartBody
import okhttp3.RequestBody
import org.superosystem.traveller.data.model.auth.*
import org.superosystem.traveller.data.model.flight.BookingResponse
import org.superosystem.traveller.data.model.flight.DeleteBookingResponse
import org.superosystem.traveller.data.model.flight.DetailHistoryResponse
import org.superosystem.traveller.data.model.flight.FlightSearchResponse
import org.superosystem.traveller.data.model.flight.HistoryResponse
import org.superosystem.traveller.data.model.ocr.KTPResultResponse
import org.superosystem.traveller.data.model.ocr.ScanIDCardResponse
import org.superosystem.traveller.data.model.ocr.UpdateBookingStatus
import org.superosystem.traveller.data.model.ocr.UpdatedKTPResponse
import org.superosystem.traveller.data.model.profile.AccessEditProfileResponse
import org.superosystem.traveller.data.model.profile.AccessProfileResponse
import retrofit2.Response
import retrofit2.http.*


/* Endpoint */
const val REGIS_ENDPOINT = "users"
const val LOGIN_ENDPOINT = "authentications"
const val GOOGLE_LOGIN_ENDPOINT = "auth/google"
const val TOKEN_HEADER = "Authorization"
const val NAME = "name"
const val EMAIL = "email"
const val FLIGHT_ENDPOINT = "flights"
const val FLIGHT_BOOKING = "flights/booking"
const val FLIGHT_BOOKING_UPDATE = "flights/booking/{id}"
const val DEPARTURE_QUERY = "departure"
const val DESTINATION_QUERY = "destination"
const val KTP_SCANIDCARD_ENDPOINT = "/ktp"
const val KTP_RESULT_ENDPOINT = "ktpresult"

interface ApiService {

    //REGISTER USER
    @POST(REGIS_ENDPOINT)
    suspend fun registerUser(
        @Body data: HashMap<String, String>
    ): Response<RegisterResponse>

    //LOGIN USER
    @POST(LOGIN_ENDPOINT)
    suspend fun loginUser(
        @Body data: HashMap<String, String>
    ): Response<LoginResponse>

    //UPDATE TOKEN
    @PUT(LOGIN_ENDPOINT)
    suspend fun updateToken(
        @Body data: HashMap<String, String?>
    ): Response<UpdateTokenResponse>

    //LOGOUT USER
    @HTTP(method = "DELETE", path = LOGIN_ENDPOINT, hasBody = true)
    suspend fun logoutUser(
        @Body data: HashMap<String, String?>
    ): Response<LogoutResponse>

    //LOGIN WITH GOOGLE
    @GET(GOOGLE_LOGIN_ENDPOINT)
    suspend fun googleLogin(): Response<LoginResponse>

    //GET PROFILE
    @GET(REGIS_ENDPOINT)
    suspend fun getProfile(
        @Header(TOKEN_HEADER) accessToken: String
    ): Response<AccessProfileResponse>

    //UPDATE PROFILE
    @Multipart
    @PUT(REGIS_ENDPOINT)
    suspend fun updateProfile(
        @Header(TOKEN_HEADER) accessToken: String,
        @Part(NAME) name: RequestBody,
        @Part(EMAIL) email: RequestBody
    ): Response<AccessEditProfileResponse>

    //UPDATE PROFILE W Image
    @Multipart
    @PUT(REGIS_ENDPOINT)
    suspend fun updateProfileWithImage(
        @Header(TOKEN_HEADER) accessToken: String,
        @Part(NAME) name: RequestBody,
        @Part(EMAIL) email: RequestBody,
        @Part profile_picture: MultipartBody.Part
    ): Response<AccessEditProfileResponse>

    //GET FLIGHT SEARCH BASED ON QUERY
    @GET(FLIGHT_ENDPOINT)
    suspend fun getFlightSearchWithQuery(
        @Header(TOKEN_HEADER) accessToken: String,
        @Query(DEPARTURE_QUERY) departure: String,
        @Query(DESTINATION_QUERY) destination: String,
    ): Response<FlightSearchResponse>

    //POST BOOKING
    @POST(FLIGHT_BOOKING)
    suspend fun flightBooking(
        @Header(TOKEN_HEADER) accessToken: String,
        @Body data: HashMap<String, Int>
    ): Response<BookingResponse>

    //DELETE ALL BOOKING
    @DELETE(FLIGHT_BOOKING)
    suspend fun deleteAllBooking(
        @Header(TOKEN_HEADER) accessToken: String
    ): Response<DeleteBookingResponse>

    //GET LIST HISTORY
    @GET(FLIGHT_BOOKING)
    suspend fun getHistory(
        @Header(TOKEN_HEADER) accessToken: String
    ): Response<HistoryResponse>

    //UPDATE BOOKING STATUS
    @PUT(FLIGHT_BOOKING_UPDATE)
    suspend fun updateBookingStatus(
        @Path("id") id: String,
        @Header(TOKEN_HEADER) accessToken: String,
        @Body data: HashMap<String, String>
    ): Response<UpdateBookingStatus>

    //GET DETAIL BOOKING HISTORY
    @GET(FLIGHT_BOOKING_UPDATE)
    suspend fun getDetailHistory(
        @Path("id") id: String,
        @Header(TOKEN_HEADER) accessToken: String
    ): Response<DetailHistoryResponse>

    //DELETE BOOKING BY ID
    @DELETE(FLIGHT_BOOKING_UPDATE)
    suspend fun deleteBookingById(
        @Path("id") id: String,
        @Header(TOKEN_HEADER) accessToken: String
    ): Response<DeleteBookingResponse>

    // POST Scan ID Card
    @Multipart
    @POST(KTP_SCANIDCARD_ENDPOINT)
    suspend fun scanIDCard(
        @Header(TOKEN_HEADER) accessToken: String,
        @Part file: MultipartBody.Part,
        @Part("data") data: RequestBody
    ): Response<ScanIDCardResponse>

    //GET OCR RESULT
    @GET(KTP_RESULT_ENDPOINT)
    suspend fun getOCRResult(
        @Header(TOKEN_HEADER) accessToken: String
    ): Response<KTPResultResponse>

    //UPDATE OCR RESULT
    @PUT(KTP_RESULT_ENDPOINT)
    suspend fun updateRetrievedDataToDatabase(
        @Header(TOKEN_HEADER) accessToken: String,
        @Body data: HashMap<String, String>
    ): Response<UpdatedKTPResponse>

}