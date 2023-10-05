package org.superosystem.traveller.viewmodel.factory

import androidx.lifecycle.ViewModel
import androidx.lifecycle.ViewModelProvider
import org.superosystem.traveller.data.repository.AuthRepository
import org.superosystem.traveller.viewmodel.AuthViewModel

class AuthViewModelFactory(private val authRepository: AuthRepository) : ViewModelProvider.Factory {
    override fun <T : ViewModel> create(modelClass: Class<T>): T {
        return AuthViewModel(authRepository) as T
    }
}