package org.superosystem.traveller.data.repository

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.liveData
import com.google.gson.Gson
import org.superosystem.traveller.data.api.RetrofitInstance
import org.superosystem.traveller.data.model.flight.BookingResponse
import org.superosystem.traveller.data.model.flight.DeleteBookingResponse
import org.superosystem.traveller.data.model.flight.DetailHistoryResponse
import org.superosystem.traveller.data.model.flight.FlightSearchResponse
import org.superosystem.traveller.data.model.flight.HistoryResponse
import org.superosystem.traveller.utils.Resources

class FlightRepository {

    fun flightSearch(
        accessToken: String,
        departure: String,
        destination: String
    ): LiveData<Resources<FlightSearchResponse?>> = liveData {
        emit(Resources.Loading)

        val returnValue = MutableLiveData<Resources<FlightSearchResponse?>>()
        val response = RetrofitInstance.API_OBJECT.getFlightSearchWithQuery(
            accessToken,
            departure,
            destination
        )

        if (response.isSuccessful) {
            returnValue.value = Resources.Success(response.body())
            emitSource(returnValue)
        } else {
            val error = Gson().fromJson(
                response.errorBody()?.stringSuspending(),
                FlightSearchResponse::class.java
            )
            response.errorBody()?.close()
            returnValue.value = Resources.Success(error)
            emitSource(returnValue)
        }
    }

    fun flightBook(
        accessToken: String,
        flightID: HashMap<String, Int>
    ): LiveData<Resources<BookingResponse?>> = liveData {
        emit(Resources.Loading)

        val returnValue = MutableLiveData<Resources<BookingResponse?>>()
        val response = RetrofitInstance.API_OBJECT.flightBooking(accessToken, flightID)

        if (response.isSuccessful) {
            returnValue.value = Resources.Success(response.body())
            emitSource(returnValue)
        } else {
            val error = Gson().fromJson(
                response.errorBody()?.stringSuspending(),
                BookingResponse::class.java
            )
            response.errorBody()?.close()
            returnValue.value = Resources.Success(error)
            emitSource(returnValue)
        }
    }

    fun history(
        accessToken: String
    ): LiveData<Resources<HistoryResponse?>> = liveData {
        emit(Resources.Loading)

        val returnValue = MutableLiveData<Resources<HistoryResponse?>>()
        val response = RetrofitInstance.API_OBJECT.getHistory(accessToken)

        if (response.isSuccessful) {
            returnValue.value = Resources.Success(response.body())
            emitSource(returnValue)
        } else {
            val error = Gson().fromJson(
                response.errorBody()?.stringSuspending(),
                HistoryResponse::class.java
            )
            response.errorBody()?.close()
            returnValue.value = Resources.Success(error)
            emitSource(returnValue)
        }
    }

    fun detailHistory(
        dataBookingID: String,
        accessToken: String
    ): LiveData<Resources<DetailHistoryResponse?>> = liveData {
        emit(Resources.Loading)

        val returnValue = MutableLiveData<Resources<DetailHistoryResponse?>>()
        val response = RetrofitInstance.API_OBJECT.getDetailHistory(dataBookingID, accessToken)

        if (response.isSuccessful) {
            returnValue.value = Resources.Success(response.body())
            emitSource(returnValue)
        } else {
            val error = Gson().fromJson(
                response.errorBody()?.stringSuspending(),
                DetailHistoryResponse::class.java
            )
            response.errorBody()?.close()
            returnValue.value = Resources.Success(error)
            emitSource(returnValue)
        }
    }

    fun deleteBookingById(
        dataBookingID: String,
        accessToken: String
    ): LiveData<Resources<DeleteBookingResponse?>> = liveData {
        emit(Resources.Loading)

        val returnValue = MutableLiveData<Resources<DeleteBookingResponse?>>()
        val response = RetrofitInstance.API_OBJECT.deleteBookingById(dataBookingID, accessToken)

        if (response.isSuccessful) {
            returnValue.value = Resources.Success(response.body())
            emitSource(returnValue)
        } else {
            val error = Gson().fromJson(
                response.errorBody()?.stringSuspending(),
                DeleteBookingResponse::class.java
            )
            response.errorBody()?.close()
            returnValue.value = Resources.Success(error)
            emitSource(returnValue)
        }
    }

    fun deleteAllBooking(
        accessToken: String
    ): LiveData<Resources<DeleteBookingResponse?>> = liveData {
        emit(Resources.Loading)

        val returnValue = MutableLiveData<Resources<DeleteBookingResponse?>>()
        val response = RetrofitInstance.API_OBJECT.deleteAllBooking(accessToken)

        if (response.isSuccessful) {
            returnValue.value = Resources.Success(response.body())
            emitSource(returnValue)
        } else {
            val error = Gson().fromJson(
                response.errorBody()?.stringSuspending(),
                DeleteBookingResponse::class.java
            )
            response.errorBody()?.close()
            returnValue.value = Resources.Success(error)
            emitSource(returnValue)
        }
    }
}