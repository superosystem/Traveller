/* eslint-disable camelcase */

exports.shorthands = undefined;

exports.up = (pgm) => {
  pgm.createTable("ktps", {
    id: {
      type: "VARCHAR(50)",
      primaryKey: true,
    },
    image_url: {
      type: "TEXT",
      notNull: true,
    },
    id_user: {
      type: "VARCHAR(50)",
      notNull: true,
    },
  });

  pgm.addConstraint(
    "ktps",
    "fk_ktps.id_user_users.id",
    "FOREIGN KEY(id_user) REFERENCES users(id) ON DELETE CASCADE"
  );
};

exports.down = (pgm) => {
  pgm.dropTable("ktps");
};
