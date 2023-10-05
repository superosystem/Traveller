package org.superosystem.traveller.ui.ocr

import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.widget.ArrayAdapter
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.ViewModelProvider
import org.superosystem.traveller.R
import org.superosystem.traveller.data.model.ocr.DataKTPResult
import org.superosystem.traveller.data.preference.SavedPreference
import org.superosystem.traveller.data.repository.AuthRepository
import org.superosystem.traveller.data.repository.OCRRepository
import org.superosystem.traveller.databinding.ActivityOcrresultBinding
import org.superosystem.traveller.utils.Constants
import org.superosystem.traveller.utils.LoadingOCRResultDialog
import org.superosystem.traveller.utils.Resources
import org.superosystem.traveller.viewmodel.AuthViewModel
import org.superosystem.traveller.viewmodel.OCRResultViewModel
import org.superosystem.traveller.viewmodel.factory.AuthViewModelFactory
import org.superosystem.traveller.viewmodel.factory.OCRResultViewModelFactory

class OCRResultActivity : AppCompatActivity() {

    //BINDING
    private lateinit var binding: ActivityOcrresultBinding

    // Session
    private lateinit var savedPreference: SavedPreference
    private lateinit var accessToken: String
    private lateinit var dataKTPResult: DataKTPResult
    private lateinit var dataBookingID: String

    // dialog stuff
    private lateinit var loadingDialog: LoadingOCRResultDialog

    // ViewModel
    private lateinit var viewModel: OCRResultViewModel
    private lateinit var authViewModel: AuthViewModel

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityOcrresultBinding.inflate(layoutInflater)
        setContentView(binding.root)

        supportActionBar?.hide()

        loadingDialog = LoadingOCRResultDialog(this)

        setupKTPResultAndBookingIDAndToken()
        setupResult(
            dataKTPResult.nik,
            dataKTPResult.name,
            dataKTPResult.sex,
            dataKTPResult.married,
            dataKTPResult.nationality,
            dataKTPResult.title
        )

        // instantiate View Model

        val authFactory = AuthViewModelFactory(AuthRepository())
        authViewModel = ViewModelProvider(this, authFactory)[AuthViewModel::class.java]

        val factory = OCRResultViewModelFactory(OCRRepository())
        viewModel = ViewModelProvider(this, factory)[OCRResultViewModel::class.java]

        viewModel.setLoadingOCRResultDialog.observe(this) {
            if (it == true) {
                loadingDialog.startLoadingDialog()
            } else if (it == false) {
                loadingDialog.dismissLoadingDialog()
            }
        }

        binding.ivRetry.setOnClickListener {
            val intent = Intent(this, OCRScreenActivity::class.java)
            intent.putExtra("id", dataBookingID)
            startActivity(intent)
            finish()
        }

        binding.btnSubmit.setOnClickListener {

            setUpDataToBeSendToAPI()

        }

    }

    override fun onBackPressed() {
        val intent = Intent(this, OCRScreenActivity::class.java)
        intent.putExtra("id", dataBookingID)
        startActivity(intent)
        finish()
    }

    override fun onResume() {
        super.onResume()
        setupAutoTextView()
    }

    private fun setupKTPResultAndBookingIDAndToken() {
        savedPreference = SavedPreference(this)
        val tokenFromPreferences = savedPreference.getData(Constants.ACCESS_TOKEN)
        accessToken = "Bearer $tokenFromPreferences"

        dataKTPResult = intent.getParcelableExtra<DataKTPResult?>("IDCardResult") as DataKTPResult
        dataBookingID = intent.getStringExtra("BookingID") as String
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
                        viewModel.setLoadingOCRResultDialog.value = false
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
                        viewModel.setLoadingOCRResultDialog.value = false
                        Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                    }
                }

            }
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
                    viewModel.setLoadingOCRResultDialog.value = false
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
                    viewModel.setLoadingOCRResultDialog.value = false
                    Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                }
            }

        }
    }

    private fun setupResult(
        nikFromAPI: String?,
        nameFromAPI: String?,
        sexFromAPI: String?,
        marriedFromAPI: String?,
        nationalityFromAPI: String?,
        titleFromAPI: String?
    ) {

        //FOR NIK
        binding.edtNik.setText(nikFromAPI.toString())

        //FOR NAME
        binding.edtName.setText(nameFromAPI)

        //FOR TITLE PART
        val title = resources.getStringArray(R.array.data_title)

        when (titleFromAPI) {
            "Mr" -> binding.edtTitle.setText(title[0])
            "Ms" -> binding.edtTitle.setText(title[1])
            "Mrs" -> binding.edtTitle.setText(title[2])
        }

        val arrayAdapterTitle = ArrayAdapter(this, R.layout.dropdown_item_ocr_result, title)
        binding.edtTitle.setAdapter(arrayAdapterTitle)

        //FOR NATIONALITY PART
        val nationality = resources.getStringArray(R.array.data_nationality)

        when (nationalityFromAPI) {
            "WNI" -> binding.edtNationality.setText(nationality[0])
            "WNA" -> binding.edtNationality.setText(nationality[1])
        }

        val arrayAdapterNationality =
            ArrayAdapter(this, R.layout.dropdown_item_ocr_result, nationality)
        binding.edtNationality.setAdapter(arrayAdapterNationality)

        //FOR GENDER PART
        val gender = resources.getStringArray(R.array.data_gender)

        when (sexFromAPI) {
            "Male" -> binding.edtGender.setText(gender[0])
            "Female" -> binding.edtGender.setText(gender[1])
        }

        val arrayAdapterGender = ArrayAdapter(this, R.layout.dropdown_item_ocr_result, gender)
        binding.edtGender.setAdapter(arrayAdapterGender)

        //FOR MARITAL STATUS PART
        val maritalStatus = resources.getStringArray(R.array.data_marital_status)

        when (marriedFromAPI) {
            "Single" -> binding.edtMaritalStatus.setText(maritalStatus[0])
            "Married" -> binding.edtMaritalStatus.setText(maritalStatus[1])
        }

        val arrayAdapterMaritalStatus =
            ArrayAdapter(this, R.layout.dropdown_item_ocr_result, maritalStatus)
        binding.edtMaritalStatus.setAdapter(arrayAdapterMaritalStatus)

        binding.tvContactDesc.text = savedPreference.getData(Constants.USERNAME)
        binding.tvContactDesc2.text = savedPreference.getData(Constants.EMAIL)
    }


    private fun setupAutoTextView() {
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
    }

}