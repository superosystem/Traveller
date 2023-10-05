package org.superosystem.traveller.ui.component

import android.content.Context
import android.graphics.Canvas
import android.graphics.drawable.Drawable
import android.text.Editable
import android.text.TextWatcher
import android.util.AttributeSet
import android.view.MotionEvent
import android.view.View
import androidx.appcompat.widget.AppCompatEditText
import androidx.core.content.ContextCompat
import org.superosystem.traveller.R

class CustomEditText : AppCompatEditText, View.OnTouchListener {
    private lateinit var deleteButtonIcon: Drawable

    //pembuatan constructor karena kita ingin men-extend AppCompatEditText
    constructor(context: Context) : super(context) {
        init()
    }

    constructor(context: Context, attributes: AttributeSet) : super(context, attributes) {
        init()
    }

    constructor(
        context: Context,
        attributes: AttributeSet,
        defStyleAttribute: Int
    )
            : super(context, attributes, defStyleAttribute) {
        init()
    }

    //untuk melakukan customisasi editText
    override fun onDraw(canvas: Canvas?) {
        super.onDraw(canvas)
        textAlignment = View.TEXT_ALIGNMENT_VIEW_START
    }

    //when EditText is being clicked
    override fun onTouch(p0: View?, p1: MotionEvent): Boolean {
        if (compoundDrawables[2] != null) {
            val deleteButtonStart: Float
            val deleteButtonEnd: Float

            //kondisi delete iconnya belum dipencet
            var isDeleteButtonClicked = false

            //pengecekan jenis handphone apakah menggunakan format Right to Left atau tidak
            if (layoutDirection == View.LAYOUT_DIRECTION_RTL) {
                deleteButtonEnd = (deleteButtonIcon.intrinsicWidth + paddingStart).toFloat()
                when {
                    p1.x < deleteButtonEnd -> isDeleteButtonClicked = true
                }
            } else {
                deleteButtonStart = (width - paddingEnd - deleteButtonIcon.intrinsicWidth).toFloat()
                when {
                    p1.x > deleteButtonStart -> isDeleteButtonClicked = true
                }
            }

            if (isDeleteButtonClicked) {
                when (p1.action) {
                    //jika tombol delete di pencet, tombol delete masih akan terlihat
                    MotionEvent.ACTION_DOWN -> {
                        deleteButtonIcon = ContextCompat
                            .getDrawable(context, R.drawable.ic_round_delete) as Drawable
                        showDeleteButton()
                        return true
                    }

                    //jika tombol delete di lepas, tombol delete akan disembunyikan
                    MotionEvent.ACTION_UP -> {
                        deleteButtonIcon = ContextCompat
                            .getDrawable(context, R.drawable.ic_round_delete) as Drawable
                        when {
                            //dan text yang sudah diketik oleh user akan di klik
                            text != null -> text?.clear()
                        }
                        showDeleteButton()
                        return true
                    }

                    else -> return false
                }
            } else {
                return false
            }
        }
        return false
    }

    //method untuk menampilkan icon delete (x)
    private fun showDeleteButton() {
        /* when showing, we put the delete icon to be left of the edit text
        artinya ditaruh di pojok kanan edit textnya */
        setButtonDrawables(endOfTheText = deleteButtonIcon)
    }

    //method untuk menyembunyikan icon delete (x)
    private fun hideDeleteButton() {
        setButtonDrawables()
    }

    //logic for positioning the text
    private fun setButtonDrawables(
        startOfTheText: Drawable? = null,
        topOfTheText: Drawable? = null,
        endOfTheText: Drawable? = null,
        bottomOfTheText: Drawable? = null
    ) {
        //menampilkan gambar pada edit text dengan parameter (left, top, right, bottom)
        setCompoundDrawablesWithIntrinsicBounds(
            startOfTheText,
            topOfTheText,
            endOfTheText,
            bottomOfTheText
        )
    }

    private fun init() {
        deleteButtonIcon = ContextCompat
            .getDrawable(context, R.drawable.ic_round_delete) as Drawable
        setOnTouchListener(this)

        addTextChangedListener(object : TextWatcher {
            override fun beforeTextChanged(p0: CharSequence?, p1: Int, p2: Int, p3: Int) {
                //Do nothing
            }

            override fun onTextChanged(p0: CharSequence?, p1: Int, p2: Int, p3: Int) {
                //condition where the edit text is filled with text or not,
                //if yes, then show the delete icon
                //if not, then hide the delete icon
                if (p0.toString().isNotEmpty()) {
                    showDeleteButton()
                } else {
                    hideDeleteButton()
                }
            }

            override fun afterTextChanged(p0: Editable?) {
                //Do Nothing
            }

        })
    }
}