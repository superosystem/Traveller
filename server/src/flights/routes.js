const {
  getFlightsHandler,
  postFlightBookingHandler,
  getBookingByUserIdHandler,
  putBookingByIdHandler,
  getBookingDetailsByBookingIdHandler,
  deleteBookingsHandler,
  deleteBookingByIdHandler,
} = require('./handler');

const routes = [
  {
    method: 'GET',
    path: '/flights',
    options: {
      auth: {
        strategy: 'traveller_jwt',
      },
      handler: getFlightsHandler,
    },
  },
  {
    method: 'POST',
    path: '/flights/booking',
    options: {
      auth: {
        strategy: 'traveller_jwt',
      },
      handler: postFlightBookingHandler,
    },
  },
  {
    method: 'GET',
    path: '/flights/booking',
    options: {
      auth: {
        strategy: 'traveller_jwt',
      },
      handler: getBookingByUserIdHandler,
    },
  },
  {
    method: 'GET',
    path: '/flights/booking/{bookingId}',
    options: {
      auth: {
        strategy: 'traveller_jwt',
      },
      handler: getBookingDetailsByBookingIdHandler,
    },
  },
  {
    method: 'PUT',
    path: '/flights/booking/{id}',
    options: {
      auth: {
        strategy: 'traveller_jwt',
      },
      handler: putBookingByIdHandler,
    },
  },
  {
    method: 'DELETE',
    path: '/flights/booking',
    options: {
      auth: {
        strategy: 'traveller_jwt',
      },
      handler: deleteBookingsHandler,
    },
  },
  {
    method: 'DELETE',
    path: '/flights/booking/{id}',
    options: {
      auth: {
        strategy: 'traveller_jwt',
      },
      handler: deleteBookingByIdHandler,
    },
  },
];

module.exports = routes;
