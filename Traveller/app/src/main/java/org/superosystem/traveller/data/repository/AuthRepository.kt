package org.superosystem.traveller.data.repository

import androidx.lifecycle.LiveData
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.liveData
import com.google.gson.Gson
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import okhttp3.ResponseBody
import org.superosystem.traveller.data.api.RetrofitInstance
import org.superosystem.traveller.data.model.auth.LoginResponse
import org.superosystem.traveller.data.model.auth.LogoutResponse
import org.superosystem.traveller.data.model.auth.RegisterResponse
import org.superosystem.traveller.data.model.auth.UpdateTokenResponse
import org.superosystem.traveller.utils.Resources

class AuthRepository {
    fun registerUser(data: HashMap<String, String>)
            : LiveData<Resources<RegisterResponse?>> = liveData {
        emit(Resources.Loading)

        val returnValue = MutableLiveData<Resources<RegisterResponse?>>()
        val response = RetrofitInstance.API_OBJECT.registerUser(data)

        if (response.isSuccessful) {
            returnValue.value = Resources.Success(response.body())
            emitSource(returnValue)
        } else {
            val error = Gson().fromJson(
                response.errorBody()?.stringSuspending(),
                RegisterResponse::class.java
            )
            response.errorBody()?.close()
            returnValue.value = Resources.Success(error)
            emitSource(returnValue)
        }
    }

    fun loginUser(data: HashMap<String, String>): LiveData<Resources<LoginResponse?>> = liveData {
        emit(Resources.Loading)
        val returnValue = MutableLiveData<Resources<LoginResponse?>>()
        val response = RetrofitInstance.API_OBJECT.loginUser(data)
        if (response.isSuccessful) {
            returnValue.value = Resources.Success(response.body())
            emitSource(returnValue)
        } else {
            val error =
                Gson().fromJson(response.errorBody()?.stringSuspending(), LoginResponse::class.java)
            response.errorBody()?.close()
            returnValue.value = Resources.Success(error)
            emitSource(returnValue)
        }
    }

    fun loginWithGoogle(): LiveData<Resources<LoginResponse?>> = liveData {
        emit(Resources.Loading)
        val returnValue = MutableLiveData<Resources<LoginResponse?>>()
        val response = RetrofitInstance.API_OBJECT.googleLogin()
        if (response.isSuccessful) {
            returnValue.value = Resources.Success(response.body())
            emitSource(returnValue)
        } else {
            val error =
                Gson().fromJson(response.errorBody()?.stringSuspending(), LoginResponse::class.java)
            response.errorBody()?.close()
            returnValue.value = Resources.Success(error)
            emitSource(returnValue)
        }
    }

    fun updateToken(data: HashMap<String, String?>): LiveData<Resources<UpdateTokenResponse?>> =
        liveData {
            emit(Resources.Loading)
            val returnValue = MutableLiveData<Resources<UpdateTokenResponse?>>()
            val response = RetrofitInstance.API_OBJECT.updateToken(data)
            if (response.isSuccessful) {
                returnValue.value = Resources.Success(response.body())
                emitSource(returnValue)
            } else {
                val error = Gson().fromJson(
                    response.errorBody()?.stringSuspending(),
                    UpdateTokenResponse::class.java
                )
                response.errorBody()?.close()
                returnValue.value = Resources.Success(error)
                emitSource(returnValue)
            }
        }

    fun logoutUser(data: HashMap<String, String?>): LiveData<Resources<LogoutResponse?>> =
        liveData {
            emit(Resources.Loading)
            val returnValue = MutableLiveData<Resources<LogoutResponse?>>()
            val response = RetrofitInstance.API_OBJECT.logoutUser(data)
            if (response.isSuccessful) {
                returnValue.value = Resources.Success(response.body())
                emitSource(returnValue)
            } else {
                val error = Gson().fromJson(
                    response.errorBody()?.stringSuspending(),
                    LogoutResponse::class.java
                )
                response.errorBody()?.close()
                returnValue.value = Resources.Success(error)
                emitSource(returnValue)
            }
        }
}

@Suppress("BlockingMethodInNonBlockingContext")
suspend fun ResponseBody.stringSuspending() = withContext(Dispatchers.IO) { string() }