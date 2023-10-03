const { nanoid } = require('nanoid');
const { Pool } = require('pg');
const bcrypt = require('bcrypt');
const path = require('path');

const ClientError = require('../commons/exceptions/ClientError');
const InvariantError = require('../commons/exceptions/InvariantError');
const NotFoundError = require('../commons/exceptions/NotFoundError');
const StorageService = require('../commons/storage/StorageService');

const pool = new Pool();

const verifyEmail = async (email, message, exceptions, userId = null) => {
  let result;

  if (exceptions === true) {
    const query = {
      text: `SELECT email FROM users where email = $1 
						EXCEPT
						SELECT email from users WHERE id = $2`,
      values: [email, userId],
    };
    result = await pool.query(query);
  } else {
    const query = {
      text: 'SELECT email FROM users WHERE email = $1',
      values: [email],
    };
    result = await pool.query(query);
  }

  if (result.rows.length > 0) {
    throw new InvariantError(
      `Failed to ${message} email. Email is already in use.`,
    );
  }
};

const postUserHandler = async (request, h) => {
  try {
    const {
 name, email, password, profilePicture,
} = request.payload;

    await verifyEmail(email, 'add', false);

    const id = `users-${nanoid(16)}`;
    const hashedPassword = await bcrypt.hash(password, 10);

    const query = {
      text: 'INSERT INTO users VALUES($1, $2, $3, $4, $5) RETURNING id',
      values: [id, name, email, hashedPassword, profilePicture],
    };

    const result = await pool.query(query);
    if (!result.rows.length) {
      throw new InvariantError('Failed to add user');
    }

    const userId = result.rows[0].id;
    const response = h.response({
      status: 'success',
      message: 'User added successfully',
      data: {
        userId,
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

const getUserByAuthHandler = async (request, h) => {
  try {
    const { id } = request.auth.credentials;

    const query = {
      text: 'SELECT * FROM users WHERE id = $1',
      values: [id],
    };
    const result = await pool.query(query);
    if (!result.rows.length) {
      throw new NotFoundError('User not found');
    }

    const dataUser = result.rows[0];

    return {
      status: 'success',
      data: { user: dataUser },
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

const putUserByAuthHandler = async (request, h) => {
  try {
    // call class of StorageService to object variable
    const { id } = request.auth.credentials;
    const { name, email, profile_picture: profilePicture } = request.payload;
    await verifyEmail(email, 'edit', true, id);

    const storage = new StorageService(path.resolve(__dirname, `images/${id}`));

    // eslint-disable-next-line no-console
    console.log(path.resolve(__dirname, 'images'));

    if (profilePicture) {
      // delete all file in folder users
      await storage.deleteFile('');

      // save file to the local storage
      const filename = await storage.writeFile(
        profilePicture,
        profilePicture.hapi,
      );

      const profileImgUrl = `${process.env.ML_SERVER}/users/images/${id}/${filename}`;

      const query = {
        text: 'UPDATE users SET name = $1, email = $2, profile_picture = $3 WHERE id = $4 RETURNING id',
        values: [name, email, profileImgUrl, id],
      };
      await pool.query(query);

      const response = h.response({
        status: 'success',
        message: 'Success updated profile',
        data: {
          imageUri: profileImgUrl,
        },
      });
      response.code(200);
      return response;
    }

    const query = {
      text: 'UPDATE users SET name = $1, email = $2 WHERE id = $3 RETURNING id',
      values: [name, email, id],
    };
    await pool.query(query);

    const response = h.response({
      status: 'success',
      message: 'Success updated profile',
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

module.exports = {
  postUserHandler,
  getUserByAuthHandler,
  putUserByAuthHandler,
};
