package org.superosystem.traveller.ui.main.fragment

import android.app.DatePickerDialog
import android.content.Context
import android.content.Intent
import android.os.Bundle
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.AdapterView
import android.widget.ArrayAdapter
import android.widget.Toast
import androidx.fragment.app.Fragment
import androidx.lifecycle.ViewModelProvider
import com.google.android.material.textfield.TextInputEditText
import org.superosystem.traveller.R
import org.superosystem.traveller.data.preference.SavedPreference
import org.superosystem.traveller.data.repository.AccessProfileRepository
import org.superosystem.traveller.data.repository.AuthRepository
import org.superosystem.traveller.databinding.FragmentFlightBinding
import org.superosystem.traveller.ui.flightsearchresult.FlightSearchResultActivity
import org.superosystem.traveller.utils.Constants
import org.superosystem.traveller.utils.Resources
import org.superosystem.traveller.viewmodel.AccessProfileViewModel
import org.superosystem.traveller.viewmodel.AuthViewModel
import org.superosystem.traveller.viewmodel.factory.AccessProfileFactory
import org.superosystem.traveller.viewmodel.factory.AuthViewModelFactory
import java.text.SimpleDateFormat
import java.util.Calendar
import java.util.Date
import java.util.Locale


class FlightFragment : Fragment(), View.OnClickListener {

    private var _binding: FragmentFlightBinding? = null
    private val binding get() = _binding!!

    private lateinit var savedPref: SavedPreference

    private lateinit var viewModel: AccessProfileViewModel
    private lateinit var authViewModel: AuthViewModel

    private var checkTo = false
    private var checkFrom = false
    private var checkDate = false

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        // Inflate the layout for this fragment
        _binding = FragmentFlightBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        val factory = AccessProfileFactory(AccessProfileRepository())
        viewModel = ViewModelProvider(this, factory)[AccessProfileViewModel::class.java]

        val authFactory = AuthViewModelFactory(AuthRepository())
        authViewModel = ViewModelProvider(this, authFactory)[AuthViewModel::class.java]

        //SETUP
        savedPref = SavedPreference(requireContext())
        setupDateEditText()
        itemOnClickListener()

        //GET DATA USER
        val token = savedPref.getData(Constants.ACCESS_TOKEN)
        val accessToken = "Bearer $token"

        // get user data if savedPref does not contain user data
        if (!savedPref.checkIfKeyExist(Constants.USERNAME)) {
            getDataUser(accessToken)
        }
    }

    private fun getDataUser(dataUser: String) {
        viewModel.profileUser(dataUser).observe(viewLifecycleOwner) { response ->
            if (response is Resources.Loading) {
                progressBar(true)
            } else if (response is Resources.Error) {
                progressBar(false)
                Toast.makeText(requireContext(), response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                progressBar(false)
                val result = response.data
                if (result != null) {
                    if (result.status.equals("success")) {

                        val username = result.data?.user?.name.toString()
                        val email = result.data?.user?.email.toString()
                        val photoProfile = result.data?.user?.foto_profil

                        savedPref.putData(Constants.USERNAME, username)
                        savedPref.putData(Constants.EMAIL, email)
                        savedPref.putData(Constants.PROFILE_PICTURE, photoProfile)

                    } else {
                        Log.d("PROFILE", result.status.toString())

                        Log.d("REGIS", result.status.toString())
                        val dataToken = hashMapOf(
                            "refreshToken" to savedPref.getData(Constants.REFRESH_TOKEN)
                        )

                        Log.d("REFRESH TOKEN", "observerFlightSearch: $dataToken")
                        observeUpdateToken(dataToken)
                    }
                } else {
                    Toast.makeText(requireContext(), getString(R.string.error), Toast.LENGTH_SHORT)
                        .show()
                }
            }
        }
    }

    private fun observeUpdateToken(dataToken: HashMap<String, String?>) {
        authViewModel.updateToken(dataToken).observe(viewLifecycleOwner) { response ->
            if (response is Resources.Loading) {
                progressBar(true)
            } else if (response is Resources.Error) {
                progressBar(false)
                Toast.makeText(requireContext(), response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                progressBar(false)
                val result = response.data
                if (result != null) {
                    if (result.status.equals("success")) {
                        val newAccessToken = result.data?.accessToken.toString()
                        //save new token
                        savedPref.putData(Constants.ACCESS_TOKEN, newAccessToken)

                        //get new token
                        val tokenFromAPI = (savedPref.getData(Constants.ACCESS_TOKEN))
                        val accessToken = "Bearer $tokenFromAPI"

                        Log.d("NEW ACCESS TOKEN", "observeUpdateToken: $accessToken")

                        getDataUser(accessToken)
                    } else {
                        Log.d("REGIS", result.status.toString())
                    }
                } else {
                    Toast.makeText(requireContext(), getString(R.string.error), Toast.LENGTH_SHORT)
                        .show()
                }
            }
        }
    }

    override fun onClick(p0: View?) {
        when (p0?.id) {
            R.id.searchBtn -> {
                //logout, go to flight search result
                if (checkFrom && checkTo && checkDate) {
                    startActivity(Intent(requireActivity(), FlightSearchResultActivity::class.java))
                } else {
                    Toast.makeText(requireContext(), "The form is still empty", Toast.LENGTH_SHORT)
                        .show()
                }
            }
        }
    }

    override fun onResume() {
        super.onResume()
        setupAutoTextView()
    }

    private fun setupAutoTextView() {
        val city = resources.getStringArray(R.array.data_dummy_city)

        val adapterCity: ArrayAdapter<String> =
            ArrayAdapter<String>(requireContext(), android.R.layout.select_dialog_item, city)

        //FROM
        binding.fromEditText.threshold = 1 //will start working from first character
        binding.fromEditText.setAdapter(adapterCity) //setting the adapter data into the AutoCompleteTextView
        //when clicked
        binding.fromEditText.onItemClickListener =
            AdapterView.OnItemClickListener { parent, _, position, _ ->
                val getFromMix = parent.getItemAtPosition(position).toString()
                val getCity: List<String> = getFromMix.split("-")

                println("getFromMix: $getFromMix")

                val onlyFromCity = getCity[0]
                println("onlyFromCity: $onlyFromCity")

                val onlyFromCode = getCity[1]
                println("onlyFromCode: $onlyFromCode")

                saveDataFromCity(onlyFromCity, onlyFromCode)
            }

        //TO
        binding.toEditText.threshold = 1 //will start working from first character
        binding.toEditText.setAdapter(adapterCity) //setting the adapter data into the AutoCompleteTextView
        //when clicked
        binding.toEditText.onItemClickListener =
            AdapterView.OnItemClickListener { parent, _, position, _ ->
                val getToMix = parent.getItemAtPosition(position).toString()
                val getCity: List<String> = getToMix.split("-")

                println("getToMix : $getToMix")

                val onlyToCity = getCity[0]
                println("onlyToCity : $onlyToCity")

                val onlyToCode = getCity[1]
                println("onlyToCode : $onlyToCode")

                saveDataToCity(onlyToCity, onlyToCode)
            }

        //PASSENGERS
        val dataPassengers = resources.getStringArray(R.array.data_dummy_passengers)
        val dataNumber = resources.getStringArray(R.array.data_number_passengers)
        val adapterPassengers: ArrayAdapter<String> =
            ArrayAdapter<String>(requireContext(), R.layout.dropdown_item, dataPassengers)
        //when clicked
        binding.passengersEditText.setText(dataPassengers[0])
        binding.passengersEditText.threshold = 1 //will start working from first character
        binding.passengersEditText.setAdapter(adapterPassengers) //setting the adapter data into the AutoCompleteTextView

        binding.passengersEditText.onItemClickListener =
            AdapterView.OnItemClickListener { parent, _, position, _ ->
                val passenger = parent.getItemAtPosition(position).toString()
                println("passenger: $passenger")

                val pax = dataNumber[position]
                println("pax : $pax")

                savedPref.putData(Constants.PAX, pax)

            }

        val totalPax: String

        val pax = savedPref.getData(Constants.PAX)

        if (pax == null) {
            totalPax = dataNumber[0]
            savedPref.putData(Constants.PAX, totalPax)

        }

        //SEAT CLASS
        val dataSeatClass = resources.getStringArray(R.array.data_dummy_seatClass)
        val adapterSeatClass: ArrayAdapter<String> =
            ArrayAdapter<String>(requireContext(), R.layout.dropdown_item, dataSeatClass)
        //when clicked
        binding.seatClassEditText.setText(dataSeatClass[0])
        binding.seatClassEditText.threshold = 1 //will start working from first character
        binding.seatClassEditText.setAdapter(adapterSeatClass) //setting the adapter data into the AutoCompleteTextView
        binding.seatClassEditText.onItemClickListener =
            AdapterView.OnItemClickListener { parent, _, position, _ ->
                val seatClass = parent.getItemAtPosition(position).toString()
                println("searClass : $seatClass")

                val seat = dataSeatClass[position]

                savedPref.putData(Constants.SEAT, seat)
            }

        val totalSeatClass = savedPref.getData(Constants.SEAT)

        if (totalSeatClass == null) {
            val totalSeat = dataSeatClass[0]
            savedPref.putData(Constants.SEAT, totalSeat)
        }
    }

    private fun saveDataToCity(onlyToCity: String?, codeTo: String?) {
        savedPref.putData(Constants.TO_ONLY_CITY, onlyToCity!!)
        savedPref.putData(Constants.TO_CODE, codeTo!!)

        if (
            savedPref.getData(Constants.TO_ONLY_CITY) != null &&
            savedPref.getData(Constants.TO_CODE) != null
        ) {
            checkTo = true
        }
    }

    private fun saveDataFromCity(fromOnlyCity: String?, codeFrom: String?) {
        savedPref.putData(Constants.FROM_ONLY_CITY, fromOnlyCity!!)
        savedPref.putData(Constants.FROM_CODE, codeFrom!!)

        if (
            savedPref.getData(Constants.FROM_ONLY_CITY) != null &&
            savedPref.getData(Constants.FROM_CODE) != null
        ) {
            checkFrom = true
        }
    }

    private fun setupDateEditText() {
        binding.dateEditText.transformIntoDatePicker(requireContext(), "EEEE, dd MMM yyyy", Date())
    }

    private fun TextInputEditText.transformIntoDatePicker(
        context: Context,
        format: String,
        minDate: Date? = null
    ) {
        isFocusableInTouchMode = false
        isClickable = true
        isFocusable = false

        val myCalendar = Calendar.getInstance()
        val datePickerOnDataSetListener =
            DatePickerDialog.OnDateSetListener { _, year, monthOfYear, dayOfMonth ->
                myCalendar.set(Calendar.YEAR, year)
                myCalendar.set(Calendar.MONTH, monthOfYear)
                myCalendar.set(Calendar.DAY_OF_MONTH, dayOfMonth)
                val sdf = SimpleDateFormat(format, Locale.US)
                setText(sdf.format(myCalendar.time))

                val date = binding.dateEditText.text.toString()
                savedPref.putData(Constants.DATE, date)

                if (savedPref.getData(Constants.DATE) != null) {
                    checkDate = true
                }
                println("date : $date")
            }

        setOnClickListener {
            DatePickerDialog(
                context, datePickerOnDataSetListener, myCalendar
                    .get(Calendar.YEAR), myCalendar.get(Calendar.MONTH),
                myCalendar.get(Calendar.DAY_OF_MONTH)
            ).run {
                minDate?.time?.also { datePicker.minDate = it }
                show()
            }
        }
    }

    private fun itemOnClickListener() {
        binding.searchBtn.setOnClickListener(this)
    }

    private fun progressBar(isLoading: Boolean) = with(binding) {
        if (isLoading) {
            this.progressBar.visibility = View.VISIBLE
        } else {
            this.progressBar.visibility = View.GONE
        }
    }
}