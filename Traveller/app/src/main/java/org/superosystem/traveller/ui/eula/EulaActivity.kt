package org.superosystem.traveller.ui.eula

import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import org.superosystem.traveller.R

class EulaActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_eula)
        supportActionBar?.hide()
    }
}