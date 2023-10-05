package org.superosystem.traveller.ui.main.fragment

import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.text.Editable
import android.text.TextWatcher
import android.util.Log
import android.util.Patterns
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.Toast
import androidx.activity.result.contract.ActivityResultContracts
import androidx.appcompat.app.AppCompatActivity
import androidx.fragment.app.Fragment
import androidx.lifecycle.ViewModelProvider
import androidx.navigation.findNavController
import com.bumptech.glide.Glide
import okhttp3.MediaType.Companion.toMediaType
import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.MultipartBody
import okhttp3.RequestBody
import okhttp3.RequestBody.Companion.asRequestBody
import okhttp3.RequestBody.Companion.toRequestBody
import org.superosystem.traveller.R
import org.superosystem.traveller.data.preference.SavedPreference
import org.superosystem.traveller.data.repository.AccessProfileRepository
import org.superosystem.traveller.data.repository.AuthRepository
import org.superosystem.traveller.databinding.FragmentEditProfileBinding
import org.superosystem.traveller.utils.Constants
import org.superosystem.traveller.utils.Resources
import org.superosystem.traveller.utils.image.uriToFile
import org.superosystem.traveller.viewmodel.AccessProfileViewModel
import org.superosystem.traveller.viewmodel.AuthViewModel
import org.superosystem.traveller.viewmodel.factory.AccessProfileFactory
import org.superosystem.traveller.viewmodel.factory.AuthViewModelFactory
import java.io.File


class EditProfileFragment : Fragment() {

    // BINDING
    private var _binding: FragmentEditProfileBinding? = null
    private val binding get() = _binding!!

    //API
    private lateinit var authViewModel: AuthViewModel
    private lateinit var viewModel: AccessProfileViewModel

    //SESSION
    private lateinit var savedPreference: SavedPreference
    private lateinit var accessToken: String

    //FILE
    private var getFile: File? = null

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View {
        // Inflate the layout for this fragment
        _binding = FragmentEditProfileBinding.inflate(inflater, container, false)
        return binding.root
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        //CREATE API CONNECTION
        val factory = AccessProfileFactory(AccessProfileRepository())
        viewModel = ViewModelProvider(this, factory)[AccessProfileViewModel::class.java]

        val authFactory = AuthViewModelFactory(AuthRepository())
        authViewModel = ViewModelProvider(this, authFactory)[AuthViewModel::class.java]

        savedPreference = SavedPreference(requireContext())

        val token = savedPreference.getData(Constants.ACCESS_TOKEN)
        accessToken = "Bearer $token"

        setUpButton()

        showDataUser()

        binding.ivBackFragmentEditProfile.setOnClickListener {
            val fragmentProfile =
                EditProfileFragmentDirections.actionEditProfileFragmentToProfileFragment()
            binding.root.findNavController().navigate(fragmentProfile)
        }

        binding.tvUploadImage.setOnClickListener {
            startGallery()
        }

        binding.btnSaveChanges.setOnClickListener {
            update()
        }

    }

    //SHOW DATA
    private fun showDataUser() {
        val username = savedPreference.getData(Constants.USERNAME)
        val email = savedPreference.getData(Constants.EMAIL)
        val profilePicture = savedPreference.getData(Constants.PROFILE_PICTURE)

        //SHOW DATA
        binding.edtUsername.setText(username)
        binding.edtEmail.setText(email)

        Glide.with(this)
            .load(profilePicture)
            .placeholder(R.drawable.avatar)
            .into(binding.ivProfilePicture)

    }

    private fun validateUsername(): String? {
        val name = binding.edtUsername.text.toString()

        if (name.isEmpty()) {
            return getString(R.string.error_empty_message)
        }
        return null
    }

    private fun validateEmail(): String? {
        val email = binding.edtEmail.text.toString()

        if (email.isEmpty()) {
            return getString(R.string.error_empty_message)
        } else if (!Patterns.EMAIL_ADDRESS.matcher(email).matches()) {
            return getString(R.string.error_invalid_email_message)
        }
        return null
    }

    private fun validateButton() {
        binding.tilUsername.helperText = validateUsername()
        binding.tilEmail.helperText = validateEmail()

        val validUsername = binding.tilUsername.helperText == null
        val validEmail = binding.tilEmail.helperText == null

        val username = binding.edtUsername.text.toString()
        val email = binding.edtEmail.text.toString()

        if (validEmail && validUsername) {
            binding.btnSaveChanges.isEnabled =
                (username.isNotEmpty()) &&
                        (email.isNotEmpty())
        } else {
            binding.btnSaveChanges.isEnabled = false
        }
    }

    private fun setUpButton() {
        binding.edtUsername.addTextChangedListener(object : TextWatcher {
            override fun beforeTextChanged(p0: CharSequence?, p1: Int, p2: Int, p3: Int) {
                //DO NOTHING
            }

            override fun onTextChanged(p0: CharSequence?, p1: Int, p2: Int, p3: Int) {
                validateButton()
            }

            override fun afterTextChanged(p0: Editable?) {
                //DO NOTHING
            }
        })

        binding.edtEmail.addTextChangedListener(object : TextWatcher {
            override fun beforeTextChanged(p0: CharSequence?, p1: Int, p2: Int, p3: Int) {
                //DO NOTHING
            }

            override fun onTextChanged(p0: CharSequence?, p1: Int, p2: Int, p3: Int) {
                validateButton()
            }

            override fun afterTextChanged(p0: Editable?) {
                //DO NOTHING
            }
        })
    }

    private fun startGallery() {
        val intent = Intent()
        intent.action = Intent.ACTION_GET_CONTENT
        intent.type = "image/*"
        val chooser = Intent.createChooser(intent, "Choose a Picture")
        launcherIntentGallery.launch(chooser)
    }

    private val launcherIntentGallery = registerForActivityResult(
        ActivityResultContracts.StartActivityForResult()
    ) { result ->
        if (result.resultCode == AppCompatActivity.RESULT_OK) {
            val selectedImg: Uri = result.data?.data as Uri
            val myFile = uriToFile(selectedImg, requireContext())

            getFile = myFile

            Glide.with(this)
                .load(getFile)
                .placeholder(R.drawable.avatar)
                .into(binding.ivProfilePicture)

        }
    }

    private fun update() {
        binding.tilUsername.helperText = validateUsername()
        binding.tilEmail.helperText = validateEmail()

        val validUsername = binding.tilUsername.helperText == null
        val validEmail = binding.tilEmail.helperText == null

        val username = binding.edtUsername.text?.toString()?.trim()
        val email = binding.edtEmail.text?.toString()?.trim()

        if (validUsername && validEmail) {
            val dataLoginUsername = username?.toRequestBody("text/plain".toMediaType())
            val dataLoginEmail = email?.toRequestBody("text/plain".toMediaType())

            if (getFile != null) {
                //api process edit profile with image
                val requestImageFile = getFile!!.asRequestBody("image/jpeg".toMediaTypeOrNull())
                val imageMultipart: MultipartBody.Part = MultipartBody.Part.createFormData(
                    "foto_profile",
                    getFile!!.name,
                    requestImageFile
                )

                observeUpdateProfile(accessToken, dataLoginUsername, dataLoginEmail, imageMultipart)

            } else {
                observeUpdateProfile(accessToken, dataLoginUsername, dataLoginEmail, null)
            }
        }
    }

    private fun observeUpdateProfile(
        accessToken: String,
        dataUserUsername: RequestBody?,
        dataUserEmail: RequestBody?,
        imageMultipart: MultipartBody.Part?
    ) {
        viewModel.updateUser(accessToken, dataUserUsername, dataUserEmail, imageMultipart)
            .observe(viewLifecycleOwner) { response ->
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

                            val username = binding.edtUsername.text?.toString()?.trim()
                            val email = binding.edtEmail.text?.toString()?.trim()

                            if (imageMultipart != null) {
                                val photoProfile = result.data?.imageUri
                                savedPreference.putData(Constants.PROFILE_PICTURE, photoProfile)
                            }

                            savedPreference.putData(Constants.USERNAME, username!!)
                            savedPreference.putData(Constants.EMAIL, email!!)

                            Toast.makeText(requireContext(), "Changes Saved", Toast.LENGTH_SHORT)
                                .show()

                            val toProfileFragment =
                                EditProfileFragmentDirections.actionEditProfileFragmentToProfileFragment()
                            binding.root.findNavController().navigate(toProfileFragment)

                        } else if (result.status.equals("fail") && result.message.equals("Failed to edit email. Email is already in use.")) {

                            Toast.makeText(
                                requireContext(),
                                "Failed to edit email. Email is already in use. Please use another email address",
                                Toast.LENGTH_SHORT
                            ).show()

                        } else if (result.status == null && result.message.equals("Payload content length greater than maximum allowed: 512000")) {

                            Toast.makeText(
                                requireContext(),
                                "File size is too big, please try again to upload image with size below 512kb",
                                Toast.LENGTH_SHORT
                            ).show()

                        } else {
                            val dataToken = hashMapOf(
                                "refreshToken" to savedPreference.getData(Constants.REFRESH_TOKEN)
                            )

                            observeUpdateTokenProfile(
                                dataToken,
                                dataUserUsername,
                                dataUserEmail,
                                imageMultipart
                            )

                        }
                    } else {
                        Toast.makeText(
                            requireContext(),
                            getString(R.string.error),
                            Toast.LENGTH_SHORT
                        ).show()
                    }
                }
            }
    }

    private fun observeUpdateTokenProfile(
        dataToken: HashMap<String, String?>,
        dataUser: RequestBody?,
        dataEmail: RequestBody?,
        imageMultipart: MultipartBody.Part?
    ) {
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
                        savedPreference.putData(Constants.ACCESS_TOKEN, newAccessToken)

                        //get new token
                        val tokenFromAPI = (savedPreference.getData(Constants.ACCESS_TOKEN))
                        val accessToken = "Bearer $tokenFromAPI"

                        observeUpdateProfile(accessToken, dataUser, dataEmail, imageMultipart)

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


    private fun progressBar(isLoading: Boolean) = with(binding) {
        if (isLoading) {
            progressBar.visibility = View.VISIBLE
        } else {
            progressBar.visibility = View.GONE
        }
    }

}