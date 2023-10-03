const { Pool } = require('pg');

const ClientError = require('../commons/exceptions/ClientError');
const InvariantError = require('../commons/exceptions/InvariantError');

const pool = new Pool();

function uppercase(string) {
  return string.toString().toUpperCase();
}

const getKtpResult = async (request, h) => {
  try {
    const { id: idUser } = request.auth.credentials;

    // const queryImageKtp = {
    // 	text: 'SELECT id FROM ktps WHERE id_user = $1',
    // 	values: [id],
    // }

    // const getIdKtp = await pool.query(queryImageKtp);
    // const id_ktp = getIdKtp.rows[0].id ;

    const query = {
      text: 'SELECT title, name, nationality, nik, sex, married FROM ktpresults WHERE id_user = $1',
      values: [idUser],
    };

    const result = await pool.query(query);
    if (!result.rows.length) {
      throw new InvariantError('Failed to get data from ktpresult');
    }

    const dataKtp = result.rows[0];

    if (
      uppercase(dataKtp.sex) !== 'MALE'
      && uppercase(dataKtp.sex) !== 'FEMALE'
    ) {
      if (
        uppercase(dataKtp.sex) === 'PEREMPUAN'
        && uppercase(dataKtp.married) === 'KAWIN'
      ) {
        dataKtp.title = 'Mrs';
        dataKtp.sex = 'Female';
        dataKtp.married = 'Married';
      } else if (
        uppercase(dataKtp.sex) === 'PEREMPUAN'
        && uppercase(dataKtp.married) !== 'KAWIN'
      ) {
        dataKtp.title = 'Ms';
        dataKtp.sex = 'Female';
        dataKtp.married = 'Single';
      } else if (
        uppercase(dataKtp.sex) === 'LAKI-LAKI'
        && uppercase(dataKtp.married) === 'KAWIN'
      ) {
        dataKtp.title = 'Mr';
        dataKtp.sex = 'Male';
        dataKtp.married = 'Married';
      } else {
        dataKtp.title = 'Mr';
        dataKtp.sex = 'Male';
        dataKtp.married = 'Single';
      }
    }

    const response = h.response({
      status: 'success',
      data: dataKtp,
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

const putKtpResult = async (request, h) => {
  try {
    const {
 title, name, nationality, nik, sex, married,
} = request.payload;
    const { id: idUser } = request.auth.credentials;

    // const queryImageKtp = {
    // 	text: 'SELECT id FROM ktps WHERE id_user = $1',
    // 	values: [id],
    // }

    // const getIdKtp = await pool.query(queryImageKtp);
    // const id_ktp = getIdKtp.rows[0].id ;

    const query = {
      text: 'UPDATE ktpresults SET title = $1, name = $2, nationality = $3, nik = $4, sex = $5, married = $6 WHERE id_user = $7 RETURNING title, name, nationality, nik, sex, married',
      values: [title, name, nationality, nik, sex, married, idUser],
    };

    const result = await pool.query(query);
    if (!result.rows.length) {
      throw new InvariantError('Failed update data from KtpResults');
    }

    const dataKtp = result.rows;
    const response = h.response({
      status: 'success',
      data: dataKtp,
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
      message: 'Sorry, our server are busy. Please, try again later.',
    });
    response.code(500);
    return response;
  }
};

module.exports = { getKtpResult, putKtpResult };
