const { Pool } = require('pg');
const bcrypt = require('bcrypt');

const AuthenticationError = require('../commons/exceptions/AuthenticationError');
const ClientError = require('../commons/exceptions/ClientError');
const InvariantError = require('../commons/exceptions/InvariantError');
const TokenManager = require('../commons/tokenize/TokenManager');

const pool = new Pool();

const verifyUserCredential = async (email, password) => {
  const query = {
    text: 'SELECT id, password FROM users WHERE email = $1',
    values: [email],
  };

  const result = await pool.query(query);
  if (!result.rows.length) {
    throw new AuthenticationError('The credentials you provided are wrong');
  }

  const { id, password: hashedPassword } = result.rows[0];
  const match = await bcrypt.compare(password, hashedPassword);
  if (!match) {
    throw new AuthenticationError('The credentials you provided are wrong');
  }

  return id;
};

const userEmailExist = async (email) => {
  const query = {
    text: 'SELECT * FROM users WHERE email = $1',
    values: [email],
  };
  const result = await pool.query(query);
  if (!result.rows.length) {
    return false;
  }
  return true;
};

const verifyRefreshToken = async (token) => {
  const query = {
    text: 'SELECT token FROM authentications WHERE token = $1',
    values: [token],
  };

  const result = await pool.query(query);
  if (!result.rows.length) {
    throw new InvariantError('Refresh token invalid');
  }
};

const deleteRefreshToken = async (token) => {
  const query = {
    text: 'DELETE FROM authentications WHERE token = $1',
    values: [token],
  };
  await pool.query(query);
};

const postAuthenticationHandler = async (request, h) => {
  try {
    const { email, password } = request.payload;

    const id = await verifyUserCredential(email, password);

    const accessToken = TokenManager.generateAccessToken({ id });
    const refreshToken = TokenManager.generateRefreshToken({ id });

    const query = {
      text: 'INSERT INTO authentications VALUES($1)',
      values: [refreshToken],
    };
    await pool.query(query);

    const response = h.response({
      status: 'success',
      message: 'Authentication success',
      data: {
        accessToken,
        refreshToken,
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
      message: 'Maaf, terjadi kegagalan pada server kami.',
    });
    response.code(500);
    return response;
  }
};

const putAuthenticationHandler = async (request, h) => {
  try {
    const { refreshToken } = request.payload;
    await verifyRefreshToken(refreshToken);
    const { id } = TokenManager.verifyRefreshToken(refreshToken);

    const accessToken = TokenManager.generateAccessToken({ id });
    return {
      status: 'success',
      message: 'Access Token berhasil diperbarui',
      data: {
        accessToken,
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
      message: 'Maaf, terjadi kegagalan pada server kami.',
    });
    response.code(500);
    return response;
  }
};

const deleteAuthenticationHandler = async (request, h) => {
  try {
    const { refreshToken } = request.payload;
    await verifyRefreshToken(refreshToken);
    await deleteRefreshToken(refreshToken);

    return {
      status: 'success',
      message: 'Authentications has been removed',
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
      message: 'Maaf, terjadi kegagalan pada server kami.',
    });
    response.code(500);
    return response;
  }
};

// Google Auth
// eslint-disable-next-line consistent-return
const getGoogleAuthenticationHandler = async (request, h) => {
  try {
    if (request.auth.isAuthenticated) {
      const user = request.auth.credentials.profile.raw;
      const id = request.auth.credentials.profile.id;

      const accessToken = TokenManager.generateAccessToken({ id });
      const refreshToken = TokenManager.generateRefreshToken({ id });

      const existEmail = await userEmailExist(user.email);
      if (existEmail === false) {
        // Input users
        const queryUser = {
          text: 'INSERT INTO users VALUES($1, $2, $3, $4, $5) RETURNING id',
          values: [id, user.name, user.email, '', user.picture],
        };
        const resultUser = await pool.query(queryUser);
        if (!resultUser.rows.length) {
          throw new InvariantError('Failed to add user');
        }
        // Create Auth
        const queryAuth = {
          text: 'INSERT INTO authentications VALUES($1)',
          values: [refreshToken],
        };
        await pool.query(queryAuth);
        const response = h.response({
          status: 'success',
          message: 'Authentication success',
          data: {
            accessToken,
            refreshToken,
          },
        });
        response.code(201);
        return response;
      }

      // if exist email is true
      // Create Auth
      const query = {
        text: 'INSERT INTO authentications VALUES($1)',
        values: [refreshToken],
      };
      await pool.query(query);
      const response = h.response({
        status: 'success',
        message: 'Authentication success',
        data: {
          accessToken,
          refreshToken,
        },
      });
      response.code(201);
      return response;
    }
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
      message: 'Maaf, terjadi kegagalan pada server kami.',
    });
    response.code(500);
    return response;
  }
};

module.exports = {
  postAuthenticationHandler,
  putAuthenticationHandler,
  deleteAuthenticationHandler,
  getGoogleAuthenticationHandler,
};
