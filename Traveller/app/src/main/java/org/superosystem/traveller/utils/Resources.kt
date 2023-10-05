package org.superosystem.traveller.utils

sealed class Resources<out R> private constructor() {
    data class Success<out T>(val data: T) : Resources<T>()
    data class Error(val error: String) : Resources<Nothing>()
    object Loading : Resources<Nothing>()
}