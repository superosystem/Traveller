package org.superosystem.traveller.viewmodel.factory

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import org.superosystem.traveller.data.repository.AccessProfileRepository
import org.superosystem.traveller.viewmodel.AccessProfileViewModel

@Suppress("UNCHECKED_CAST")
class AccessProfileFactory(private val repo: AccessProfileRepository) : ViewModelProvider.Factory {
    override fun <T : ViewModel> create(modelClass: Class<T>): T {
        return AccessProfileViewModel(repo) as T
    }
}