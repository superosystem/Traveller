package org.superosystem.traveller.viewmodel

import androidx.lifecycle.ViewModel
import org.superosystem.traveller.data.repository.AuthRepository

class AuthViewModel(
    private val repo: AuthRepository,
) : ViewModel() {
    fun registerUser(data: HashMap<String, String>) = repo.registerUser(data)
    fun loginUser(data: HashMap<String, String>) = repo.loginUser(data)
    fun loginWithGoogle() = repo.loginWithGoogle()
    fun updateToken(data: HashMap<String, String?>) = repo.updateToken(data)
    fun logoutUser(data: HashMap<String, String?>) = repo.logoutUser(data)
}