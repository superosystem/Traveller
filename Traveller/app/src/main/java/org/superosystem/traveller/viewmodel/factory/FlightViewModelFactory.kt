package org.superosystem.traveller.viewmodel.factory

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import org.superosystem.traveller.data.repository.FlightRepository
import org.superosystem.traveller.viewmodel.FlightViewModel

@Suppress("UNCHECKED_CAST")
class FlightViewModelFactory(
    private val repo: FlightRepository
) : ViewModelProvider.Factory {
    override fun <T : ViewModel> create(modelClass: Class<T>): T {
        return FlightViewModel(repo) as T
    }
}