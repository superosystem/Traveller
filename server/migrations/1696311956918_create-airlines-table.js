/* eslint-disable camelcase */

exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable("airlines", {
    id: {
      type: "VARCHAR(50)",
      primaryKey: true,
    },
    airline: {
      type: "VARCHAR(50)",
      notNull: true,
    },
    icon: {
      type: "TEXT",
      notNull: true,
    },
  });
};

exports.down = (pgm) => {
  pgm.dropTable("airlines");
};
