package org.superosystem.traveller.ui.ocr

import android.content.Intent
import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import org.superosystem.traveller.databinding.ActivitySuccessPageBinding
import org.superosystem.traveller.ui.main.HomeActivity
import org.superosystem.traveller.ui.main.HomeActivityTwo

class SuccessPageActivity : AppCompatActivity() {
    //BINDING
    private lateinit var binding: ActivitySuccessPageBinding

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivitySuccessPageBinding.inflate(layoutInflater)
        setContentView(binding.root)

        supportActionBar?.hide()

        binding.cvBtnSuccessToFlight.setOnClickListener {

            val intent = Intent(this, HomeActivity::class.java)
            intent.flags = Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK
            startActivity(intent)

        }

        binding.cvBtnSuccessToHistory.setOnClickListener {

            val intent = Intent(this, HomeActivityTwo::class.java)
            intent.flags = Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK
            startActivity(intent)

        }

    }

    override fun onBackPressed() {
        val intent = Intent(this, HomeActivity::class.java)
        intent.flags = Intent.FLAG_ACTIVITY_NEW_TASK or Intent.FLAG_ACTIVITY_CLEAR_TASK
        startActivity(intent)
    }
}