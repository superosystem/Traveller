package org.superosystem.traveller.ui.main.fragment

import android.content.res.Configuration
import android.os.Bundle
import android.util.Log
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.appcompat.app.AlertDialog
import androidx.fragment.app.Fragment
import androidx.lifecycle.ViewModelProvider
import androidx.navigation.findNavController
import androidx.navigation.fragment.findNavController
import androidx.recyclerview.widget.GridLayoutManager
import androidx.recyclerview.widget.LinearLayoutManager
import org.superosystem.traveller.R
import org.superosystem.traveller.data.preference.SavedPreference
import org.superosystem.traveller.data.repository.AuthRepository
import org.superosystem.traveller.data.repository.FlightRepository
import org.superosystem.traveller.databinding.FragmentHistoryBinding
import org.superosystem.traveller.ui.adapter.HistoryAdapter
import org.superosystem.traveller.utils.Constants
import org.superosystem.traveller.utils.Resources
import org.superosystem.traveller.viewmodel.AuthViewModel
import org.superosystem.traveller.viewmodel.FlightViewModel
import org.superosystem.traveller.viewmodel.factory.AuthViewModelFactory
import org.superosystem.traveller.viewmodel.factory.FlightViewModelFactory

class HistoryFragment : Fragment() {

    private var _binding: FragmentHistoryBinding? = null
    private val binding get() = _binding!!

    private lateinit var savedPref: SavedPreference
    private lateinit var accessToken: String

    private lateinit var viewModel: FlightViewModel
    private lateinit var authViewModel: AuthViewModel
    private lateinit var list: HistoryAdapter

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        // Inflate the layout for this fragment
        _binding = FragmentHistoryBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        val factory = FlightViewModelFactory(FlightRepository())
        viewModel = ViewModelProvider(this, factory)[FlightViewModel::class.java]

        val authFactory = AuthViewModelFactory(AuthRepository())
        authViewModel = ViewModelProvider(this, authFactory)[AuthViewModel::class.java]

        //SETUP
        savedPref = SavedPreference(requireContext())
        configRecyclerView()

        list.setOnItemClickListener {
            val bundle = Bundle().apply { putString("bookingid", it.id) }
            findNavController().navigate(
                R.id.action_historyFragment_to_historyDetailActivity,
                bundle
            )
        }

        val tokenFromApi = savedPref.getData(Constants.ACCESS_TOKEN)
        accessToken = "Bearer $tokenFromApi"

        binding.ivDelete.setOnClickListener {

            if (list.itemCount > 0) {
                val viewDialog = View.inflate(context, R.layout.delete_all_dialog, null)

                context?.let {
                    AlertDialog.Builder(it, R.style.MyAlertDialogTheme)
                        .setView(viewDialog)
                        .setNegativeButton("No") { p0, _ ->
                            p0.dismiss()
                        }
                        .setPositiveButton("Yes") { _, _ ->
                            observerDeleteAll(accessToken)
                        }
                        .show()
                }
            } else {
                observerDeleteAll(accessToken)
            }

        }
    }

    override fun onResume() {
        super.onResume()
        observerHistory(accessToken)
    }

    private fun observerHistory(accessToken: String) {
        viewModel.history(accessToken).observe(viewLifecycleOwner) { response ->
            if (response is Resources.Loading) {
                enableProgressBar()
            } else if (response is Resources.Error) {
                disableProgressBar()
                Toast.makeText(requireContext(), response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                disableProgressBar()
                val result = response.data
                if (result != null) {
                    if (result.status.equals("success")) {
                        list.differAsync.submitList(result.data?.bookings)

                        if ((result.data?.bookings)!!.isEmpty()) {
                            binding.containerLl.visibility = View.VISIBLE
                        }
                    } else {
                        println("Result : ${result.status.toString()}")

                        val tokenFromApi = savedPref.getData(Constants.REFRESH_TOKEN)
                        println("refresh token : $tokenFromApi")

                        val dataToken = hashMapOf(
                            "refreshToken" to tokenFromApi
                        )

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
                enableProgressBar()
            } else if (response is Resources.Error) {
                disableProgressBar()
                Toast.makeText(requireContext(), response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                disableProgressBar()
                val result = response.data
                if (result != null) {
                    if (result.status.equals("success")) {
                        val newAccessToken = result.data?.accessToken.toString()
                        println("new access token : $newAccessToken")

                        //save new token
                        savedPref.putData(Constants.ACCESS_TOKEN, newAccessToken)

                        //get new token
                        val tokenFromAPI = (savedPref.getData(Constants.ACCESS_TOKEN))
                        println("token from api : $tokenFromAPI")

                        val accessToken = "Bearer $tokenFromAPI"
                        println("access token : $accessToken")

                        observerHistory(accessToken)
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

    private fun configRecyclerView() {
        list = HistoryAdapter(requireContext())
        binding.rvHistoryTickets.apply {
            adapter = list
            layoutManager =
                if (this.resources.configuration.orientation == Configuration.ORIENTATION_LANDSCAPE) {
                    GridLayoutManager(context, 2)
                } else {
                    LinearLayoutManager(context)
                }
        }
    }

    private fun observerDeleteAll(accessToken: String) {
        viewModel.deleteAllBooking(accessToken).observe(viewLifecycleOwner) { response ->
            if (response is Resources.Loading) {
                enableProgressBar()
            } else if (response is Resources.Error) {
                disableProgressBar()
                Toast.makeText(requireContext(), response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                disableProgressBar()
                val result = response.data
                if (result != null) {
                    if (result.status.equals("success")) {

                        val fragmentHistory = HistoryFragmentDirections.actionHistoryFragmentSelf()
                        binding.root.findNavController().navigate(fragmentHistory)

                    } else {
                        println("Result : ${result.status.toString()}")

                        val tokenFromApi = savedPref.getData(Constants.REFRESH_TOKEN)
                        println("refresh token : $tokenFromApi")

                        val dataToken = hashMapOf(
                            "refreshToken" to tokenFromApi
                        )

                        observeUpdateToken(dataToken)
                    }
                } else {
                    Toast.makeText(requireContext(), getString(R.string.error), Toast.LENGTH_SHORT)
                        .show()
                }
            }
        }
    }

    private fun enableProgressBar() {
        binding.progressBar.visibility = View.VISIBLE
    }

    private fun disableProgressBar() {
        binding.progressBar.visibility = View.INVISIBLE
        binding.containerLl.visibility = View.INVISIBLE
    }
}