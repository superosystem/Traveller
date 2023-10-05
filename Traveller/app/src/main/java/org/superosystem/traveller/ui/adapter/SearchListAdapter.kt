package org.superosystem.traveller.ui.adapter

import android.content.Context
import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.AsyncListDiffer
import androidx.recyclerview.widget.DiffUtil
import androidx.recyclerview.widget.RecyclerView
import com.bumptech.glide.Glide
import com.bumptech.glide.request.RequestOptions
import org.superosystem.traveller.R
import org.superosystem.traveller.data.model.flight.Flights
import org.superosystem.traveller.data.preference.SavedPreference
import org.superosystem.traveller.databinding.ItemRowFlightBinding
import org.superosystem.traveller.utils.Constants
import java.text.NumberFormat
import java.util.*


class SearchListAdapter(
    val context: Context
) : RecyclerView.Adapter<SearchListAdapter.ListUsersViewHolder>() {

    private lateinit var savedPref: SavedPreference

    private val differCallback = object : DiffUtil.ItemCallback<Flights>() {
        override fun areItemsTheSame(
            oldItem: Flights,
            newItem: Flights
        ): Boolean {
            //COMPARE ID BECAUSE IT'S UNIQUE
            return oldItem.id == newItem.id
        }

        override fun areContentsTheSame(
            oldItem: Flights,
            newItem: Flights
        ): Boolean {
            return oldItem == newItem
        }
    }

    val differAsync = AsyncListDiffer(this, differCallback)

    private lateinit var onItemClickCallback: OnItemClickCallback

    fun setOnItemClickCallback(onItemClickCallback: OnItemClickCallback) {
        this.onItemClickCallback = onItemClickCallback
    }

    inner class ListUsersViewHolder(var binding: ItemRowFlightBinding) : RecyclerView
    .ViewHolder(binding.root)

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ListUsersViewHolder {
        val binding = ItemRowFlightBinding
            .inflate(
                LayoutInflater.from(parent.context),
                parent,
                false
            )
        return ListUsersViewHolder(binding)
    }

    override fun onBindViewHolder(holder: ListUsersViewHolder, position: Int) {
        val data = differAsync.currentList[position]
        savedPref = SavedPreference(context)

        holder.itemView.apply {
            Glide
                .with(context)
                .load(data.icon)
                .apply(RequestOptions().override(32, 32))
                .into(holder.binding.airplaneImageIv)
            holder.binding.timeDepartTv.text = data.depart_time
            holder.binding.timeArriveTv.text = data.arrival_time
            holder.binding.airplaneNameTv.text = data.airline

            val myIndonesianLocale = Locale("in", "ID")
            val numberFormat = NumberFormat.getCurrencyInstance(myIndonesianLocale)
            numberFormat.maximumFractionDigits = 0
            val convert = numberFormat.format(data.price)

            holder.binding.priceTv.text = convert
            holder.binding.cityDepartCodeTv.text = savedPref.getData(Constants.FROM_CODE)
            holder.binding.cityArriveCodeTv.text = savedPref.getData(Constants.TO_CODE)

            val timeDepart = data.depart_time
            val getTimeDepart: List<String> = timeDepart!!.split(":")

            val timeArrive = data.arrival_time
            val getTimeArrive: List<String> = timeArrive!!.split(":")

            val timeDepartToInt = (getTimeDepart[0] + getTimeDepart[1]).toInt()
            val timeArriveToInt = (getTimeArrive[0] + getTimeArrive[1]).toInt()

            val durationTotal = (timeArriveToInt - timeDepartToInt).toString()

            val hour: String
            val minute: String

            if (durationTotal.length == 4) {
                hour = (durationTotal[0]).toString() + (durationTotal[1]).toString()
                minute = (durationTotal[2]).toString() + (durationTotal[3]).toString()
            } else {
                hour = (durationTotal[0]).toString()
                minute = (durationTotal[1]).toString() + (durationTotal[2]).toString()
            }

            holder.binding.flightDurationTv.text = String.format(
                resources.getString(R.string.SearchListAdapter_flightDuration),
                hour,
                minute
            )
            holder.binding.flightTypeTv.text =
                resources.getString(R.string.SearchListAdapter_direct)

            holder.binding.buyBtn.setOnClickListener {
                onItemClickCallback.onItemClicked(data.id)
            }

        }
    }

    override fun getItemCount(): Int {
        return differAsync.currentList.size
    }

    interface OnItemClickCallback {
        fun onItemClicked(id: String?)
    }

}
