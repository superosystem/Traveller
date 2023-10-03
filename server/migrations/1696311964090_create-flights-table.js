/* eslint-disable camelcase */

exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable("flights", {
    id: {
      type: "VARCHAR(50)",
      primaryKey: true,
    },
    depart_time: {
      type: "TEXT",
      notNull: true,
    },
    arrival_time: {
      type: "TEXT",
      notNull: true,
    },
    departure: {
      type: "VARCHAR(50)",
      notNull: true,
    },
    destination: {
      type: "VARCHAR(50)",
      notNull: true,
    },
    price: {
      type: "INTEGER",
      notNull: true,
    },
    id_airline: {
      type: "VARCHAR(50)",
      notNull: true,
    },
  });

  pgm.addConstraint(
    "flights",
    "fk_flights.id_airline_airlines.id",
    "FOREIGN KEY(id_airline) REFERENCES airlines(id) ON DELETE CASCADE"
  );
};

exports.down = (pgm) => {
  pgm.dropTable("flights");
};
