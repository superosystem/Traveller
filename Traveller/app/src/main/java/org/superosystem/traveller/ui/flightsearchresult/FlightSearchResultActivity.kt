package org.superosystem.traveller.ui.flightsearchresult

import android.content.Intent
import android.content.res.Configuration
import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.Toast
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.ViewModelProvider
import androidx.recyclerview.widget.GridLayoutManager
import androidx.recyclerview.widget.LinearLayoutManager
import org.superosystem.traveller.R
import org.superosystem.traveller.data.preference.SavedPreference
import org.superosystem.traveller.data.repository.AuthRepository
import org.superosystem.traveller.data.repository.FlightRepository
import org.superosystem.traveller.databinding.ActivityFlightSearchResultBinding
import org.superosystem.traveller.ui.adapter.SearchListAdapter
import org.superosystem.traveller.ui.ocr.ManualInputActivity
import org.superosystem.traveller.ui.ocr.OCRScreenActivity
import org.superosystem.traveller.utils.Constants
import org.superosystem.traveller.utils.Resources
import org.superosystem.traveller.viewmodel.AuthViewModel
import org.superosystem.traveller.viewmodel.FlightViewModel
import org.superosystem.traveller.viewmodel.factory.AuthViewModelFactory
import org.superosystem.traveller.viewmodel.factory.FlightViewModelFactory


class FlightSearchResultActivity : AppCompatActivity() {

    private lateinit var binding: ActivityFlightSearchResultBinding

    private lateinit var savedPref: SavedPreference
    private lateinit var accessToken: String
    private lateinit var list: SearchListAdapter

    //API
    private lateinit var viewModel: FlightViewModel
    private lateinit var authViewModel: AuthViewModel

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityFlightSearchResultBinding.inflate(layoutInflater)
        setContentView(binding.root)

        supportActionBar?.hide()

        savedPref = SavedPreference(this)
        configRecyclerView()

        val cityTo = (savedPref.getData(Constants.TO_ONLY_CITY))?.lowercase()
        val cityFrom = (savedPref.getData(Constants.FROM_ONLY_CITY))?.lowercase()

        val tokenFromAPI = (savedPref.getData(Constants.ACCESS_TOKEN))!!
        accessToken = "Bearer $tokenFromAPI"

        //CREATE API CONNECTION
        val factory = FlightViewModelFactory(FlightRepository())
        viewModel = ViewModelProvider(this, factory)[FlightViewModel::class.java]
        observerFlightSearch(accessToken, cityFrom, cityTo)

        val authFactory = AuthViewModelFactory(AuthRepository())
        authViewModel = ViewModelProvider(this, authFactory)[AuthViewModel::class.java]

        binding.fromTv.text = savedPref.getData(Constants.FROM_ONLY_CITY)
        binding.toTv.text = savedPref.getData(Constants.TO_ONLY_CITY)
        binding.dateTv.text = savedPref.getData(Constants.DATE)
        binding.paxTv.text = String.format(
            getString(R.string.FlightSearchResultActivity_paxTv),
            savedPref.getData(Constants.PAX)
        )
        binding.seatClassTv.text = savedPref.getData(Constants.SEAT)

        binding.ivBack.setOnClickListener {
            savedPref.putData(Constants.PAX, null)
            savedPref.putData(Constants.SEAT, null)
            finish()
        }
    }

    private fun observerFlightBook(accessToken: String, flightID: HashMap<String, Int>) {
        viewModel.flightBook(accessToken, flightID).observe(this) { response ->
            if (response is Resources.Loading) {
                enableProgressBar()
            } else if (response is Resources.Error) {
                disableProgressBar()
                Toast.makeText(this, response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                val result = response.data
                if (result != null) {
                    if (result.status == "success") {

                        val view = View.inflate(this, R.layout.scanning_option_dialog, null)

                        disableProgressBar()

                        AlertDialog.Builder(this, R.style.MyAlertDialogTheme)
                            .setView(view)
                            .setNegativeButton("No, let me fill in manually") { _, _ ->
                                //Intent to manual input screen
                                val intent = Intent(
                                    this@FlightSearchResultActivity,
                                    ManualInputActivity::class.java
                                )
                                intent.putExtra("id", result.data?.bookingId)
                                startActivity(intent)
                            }
                            .setPositiveButton("Continue") { _, _ ->
                                //Intent to OCR Screen
                                val intent = Intent(
                                    this@FlightSearchResultActivity,
                                    OCRScreenActivity::class.java
                                )
                                intent.putExtra("id", result.data?.bookingId)
                                startActivity(intent)
                            }
                            .show()

                    } else {
                        val dataToken = hashMapOf(
                            "refreshToken" to savedPref.getData(Constants.REFRESH_TOKEN)
                        )
                        observeUpdateTokenObserverFlightBook(dataToken, flightID)
                    }
                } else {
                    disableProgressBar()
                    Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    private fun observeUpdateTokenObserverFlightBook(
        dataToken: HashMap<String, String?>,
        flightID: HashMap<String, Int>
    ) {
        authViewModel.updateToken(dataToken).observe(this) { response ->
            if (response is Resources.Error) {
                disableProgressBar()
                Toast.makeText(this, response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                val result = response.data
                if (result != null) {
                    if (result.status.equals("success")) {

                        val newAccessToken = result.data?.accessToken.toString()
                        //save new token
                        savedPref.putData(Constants.ACCESS_TOKEN, newAccessToken)

                        //get new token
                        val tokenFromAPI = (savedPref.getData(Constants.ACCESS_TOKEN))
                        val accessToken = "Bearer $tokenFromAPI"

                        observerFlightBook(accessToken, flightID)

                    } else {
                        Log.d("REGIS", result.status.toString())
                    }
                } else {
                    disableProgressBar()
                    Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    private fun observerFlightSearch(accessToken: String, cityFrom: String?, cityTo: String?) {
        viewModel.flightSearch(accessToken, cityFrom!!, cityTo!!).observe(this) { response ->
            if (response is Resources.Loading) {
                enableProgressBar()
            } else if (response is Resources.Error) {
                disableProgressBar()
                Toast.makeText(this, response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                val result = response.data
                if (result != null) {
                    if (result.status.equals("success")) {
                        disableProgressBar()
                        list.differAsync.submitList((result.data?.flights)?.reversed())
                    } else {

                        val dataToken = hashMapOf(
                            "refreshToken" to savedPref.getData(Constants.REFRESH_TOKEN)
                        )

                        observeUpdateTokenObserverFlightSearch(dataToken)

                    }
                } else {
                    disableProgressBar()
                    Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    private fun observeUpdateTokenObserverFlightSearch(dataToken: HashMap<String, String?>) {
        authViewModel.updateToken(dataToken).observe(this) { response ->
            if (response is Resources.Error) {
                disableProgressBar()
                Toast.makeText(this, response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                val result = response.data
                if (result != null) {
                    if (result.status.equals("success")) {
                        val cityTo = (savedPref.getData(Constants.TO_ONLY_CITY))?.lowercase()
                        val cityFrom = (savedPref.getData(Constants.FROM_ONLY_CITY))?.lowercase()

                        val newAccessToken = result.data?.accessToken.toString()
                        //save new token
                        savedPref.putData(Constants.ACCESS_TOKEN, newAccessToken)

                        //get new token
                        val tokenFromAPI = (savedPref.getData(Constants.ACCESS_TOKEN))
                        val accessToken = "Bearer $tokenFromAPI"

                        observerFlightSearch(accessToken, cityFrom, cityTo)
                    } else {
                        Log.d("REGIS", result.status.toString())
                    }
                } else {
                    disableProgressBar()
                    Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    private fun configRecyclerView() {

        list = SearchListAdapter(this)
        binding.rvSearchResultTickets.apply {
            adapter = list
            layoutManager =
                if (this.resources.configuration.orientation == Configuration.ORIENTATION_LANDSCAPE) {
                    GridLayoutManager(context, 2)
                } else {
                    LinearLayoutManager(context)
                }
        }

        list.setOnItemClickCallback(object : SearchListAdapter.OnItemClickCallback {
            override fun onItemClicked(id: String?) {

                val flightID = hashMapOf(
                    "id" to id?.toInt()!!
                )

                observerFlightBook(accessToken, flightID)
            }

        })

    }

    private fun enableProgressBar() {
        binding.progressBar.visibility = View.VISIBLE
    }

    private fun disableProgressBar() {
        binding.progressBar.visibility = View.INVISIBLE
    }
}