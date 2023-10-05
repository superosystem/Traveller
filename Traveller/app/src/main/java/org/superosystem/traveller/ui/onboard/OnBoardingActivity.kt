package org.superosystem.traveller.ui.onboard

import android.content.Intent
import android.os.Bundle
import android.view.View
import androidx.appcompat.app.AppCompatActivity
import org.superosystem.traveller.R
import org.superosystem.traveller.databinding.ActivityOnBoardingBinding
import org.superosystem.traveller.ui.auth.RegistrationActivity

class OnBoardingActivity : AppCompatActivity(), View.OnClickListener {
    //BINDING
    private lateinit var binding: ActivityOnBoardingBinding

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        binding = ActivityOnBoardingBinding.inflate(layoutInflater)
        setContentView(binding.root)

        //SETUP
        supportActionBar?.hide()

        binding.btnRegister.setOnClickListener(this)
    }

    override fun onClick(v: View?) {
        when (v?.id) {
            R.id.btn_register -> {
                //go to register activity
                startActivity(Intent(this@OnBoardingActivity, RegistrationActivity::class.java))
                finish()
            }
        }
    }
}