const firebaseAdmin = require('firebase-admin');
const { v4: uuidv4 } = require('uuid');
const fs = require('fs');
const { Pool } = require('pg');
const { nanoid } = require('nanoid');
const path = require('path');
const axios = require('axios');

const ClientError = require('../commons/exceptions/ClientError');
const InvariantError = require('../commons/exceptions/InvariantError');

const pool = new Pool();

// Import creations.json from firebase
// eslint-disable-next-line import/no-unresolved
const serviceAccount = require('../../serviceAccountKey.json');

// Initialize App
const admin = firebaseAdmin.initializeApp({
	credential: firebaseAdmin.credential.cert(serviceAccount),
});

// Set the bucket
const storageRef = admin.storage().bucket(process.env.FIREBASE_STORAGE_NAME);

let allName;

// Function to upload and store the file in firebase storage
async function uploadFile(_path, filename) {
	// THIS WILL RETURN THE IMAGE LOCATION
	const imageLocation = `${process.env.FIREBASE_STORAGE_NAME}/ktpimage/${filename}`;
	return imageLocation;
}

// Function for storing the file locally before uploading it to the firebase storage
async function storeFileUpload(file) {
	// **nama file asli
	const { filename } = file.hapi;

	//* * image extension validation
	const ext = path.extname(filename);
	const validExt = ['.jpg', '.png', '.jpeg'];

	if (validExt.indexOf(ext) === -1) {
		throw new InvariantError('Not allowed file type');
	}

	// **file custom name
	const filenameCustom = `${allName}ktp${ext}`;
	const data = file._data;
	const ktpFolder = './ktp';

	// create the ktp folder if doesnt exist
	if (!fs.existsSync(ktpFolder)) {
		fs.mkdirSync(ktpFolder);
	}

	// eslint-disable-next-line consistent-return
	fs.writeFile(`./ktp/${filenameCustom}`, data, (err) => {
		if (err) {
			return err;
		}
	});

	// call uploadFile function
	const imagePath = `./ktp/${filenameCustom}`;
	return uploadFile(imagePath, filenameCustom);
}

// Function for deleting the previous files
async function deletePrevFile(imageName, jsonName) {
	const path1 = `./ktp/${imageName}`;
	const path2 = `./ktp/${jsonName}`;
	// file removed from local storage
	try {
		fs.unlinkSync(path1);
		fs.unlinkSync(path2);
	} catch (error) {
		throw new InvariantError('Failed to delete the file locally');
	}

	// delete image from firebase storage
	try {
		await storageRef.file(`ktpimage/${imageName}`).delete();
		await storageRef.file(`ktpimage/${jsonName}`).delete();
	} catch (error) {
		// console.error(err);
		throw new InvariantError('Failed to delete the files in the server');
	}
}

// Function for getting and writing coordinates into file.json
async function writeCoordinates(dataClassString, imageUrl) {
	const dataClassObject = JSON.parse(dataClassString);
	dataClassObject.image = imageUrl;
	const newDataClassString = JSON.stringify(dataClassObject, null, 4);

	const filenameCustom = `${allName}ktp.json`;

	fs.writeFileSync(`./ktp/${filenameCustom}`, newDataClassString);

	const jsonPath = `./ktp/${filenameCustom}`;
	uploadFile(jsonPath, filenameCustom);
}

// handler function POST ktp
const addImageKtp = async (request, h) => {
	try {
		const { payload } = request;
		const id = nanoid(16);
		const { id: idUser } = request.auth.credentials;

		allName = Date.now();

		// BUAT HAPUS DATA KTP DI DATABASE --------------------------------
		const queryCheckRowResults = {
			text: 'SELECT id FROM ktpresults WHERE id_user = $1',
			values: [idUser],
		};
		const checkRow = await pool.query(queryCheckRowResults);

		// Checking the ktps tabel rows
		if (Object.keys(checkRow.rows).length !== 0) {
			// Delete
			const queryDeleteResult = {
				text: 'DELETE from ktpresults WHERE id_user = $1',
				values: [idUser],
			};
			await pool.query(queryDeleteResult);
			// eslint-disable-next-line no-console
			console.log('Udah dihapus data ktp result-nya');
		}
		// --------------------------------------------------

		const imageUrl = await storeFileUpload(payload.file);

		// fill the database
		const query = {
			text: 'INSERT INTO ktps VALUES($1, $2, $3) RETURNING id',
			values: [id, imageUrl, idUser],
		};

		const result = await pool.query(query);
		if (!result.rows.length) {
			throw new InvariantError('Failed to add KTP image to Database');
		}

		// console.log((fileku.hapi).filename);
		const forpyFilename = payload.file.hapi.filename;
		const ext = path.extname(forpyFilename);

		const filenameCustom = `${allName}ktp${ext}`;

		const pyFilename = `testing/${filenameCustom}`;

		writeCoordinates(payload.data, pyFilename);
		// writeCoordinates(payload.data, imageUrl);

		// PAKAI AXIOS
		// await axios.post('http://localhost:5000/', {filenameCustom})
		await axios
			.post(process.env.ML_SERVER, { filenameCustom })

			.then((res) => {
				// console.log(`Status: ${res.status}`);
				// console.log('Body: ', res.data);

				const id_ktpresult = nanoid(16);
				const title = 'mr';

				const queryKtpR = {
					text: 'INSERT INTO ktpresults VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *',
					values: [
						id_ktpresult,
						title,
						res.data.name,
						res.data.nationality,
						res.data.nik,
						res.data.sex,
						res.data.married,
						idUser,
					],
				};
				pool.query(queryKtpR);
			})
			.catch(() => {
				// console.log(error.response);
				// console.log("Terjadi Error!")
				throw new InvariantError('Failed to get data from model');
			});
		// -----------------------------------------------------------------------

		// Deleting image

		const queryGet = {
			text: 'SELECT image_url FROM ktps WHERE id_user = $1',
			values: [idUser],
		};

		const getKtpUrl = await pool.query(queryGet);
		const imageName = await getKtpUrl.rows[0].image_url.substr(40);
		const jsonName = `${await getKtpUrl.rows[0].image_url.slice(40, 56)}.json`;

		await deletePrevFile(imageName, jsonName);

		// If Delete
		const queryDelete = {
			text: 'DELETE from ktps WHERE id_user = $1',
			values: [idUser],
		};
		await pool.query(queryDelete);
		// console.log("Udah dihapus row-nya");

		const imageId = result.rows[0].id;
		const response = h.response({
			status: 'Success',
			message: 'KTP image successfully added',
			data: {
				imageId,
			},
		});
		response.code(201);
		return response;
	} catch (error) {
		if (error instanceof ClientError) {
			const response = h.response({
				status: 'failed',
				message: error.message,
			});
			response.code(error.statusCode);
			return response;
		}

		// Server error
		const response = h.response({
			status: 'error',
			message: 'Sorry, there was a failure on our server.',
		});
		response.code(500);
		return response;
	}
};
module.exports = { addImageKtp };
