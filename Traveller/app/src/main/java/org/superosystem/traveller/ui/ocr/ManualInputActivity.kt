package org.superosystem.traveller.ui.ocr

import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.widget.ArrayAdapter
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.ViewModelProvider
import org.superosystem.traveller.R
import org.superosystem.traveller.data.preference.SavedPreference
import org.superosystem.traveller.data.repository.OCRRepository
import org.superosystem.traveller.databinding.ActivityManualInputBinding
import org.superosystem.traveller.utils.Constants
import org.superosystem.traveller.utils.Resources
import org.superosystem.traveller.viewmodel.AuthViewModel
import org.superosystem.traveller.viewmodel.OCRResultViewModel
import org.superosystem.traveller.viewmodel.factory.OCRResultViewModelFactory

class ManualInputActivity : AppCompatActivity() {
    //BINDING
    private lateinit var binding: ActivityManualInputBinding

    // Session
    private lateinit var savedPreference: SavedPreference
    private lateinit var accessToken: String
    private lateinit var dataBookingID: String

    // ViewModel
    private lateinit var viewModel: OCRResultViewModel
    private lateinit var authViewModel: AuthViewModel

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityManualInputBinding.inflate(layoutInflater)
        setContentView(binding.root)

        supportActionBar?.hide()

        savedPreference = SavedPreference(this)
        val tokenFromPreferences = savedPreference.getData(Constants.ACCESS_TOKEN)
        accessToken = "Bearer $tokenFromPreferences"

        dataBookingID = (intent.getStringExtra("id") as String).toString()

        val factory = OCRResultViewModelFactory(OCRRepository())
        viewModel = ViewModelProvider(this, factory)[OCRResultViewModel::class.java]

        setUpForm()

        binding.btnSubmit.setOnClickListener {
            setUpDataToBeSendToAPI()
        }
    }

    private fun setUpDataToBeSendToAPI() {

        val title = binding.edtTitle.text.toString().trim()
        val name = binding.edtName.text.toString().trim()
        val nationality = binding.edtNationality.text.toString().trim()
        val nik = binding.edtNik.text.toString().trim()
        val gender = binding.edtGender.text.toString().trim()
        val maritalStatus = binding.edtMaritalStatus.text.toString().trim()

        val checkIfAnyEmptyFields = title.isEmpty() ||
                name.isEmpty() ||
                nationality.isEmpty() ||
                nik.isEmpty() ||
                gender.isEmpty() ||
                maritalStatus.isEmpty()

        if (checkIfAnyEmptyFields) {
            Toast.makeText(
                this,
                "Please fill in all the details before submitting",
                Toast.LENGTH_SHORT
            ).show()
        } else {

            viewModel.setLoadingOCRResultDialog.value = true

            val dataToBeSendToAPI = hashMapOf(
                "name" to name,
                "nik" to nik,
                "nationality" to nationality,
                "sex" to gender,
                "married" to maritalStatus,
                "title" to title
            )
            observerUpdateRetrievedDataToDatabase(accessToken, dataToBeSendToAPI, dataBookingID)
        }

    }

    private fun observerUpdateRetrievedDataToDatabase(
        accessToken: String,
        dataToBeSendToAPI: HashMap<String, String>,
        dataBookingID: String
    ) {

        viewModel.updateRetrievedDataToDatabase(accessToken, dataToBeSendToAPI)
            .observe(this) { response ->

                if (response is Resources.Error) {
                    viewModel.setLoadingOCRResultDialog.value = false
                    Toast.makeText(this, response.error, Toast.LENGTH_SHORT).show()
                } else if (response is Resources.Success) {
                    val result = response.data
                    if (result != null) {
                        if (result.status == "success") {

                            val dataToBeSendToAPI1 = hashMapOf(
                                "title" to dataToBeSendToAPI["title"]!!,
                                "name" to dataToBeSendToAPI["name"]!!
                            )

                            observerUpdateBookingStatus(
                                accessToken,
                                dataBookingID,
                                dataToBeSendToAPI1
                            )

                        } else {

                            val dataToken = hashMapOf(
                                "refreshToken" to savedPreference.getData(Constants.REFRESH_TOKEN)
                            )

                            observeUpdateTokenForObserverUpdateRetrievedDataToDatabase(
                                dataToken,
                                dataToBeSendToAPI,
                                dataBookingID
                            )

                        }
                    } else {
                        Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                    }
                }

            }

    }

    private fun observerUpdateBookingStatus(
        accessToken: String,
        dataBookingID: String,
        dataToBeSendToAPI1: HashMap<String, String>
    ) {
        viewModel.updateBookingStatus(accessToken, dataBookingID, dataToBeSendToAPI1)
            .observe(this) { response ->

                if (response is Resources.Error) {
                    viewModel.setLoadingOCRResultDialog.value = false
                    Toast.makeText(this, response.error, Toast.LENGTH_SHORT).show()
                } else if (response is Resources.Success) {
                    val result = response.data
                    if (result != null) {
                        if (result.status == "success") {

                            val intent = Intent(this, SuccessPageActivity::class.java)

                            viewModel.setLoadingOCRResultDialog.value = false

                            startActivity(intent)
                            finish()
                        } else {
                            val dataToken = hashMapOf(
                                "refreshToken" to savedPreference.getData(Constants.REFRESH_TOKEN)
                            )
                            observeUpdateTokenForObserverUpdateBookingStatus(
                                dataToken,
                                dataBookingID,
                                dataToBeSendToAPI1
                            )
                        }
                    } else {
                        Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                    }
                }

            }
    }

    private fun observeUpdateTokenForObserverUpdateRetrievedDataToDatabase(
        dataToken: HashMap<String, String?>,
        dataToBeSendToAPI: HashMap<String, String>,
        dataBookingID: String
    ) {
        authViewModel.updateToken(dataToken).observe(this) { response ->

            if (response is Resources.Error) {
                viewModel.setLoadingOCRResultDialog.value = false
                Toast.makeText(this, response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                val result = response.data
                if (result != null) {
                    if (result.status.equals("success")) {

                        val newAccessToken = result.data?.accessToken.toString()
                        //save new token
                        savedPreference.putData(Constants.ACCESS_TOKEN, newAccessToken)

                        //get new token
                        val tokenFromAPI = (savedPreference.getData(Constants.ACCESS_TOKEN))
                        val accessToken = "Bearer $tokenFromAPI"

                        observerUpdateRetrievedDataToDatabase(
                            accessToken,
                            dataToBeSendToAPI,
                            dataBookingID
                        )
                    } else {
                        Log.d("REGIS", result.status.toString())
                    }
                } else {
                    Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                }
            }

        }
    }

    private fun setUpForm() {
        //FOR NIK
        binding.edtNik.setText("")

        //FOR NAME
        binding.edtName.setText("")

        //FOR TITLE PART
        val title = resources.getStringArray(R.array.data_title)
        val arrayAdapterTitle = ArrayAdapter(this, R.layout.dropdown_item_ocr_result, title)
        binding.edtTitle.setAdapter(arrayAdapterTitle)

        //FOR NATIONALITY PART
        val nationality = resources.getStringArray(R.array.data_nationality)
        val arrayAdapterNationality =
            ArrayAdapter(this, R.layout.dropdown_item_ocr_result, nationality)
        binding.edtNationality.setAdapter(arrayAdapterNationality)

        //FOR GENDER PART
        val gender = resources.getStringArray(R.array.data_gender)
        val arrayAdapterGender = ArrayAdapter(this, R.layout.dropdown_item_ocr_result, gender)
        binding.edtGender.setAdapter(arrayAdapterGender)

        //FOR MARITAL STATUS PART
        val maritalStatus = resources.getStringArray(R.array.data_marital_status)
        val arrayAdapterMaritalStatus =
            ArrayAdapter(this, R.layout.dropdown_item_ocr_result, maritalStatus)
        binding.edtMaritalStatus.setAdapter(arrayAdapterMaritalStatus)

        binding.tvContactDesc.text = savedPreference.getData(Constants.USERNAME)
        binding.tvContactDesc2.text = savedPreference.getData(Constants.EMAIL)
    }

    private fun observeUpdateTokenForObserverUpdateBookingStatus(
        dataToken: HashMap<String, String?>,
        dataBookingID: String,
        dataToBeSendToAPI1: HashMap<String, String>
    ) {
        authViewModel.updateToken(dataToken).observe(this) { response ->

            if (response is Resources.Error) {
                viewModel.setLoadingOCRResultDialog.value = false
                Toast.makeText(this, response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                val result = response.data
                if (result != null) {
                    if (result.status.equals("success")) {

                        val newAccessToken = result.data?.accessToken.toString()
                        //save new token
                        savedPreference.putData(Constants.ACCESS_TOKEN, newAccessToken)

                        //get new token
                        val tokenFromAPI = (savedPreference.getData(Constants.ACCESS_TOKEN))
                        val accessToken = "Bearer $tokenFromAPI"


                        observerUpdateBookingStatus(accessToken, dataBookingID, dataToBeSendToAPI1)
                    } else {
                        Log.d("REGIS", result.status.toString())
                    }
                } else {
                    Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                }
            }

        }
    }
}