package org.superosystem.traveller.ui.main.historydetail

import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.Toast
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import androidx.core.content.ContextCompat
import androidx.lifecycle.ViewModelProvider
import androidx.navigation.navArgs
import com.bumptech.glide.Glide
import com.bumptech.glide.request.RequestOptions
import org.superosystem.traveller.R
import org.superosystem.traveller.data.preference.SavedPreference
import org.superosystem.traveller.data.repository.AuthRepository
import org.superosystem.traveller.data.repository.FlightRepository
import org.superosystem.traveller.databinding.ActivityHistoryDetailBinding
import org.superosystem.traveller.utils.Constants
import org.superosystem.traveller.utils.Resources
import org.superosystem.traveller.viewmodel.AuthViewModel
import org.superosystem.traveller.viewmodel.FlightViewModel
import org.superosystem.traveller.viewmodel.factory.AuthViewModelFactory
import org.superosystem.traveller.viewmodel.factory.FlightViewModelFactory
import java.text.NumberFormat
import java.util.*

class HistoryDetailActivity : AppCompatActivity() {

    private lateinit var binding: ActivityHistoryDetailBinding

    private lateinit var savedPref: SavedPreference
    private lateinit var accessToken: String

    //API
    private lateinit var viewModel: FlightViewModel
    private lateinit var authViewModel: AuthViewModel

    //PARSE DATA
    private val args: HistoryDetailActivityArgs by navArgs()

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityHistoryDetailBinding.inflate(layoutInflater)
        setContentView(binding.root)
        supportActionBar?.hide()

        val factory = FlightViewModelFactory(FlightRepository())
        viewModel = ViewModelProvider(this, factory)[FlightViewModel::class.java]

        val authFactory = AuthViewModelFactory(AuthRepository())
        authViewModel = ViewModelProvider(this, authFactory)[AuthViewModel::class.java]

        //SETUP
        savedPref = SavedPreference(this)
        val tokenFromApi = savedPref.getData(Constants.ACCESS_TOKEN)
        accessToken = "Bearer $tokenFromApi"

        val bookingID = args.bookingid

        observeDetailHistory(bookingID, accessToken)

        binding.ivBack.setOnClickListener {
            finish()
        }

        binding.ivDelete.setOnClickListener {

            val view = View.inflate(this, R.layout.delete_dialog, null)

            AlertDialog.Builder(this, R.style.MyAlertDialogTheme)
                .setView(view)
                .setNegativeButton("No") { p0, _ ->
                    p0.dismiss()
                }
                .setPositiveButton("YES") { _, _ ->
                    alertDelete(bookingID, accessToken)
                }
                .show()
        }

    }

    private fun observeDetailHistory(bookingID: String, accessToken: String) {
        viewModel.detailHistory(bookingID, accessToken).observe(this) { response ->
            if (response is Resources.Loading) {
                progressBar(true)
            } else if (response is Resources.Error) {
                progressBar(false)
                Toast.makeText(this, response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                progressBar(false)
                val result = response.data
                if (result != null) {
                    if (result.status.equals("success")) {
                        binding.fromTv.text = result.data?.departure
                        binding.toTv.text = result.data?.destination
                        binding.tvBookingCode.text = (result.data?.booking_code).toString()

                        Glide
                            .with(this)
                            .load(result.data?.icon)
                            .apply(RequestOptions().override(32, 32))
                            .into(binding.ivDestination)

                        binding.tvAirplaneName.text = result.data?.airline
                        binding.tvDepartTime.text = result.data?.depart_time
                        binding.tvArrivalTime.text = result.data?.arrival_time

                        binding.tvDeparture.text = result.data?.departure
                        binding.tvDestination.text = result.data?.destination

                        if ((result.data?.passenger_title != null) && (result.data.passenger_name != null)) {
                            binding.tvDefaultStatus.text = result.data.passenger_title
                            binding.tvDefaultUser.text = result.data.passenger_name
                        } else {
                            binding.tvDefaultStatus.text = result.data?.passenger_title
                            binding.tvDefaultUser.text = getString(R.string.detail_user_info)
                            binding.tvDefaultUser.setTextColor(
                                ContextCompat.getColor(
                                    this,
                                    R.color.customColorFont
                                )
                            )
                            binding.purchaseStatus.setBackgroundResource(R.color.failed)
                        }

                        //CONVERT PRICE TO RP
                        val myIndonesianLocale = Locale("in", "ID")
                        val numberFormat = NumberFormat.getCurrencyInstance(myIndonesianLocale)
                        numberFormat.maximumFractionDigits = 0
                        val convert = numberFormat.format(result.data?.price)
                        binding.tvPrice.text = convert

                        val status = result.data?.status

                        when (status) {
                            "success" -> {

                                binding.purchaseStatus.setBackgroundResource(R.color.success)

                            }

                            "pending" -> {

                                binding.purchaseStatus.setBackgroundResource(R.color.pending)

                            }

                            "canceled" -> {

                                binding.purchaseStatus.setBackgroundResource(R.color.failed)

                            }
                        }

                        binding.tvStatus.text =
                            String.format(getString(R.string.HistoryAdapter_PurchaseStatus, status))
                        binding.tvStatus.setTextColor(
                            ContextCompat.getColor(
                                this,
                                R.color.customColorFont
                            )
                        )

                    } else {

                        val tokenFromApi = savedPref.getData(Constants.REFRESH_TOKEN)

                        val dataToken = hashMapOf(
                            "refreshToken" to tokenFromApi
                        )

                        observeUpdateToken(dataToken)
                    }
                } else {
                    Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    private fun observeUpdateToken(dataToken: HashMap<String, String?>) {
        authViewModel.updateToken(dataToken).observe(this) { response ->
            if (response is Resources.Loading) {
                progressBar(true)
            } else if (response is Resources.Error) {
                progressBar(false)
                Toast.makeText(this, response.error, Toast.LENGTH_SHORT).show()
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

                        val bookingID = args.bookingid

                        observeDetailHistory(bookingID, accessToken)
                    } else {
                        Log.d("REGIS", result.status.toString())
                    }
                } else {
                    Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    private fun alertDelete(bookingID: String, accessToken: String) {
        viewModel.deleteBookingById(bookingID, accessToken).observe(this) { response ->
            if (response is Resources.Loading) {
                progressBar(true)
            } else if (response is Resources.Error) {
                progressBar(false)
                Toast.makeText(this, response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                val result = response.data
                if (result != null) {
                    if (result.status == "success") {

                        progressBar(false)
                        finish()

                    } else {
                        val dataToken = hashMapOf(
                            "refreshToken" to savedPref.getData(Constants.REFRESH_TOKEN)
                        )
                        observeUpdateToken(dataToken)
                    }
                } else {
                    Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    private fun progressBar(isLoading: Boolean) = with(binding) {
        if (isLoading) {
            this.progressBar.visibility = View.VISIBLE
        } else {
            this.progressBar.visibility = View.GONE
        }
    }
}