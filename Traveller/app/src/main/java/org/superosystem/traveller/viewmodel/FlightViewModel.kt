package org.superosystem.traveller.viewmodel

import androidx.lifecycle.ViewModel
import org.superosystem.traveller.data.repository.FlightRepository

class FlightViewModel(
    private val repo: FlightRepository
) : ViewModel() {

    fun flightSearch(
        accessToken: String,
        departure: String,
        destination: String
    ) = repo.flightSearch(accessToken, departure, destination)

    fun flightBook(
        accessToken: String,
        flightID: HashMap<String, Int>
    ) = repo.flightBook(accessToken, flightID)

    fun history(accessToken: String) = repo.history(accessToken)

    fun detailHistory(
        dataBookingID: String,
        accessToken: String
    ) = repo.detailHistory(dataBookingID, accessToken)

    fun deleteBookingById(
        dataBookingID: String,
        accessToken: String
    ) = repo.deleteBookingById(dataBookingID, accessToken)

    fun deleteAllBooking(accessToken: String) = repo.deleteAllBooking(accessToken)
}