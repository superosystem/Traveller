/* eslint-disable camelcase */

exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable("ktpresults", {
    id: {
      type: "VARCHAR(50)",
      primaryKey: true,
    },
    title: {
      type: "VARCHAR(5)",
      notNull: false,
    },
    name: {
      type: "TEXT",
      notNull: true,
    },
    nationality: {
      type: "TEXT",
      notNull: true,
    },
    nik: {
      type: "BIGINT",
      notNull: true,
    },
    sex: {
      type: "VARCHAR(20)",
      notNull: true,
    },
    married: {
      type: "VARCHAR(20)",
      notNull: false,
    },
    id_user: {
      type: "VARCHAR(50)",
      notNull: true,
    },
  });

  pgm.addConstraint(
    "ktpresults",
    "fk_ktpresults.id_user_users.id",
    "FOREIGN KEY(id_user) REFERENCES users(id) ON DELETE CASCADE"
  );
};

exports.down = (pgm) => {
  pgm.dropTable("ktpresults");
};
