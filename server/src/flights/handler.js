const { nanoid } = require('nanoid');
const { Pool } = require('pg');

const ClientError = require('../commons/exceptions/ClientError');
const InvariantError = require('../commons/exceptions/InvariantError');

const pool = new Pool();

// eslint-disable-next-line no-undef
checkExistBooking = async (userId) => {
  const query = {
    text: 'SELECT * FROM bookings WHERE id_user = $1',
    values: [userId],
  };
  const result = await pool.query(query);

  if (!result.rows[0]) {
    throw new InvariantError('Booking data does not exist');
  }
};

const getFlightsHandler = async (request, h) => {
  const { departure, destination } = request.query;
  const query = {
    text: `SELECT * FROM airlines
        	JOIN flights ON flights.id_airline = airlines.id`,
  };
  const result = await pool.query(query);
  const flights = result.rows;
  let flighFilter = flights;

  // filter data fligths
  if (departure !== undefined && destination !== undefined) {
    flighFilter = flights.filter(
      (flight) => flight.departure.toLowerCase() === departure.toLowerCase()
        && flight.destination.toLowerCase() === destination.toLowerCase(),
    );
  }

  const response = h.response({
    status: 'success',
    data: {
      flights: flighFilter.map((flight) => ({
        id: flight.id,
        airline: flight.airline,
        icon: flight.icon,
        depart_time: flight.depart_time,
        arrival_time: flight.arrival_time,
        departure: flight.departure,
        destination: flight.destination,
        price: flight.price,
      })),
    },
  });
  response.code(200);
  return response;
};

const postFlightBookingHandler = async (request, h) => {
  try {
    const { id: idFlight } = request.payload;
    const { id: idUser } = request.auth.credentials;

    const id = `booking-${nanoid(16)}`;
    const status = 'pending';
    const bookingCode = Math.floor(Math.random() * 1000000000);
    const createdAt = new Date().getTime();
    const updatedAt = createdAt;

    const query = {
      text: 'INSERT INTO bookings VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id',
      values: [id, status, bookingCode, idUser, idFlight, createdAt, updatedAt],
    };
    const result = await pool.query(query);
    if (!result.rows[0].id) {
      throw new InvariantError('Booking failed to add');
    }

    const response = h.response({
      status: 'success',
      message: 'Booking success',
      data: {
        bookingId: id,
      },
    });
    response.code(201);
    return response;
  } catch (error) {
    if (error instanceof ClientError) {
      const response = h.response({
        status: 'fail',
        message: error.message,
      });
      response.code(error.statusCode);
      return response;
    }

    // Server ERROR!
    const response = h.response({
      status: 'error',
      message: 'Sorry, there was a failure on our server.',
    });
    response.code(500);
    return response;
  }
};

const getBookingByUserIdHandler = async (request, h) => {
  try {
    const { id: idUser } = request.auth.credentials;
    const query = {
      text: `SELECT * FROM flights
			 			JOIN bookings ON bookings.id_flight = flights.id 
						WHERE bookings.id_user = $1`,
      values: [idUser],
    };
    const result = await pool.query(query);
    const bookings = result.rows;
    const dateNow = new Date().getTime();

    // eslint-disable-next-line array-callback-return
    bookings.map((booking) => {
      const duration = dateNow - booking.created_at;
      if (booking.status === 'pending' && duration > 300000) {
        // eslint-disable-next-line no-param-reassign
        booking.status = 'canceled';
      }
    });

    return {
      status: 'success',
      data: {
        bookings: bookings.map((booking) => ({
          id: booking.id,
          departure: booking.departure,
          destination: booking.destination,
          booking_code: booking.booking_code,
          price: booking.price,
          status: booking.status,
        })),
      },
    };
  } catch (error) {
    if (error instanceof ClientError) {
      const response = h.response({
        status: 'fail',
        message: error.message,
      });
      response.code(error.statusCode);
      return response;
    }

    // Server ERROR!
    const response = h.response({
      status: 'error',
      message: 'Sorry, there was a failure on our server.',
    });
    response.code(500);
    return response;
  }
};

const getBookingDetailsByBookingIdHandler = async (request, h) => {
  try {
    const { bookingId: idBooking } = request.params;

    const query = {
      text: `SELECT bookings.id, depart_time, arrival_time, departure, destination, status, price, booking_code, created_at, passenger_name, passenger_title, airline, icon 
						FROM flights 
						INNER JOIN bookings ON flights.id=bookings.id_flight 
						INNER JOIN airlines 
						ON flights.id_airline=airlines.id 
						WHERE bookings.id = $1`,
      values: [idBooking],
    };
    const result = await pool.query(query);
    const bookingDetail = result.rows[0];
    const dateNow = new Date().getTime();

    const duration = dateNow - bookingDetail.created_at;
    if (bookingDetail.status === 'pending' && duration > 300000) {
      bookingDetail.status = 'canceled';
    }

    return {
      status: 'success',
      data: bookingDetail,

      // id: bookingDetail.id,
      // departure: bookingDetail.departure,
      // destination: bookingDetail.destination,
      // status: bookingDetail.status,
      // price: bookingDetail.price,
      // booking_code: bookingDetail.booking_code,
      // passenger_name: bookingDetail.passenger_name,
      // passenger_title: bookingDetail.passenger_title,
      // depart_time: bookingDetail.depart_time,
      // arrival_time: bookingDetail.arrival_time,
      // airline: bookingDetail.airline,
      // icon: bookingDetail.icon
    };
  } catch (error) {
    if (error instanceof ClientError) {
      const response = h.response({
        status: 'fail',
        message: error.message,
      });
      response.code(error.statusCode);
      return response;
    }

    // Server ERROR!
    const response = h.response({
      status: 'error',
      message: 'Sorry, there was a failure on our server.',
    });
    response.code(500);
    return response;
  }
};

const putBookingByIdHandler = async (request, h) => {
  try {
    const { id: idBooking } = request.params;
    const { title, name } = request.payload;

    // checking the exist title and name
    if (!title || !name) {
      throw new InvariantError('Please include the title and name');
    }

    const status = 'success';
    const updatedAt = new Date().toISOString();

    const query = {
      text: 'UPDATE bookings SET status = $1, updated_at = $2, passenger_name = $3, passenger_title = $4 WHERE id = $5 RETURNING id',
      values: [status, updatedAt, name, title, idBooking],
    };

    const result = await pool.query(query);
    if (!result.rows[0].id) {
      throw new InvariantError('Booking failed to update');
    }

    const response = h.response({
      status: 'success',
      message: 'Booking updated',
      data: {
        bookingId: idBooking,
      },
    });
    response.code(200);
    return response;
  } catch (error) {
    if (error instanceof ClientError) {
      const response = h.response({
        status: 'fail',
        message: error.message,
      });
      response.code(error.statusCode);
      return response;
    }

    // Server ERROR!
    const response = h.response({
      status: 'error',
      message: 'Sorry, there was a failure on our server.',
    });
    response.code(500);
    return response;
  }
};

const deleteBookingsHandler = async (request, h) => {
  try {
    const { id: idUser } = request.auth.credentials;
    // eslint-disable-next-line no-undef
    await checkExistBooking(idUser);

    const query = {
      text: 'DELETE FROM bookings WHERE id_user = $1 RETURNING id',
      values: [idUser],
    };
    const result = await pool.query(query);

    if (!result.rows.length) {
      throw new InvariantError('All bookings failed to be deleted');
    }

    return {
      status: 'success',
      message: 'Bookings has been deleted',
    };
  } catch (error) {
    if (error instanceof ClientError) {
      const response = h.response({
        status: 'fail',
        message: error.message,
      });
      response.code(error.statusCode);
      return response;
    }

    // Server ERROR!
    const response = h.response({
      status: 'error',
      message: 'Sorry, there was a failure on our server.',
    });
    response.code(500);
    return response;
  }
};
const deleteBookingByIdHandler = async (request, h) => {
  try {
    const { id: idUser } = request.auth.credentials;
    const { id: idBooking } = request.params;

    const query = {
      text: 'DELETE FROM bookings WHERE id = $1 AND id_user = $2 RETURNING id',
      values: [idBooking, idUser],
    };
    const result = await pool.query(query);

    if (!result.rows.length) {
      throw new InvariantError('A booking failed to be deleted');
    }

    return {
      status: 'success',
      message: 'A booking history has been deleted',
    };
  } catch (error) {
    if (error instanceof ClientError) {
      const response = h.response({
        status: 'fail',
        message: error.message,
      });
      response.code(error.statusCode);
      return response;
    }

    // Server ERROR!
    const response = h.response({
      status: 'error',
      message: 'Sorry, there was a failure on our server.',
    });
    response.code(500);
    return response;
  }
};

module.exports = {
  getFlightsHandler,
  postFlightBookingHandler,
  getBookingByUserIdHandler,
  putBookingByIdHandler,
  getBookingDetailsByBookingIdHandler,
  deleteBookingsHandler,
  deleteBookingByIdHandler,
};
