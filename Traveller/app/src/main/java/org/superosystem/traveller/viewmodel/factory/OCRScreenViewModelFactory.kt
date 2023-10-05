package org.superosystem.traveller.viewmodel.factory

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import org.superosystem.traveller.data.repository.OCRRepository
import org.superosystem.traveller.viewmodel.OCRScreenViewModel

@Suppress("UNCHECKED_CAST")
class OCRScreenViewModelFactory(
    private val repo: OCRRepository
) : ViewModelProvider.Factory {
    override fun <T : ViewModel> create(modelClass: Class<T>): T {
        return OCRScreenViewModel(repo) as T
    }
}