package org.superosystem.traveller.ui.component

import android.content.Context
import android.graphics.*
import android.os.Build
import android.util.AttributeSet
import android.view.View
import androidx.annotation.ColorInt
import org.superosystem.traveller.R

/**
 * This class draws a background with a hole in the middle of it.
 */
class ViewFinderBackground(
    context: Context,
    attrs: AttributeSet?
) : View(context, attrs) {
    private var viewFinderRect: RectF? = null
    private var onDrawListener: (() -> Unit)? = null
    private val cornerRadius = 25f

    fun setViewFinderRect(viewFinderRect: RectF) {
        this.viewFinderRect = viewFinderRect
        requestLayout()
    }

    override fun setBackgroundColor(@ColorInt color: Int) {
        paintBackground.color = color
        requestLayout()
    }

    private val theme = context.theme
    private val attributes =
        theme.obtainStyledAttributes(attrs, R.styleable.ViewFinderBackground, 0, 0)
    private val backgroundColor =
        attributes.getColor(
            R.styleable.ViewFinderBackground_backgroundColor,
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
                resources.getColor(R.color.notFoundBackgroundCamera, theme)
            } else {
                @Suppress("deprecation")
                resources.getColor(R.color.notFoundBackgroundCamera)
            }
        )

    private var paintBackground = Paint(Paint.ANTI_ALIAS_FLAG).apply {
        color = backgroundColor
        style = Paint.Style.FILL
    }

    private val paintWindow = Paint(Paint.ANTI_ALIAS_FLAG).apply {
        xfermode = PorterDuffXfermode(PorterDuff.Mode.CLEAR)
        style = Paint.Style.FILL
    }

    init {
        setLayerType(LAYER_TYPE_HARDWARE, null)
    }

    override fun onDraw(canvas: Canvas) {
        super.onDraw(canvas)
        canvas.drawPaint(paintBackground)

        val viewFinderRect = this.viewFinderRect
        if (viewFinderRect != null) {
            canvas.drawRoundRect(viewFinderRect, cornerRadius, cornerRadius, paintWindow)
        }

        val onDrawListener = this.onDrawListener
        if (onDrawListener != null) {
            onDrawListener()
        }
    }
}