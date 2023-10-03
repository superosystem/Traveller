/* eslint-disable camelcase */

exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable("bookings", {
    id: {
      type: "VARCHAR(50)",
      primaryKey: true,
    },
    status: {
      type: "VARCHAR(10)",
      notNull: true,
    },
    booking_code: {
      type: "INTEGER",
      notNull: true,
    },
    id_user: {
      type: "VARCHAR(50)",
      notNull: true,
    },
    id_flight: {
      type: "VARCHAR(50)",
      notNull: true,
    },
    created_at: {
      type: "TEXT",
      notNull: true,
    },
    updated_at: {
      type: "TEXT",
      notNull: true,
    },
    passenger_name: {
      type: "TEXT",
    },
    passenger_title: {
      type: "VARCHAR(5)",
    },
  });

  /*
    Menambahkan constraint UNIQUE, kombinasi dari kolom id_user dan id_flight.
    Guna menghindari duplikasi data antara nilai keduanya.
  */
  pgm.addConstraint(
    "bookings",
    "unique_id_user_and_id_flight",
    "UNIQUE(id_user, id_flight)"
  );

  // memberikan constraint foreign key pada kolom id_user dan id_flight terhadap users.id dan flights.id
  pgm.addConstraint(
    "bookings",
    "fk_bookings.id_user_users.id",
    "FOREIGN KEY(id_user) REFERENCES users(id) ON DELETE CASCADE"
  );
  pgm.addConstraint(
    "bookings",
    "fk_bookings.id_flight_flights.id",
    "FOREIGN KEY(id_flight) REFERENCES flights(id) ON DELETE CASCADE"
  );
};

exports.down = (pgm) => {
  pgm.dropTable("bookings");
};
