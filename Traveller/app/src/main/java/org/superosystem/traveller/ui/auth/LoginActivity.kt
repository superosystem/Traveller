package org.superosystem.traveller.ui.auth

import android.content.Context
import android.content.Intent
import android.graphics.Rect
import android.os.Bundle
import android.text.Editable
import android.text.TextWatcher
import android.util.Log
import android.util.Patterns
import android.view.MotionEvent
import android.view.View
import android.view.inputmethod.InputMethodManager
import android.widget.EditText
import android.widget.Toast
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.ViewModelProvider
import org.superosystem.traveller.ui.home.HomeActivity
import org.superosystem.traveller.R
import org.superosystem.traveller.data.preference.SavedPreference
import org.superosystem.traveller.data.repository.AuthRepository
import org.superosystem.traveller.databinding.ActivityLoginBinding
import org.superosystem.traveller.ui.eula.EulaActivity
import org.superosystem.traveller.ui.onboard.OnBoardingActivity
import org.superosystem.traveller.utils.Constants
import org.superosystem.traveller.utils.Resources
import org.superosystem.traveller.viewmodel.AuthViewModel
import org.superosystem.traveller.viewmodel.AuthViewModelFactory

class LoginActivity : AppCompatActivity(), View.OnClickListener {

    private lateinit var binding: ActivityLoginBinding
    private lateinit var viewModel: AuthViewModel
    private lateinit var savedPref: SavedPreference

    override fun onStart() {
        super.onStart()
        if (savedPref.getInstall(Constants.FIRST_INSTALL)) {
            savedPref.putInstall(Constants.FIRST_INSTALL, false)
            startActivity(Intent(this@LoginActivity, OnBoardingActivity::class.java))
            killActivity()
        }

        if (savedPref.getSession(Constants.IS_LOGIN)) {
            startActivity(Intent(this@LoginActivity, HomeActivity::class.java))
            killActivity()
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityLoginBinding.inflate(layoutInflater)
        setContentView(binding.root)

        //SETUP
        supportActionBar?.hide()
        emailFocusListener()
        setUpButton()
        itemOnClickListener()

        //CREATE API CONNECTION
        val factory = AuthViewModelFactory(AuthRepository())
        viewModel = ViewModelProvider(this, factory)[AuthViewModel::class.java]
        savedPref = SavedPreference(this)
    }

    override fun onClick(p0: View?) {
        when (p0?.id) {
            R.id.btn_login -> {
                //go to flight activity
                enableProgressBar()
                loginForm()
            }

            R.id.sign_up -> {
                //go to sign up activity
                startActivity(Intent(this@LoginActivity, RegistrationActivity::class.java))
            }

            R.id.eula_hyperlink -> {
                //go to eula page
                startActivity(Intent(this@LoginActivity, EulaActivity::class.java))
            }
        }
    }

    private fun itemOnClickListener() {
        binding.btnLogin.setOnClickListener(this)
        binding.signUp.setOnClickListener(this)
        binding.eulaHyperlink.setOnClickListener(this)
    }

    override fun onBackPressed() {
        super.onBackPressed()
        finishAffinity()
    }

    @Suppress("IMPLICIT_BOXING_IN_IDENTITY_EQUALS")
    override fun dispatchTouchEvent(ev: MotionEvent?): Boolean {
        //the focus on edit text will be cleared when user touch anything outside the edittext
        if (ev?.action === MotionEvent.ACTION_DOWN) {
            val v = currentFocus
            if (v is EditText) {
                val outRect = Rect()
                v.getGlobalVisibleRect(outRect)
                if (!outRect.contains(ev.rawX.toInt(), ev.rawY.toInt())) {
                    v.clearFocus()
                    val imm: InputMethodManager =
                        getSystemService(Context.INPUT_METHOD_SERVICE) as InputMethodManager
                    imm.hideSoftInputFromWindow(v.getWindowToken(), 0)
                }
            }
        }
        return super.dispatchTouchEvent(ev)
    }

    private fun loginForm() {
        binding.tilEmail.helperText = validateEmail()
        binding.tilPassword.helperText = validatePassword()

        val validEmail = binding.tilEmail.helperText == null
        val validPassword = binding.tilPassword.helperText == null

        val email = binding.etvEmail.text.toString()
        val password = binding.etvPassword.text.toString()

        if (validEmail && validPassword) {
            //api process
            val dataLogin = hashMapOf(
                "email" to email,
                "password" to password
            )

            observerLogin(dataLogin)

        } else {
            Toast.makeText(this, "Login gagal, password/email salah", Toast.LENGTH_LONG).show()
            disableProgressBar()
            binding.btnLogin.isEnabled = false
        }
    }

    private fun observerLogin(dataLogin: HashMap<String, String>) {
        viewModel.loginUser(dataLogin).observe(this) { response ->
            if (response is Resources.Loading) {
                enableProgressBar()
            } else if (response is Resources.Error) {
                disableProgressBar()
                Toast.makeText(this, response.error, Toast.LENGTH_SHORT).show()
            } else if (response is Resources.Success) {
                disableProgressBar()
                val result = response.data
                if (result != null) {
                    if (result.status.equals("success")) {
                        //saved userid
                        val accessToken = result.data?.accessToken.toString()
                        val refreshToken = result.data?.refreshToken.toString()
                        saveSessionLogin(accessToken, refreshToken)

                        Toast.makeText(this, result.message, Toast.LENGTH_SHORT).show()
                        Log.d("REGIS", result.message.toString())

                        //intent to home directly
                        val intent = Intent(this@LoginActivity, HomeActivity::class.java)
                        startActivity(intent)
                        killActivity()

                    } else {
                        alertUserError(result.message.toString())
                        Log.d("REGIS", result.message.toString())
                    }
                } else {
                    Toast.makeText(this, getString(R.string.error), Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    private fun saveSessionLogin(accessToken: String, refreshToken: String) {
        savedPref.putData(Constants.ACCESS_TOKEN, accessToken)
        savedPref.putData(Constants.REFRESH_TOKEN, refreshToken)
        savedPref.putSession(Constants.IS_LOGIN, true)
    }

    // Listener
    private fun emailFocusListener() {
        binding.etvEmail.setOnFocusChangeListener { _, focused ->
            //WHEN THE EMAIL EDIT TEXT ISN'T FOCUS
            if (!focused) {
                binding.tilEmail.helperText = validateEmail()
            } else {
                passwordFocusListener()
            }
        }
    }

    private fun passwordFocusListener() {
        binding.etvPassword.setOnFocusChangeListener { _, focused ->
            //WHEN THE PASSWORD EDIT TEXT ISN'T FOCUS
            if (!focused) {
                binding.tilPassword.helperText = validatePassword()
            }
        }
    }

    private fun setUpButton() {
        binding.etvEmail.addTextChangedListener(object : TextWatcher {
            override fun beforeTextChanged(p0: CharSequence?, p1: Int, p2: Int, p3: Int) {
                validateButton()
            }

            override fun onTextChanged(p0: CharSequence?, p1: Int, p2: Int, p3: Int) {
                validateButton()
            }

            override fun afterTextChanged(p0: Editable?) {
                validateButton()
            }
        })

        binding.etvPassword.addTextChangedListener(object : TextWatcher {
            override fun beforeTextChanged(p0: CharSequence?, p1: Int, p2: Int, p3: Int) {
                validateButton()
            }

            override fun onTextChanged(p0: CharSequence?, p1: Int, p2: Int, p3: Int) {
                validateButton()
            }

            override fun afterTextChanged(p0: Editable?) {
                validateButton()
            }
        })
    }

    private fun validateButton() {
        val validEmail = binding.tilEmail.helperText == null
        val validPassword = binding.tilPassword.helperText == null

        val email = binding.etvEmail.text.toString()
        val password = binding.etvPassword.text.toString()

        if (validEmail && validPassword) {
            binding.btnLogin.isEnabled =
                (email.isNotEmpty()) &&
                        (password.isNotEmpty())
        }
    }

    private fun validateEmail(): String? {
        val email = binding.etvEmail.text.toString()

        if (email.isEmpty() && email.isBlank()) {
            return getString(R.string.error_empty_message)
        } else if (!Patterns.EMAIL_ADDRESS.matcher(email).matches()) {
            return getString(R.string.error_invalid_email_message)
        }
        return null
    }

    private fun validatePassword(): String? {
        val password = binding.etvPassword.text.toString()

        if (password.isEmpty() && password.isBlank()) {
            return getString(R.string.error_empty_message)
        } else if (password.length < 6) {
            return getString(R.string.error_password_message)
        }
        return null
    }

    private fun enableProgressBar() {
        binding.progressBar.visibility = View.VISIBLE
    }

    private fun disableProgressBar() {
        binding.progressBar.visibility = View.INVISIBLE
    }

    private fun alertUserError(message: String) {
        val view = View.inflate(this, R.layout.error_401_dialog, null)

        AlertDialog.Builder(this, R.style.MyAlertDialogTheme)
            .setTitle(message)
            .setView(view)
            .setPositiveButton("Back to Login") { _, _ ->
                //DO NOTHING
            }.show()
    }

    private fun killActivity() {
        finish()
    }
}