package org.superosystem.traveller.utils

import android.app.Activity
import android.app.AlertDialog
import org.superosystem.traveller.R

class LoadingOCRResultDialog(val activity: Activity) {

    private var dialog: AlertDialog? = null

    fun startLoadingDialog() {
        val builder = AlertDialog.Builder(activity)
        val inflater = activity.layoutInflater

        builder.setView(
            inflater.inflate(
                R.layout.please_wait_while_app_submitting_booking_dialog,
                null
            )
        )
        builder.setCancelable(false)

        dialog = builder.create()
        dialog?.show()
    }

    fun dismissLoadingDialog() {
        if (dialog != null) {
            dialog?.dismiss()
        }
    }
}