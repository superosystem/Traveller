// Local storage
const fs = require('fs');
const path = require('path');
const { v4: uuidv4 } = require('uuid');

const InvariantError = require('../exceptions/InvariantError');

class StorageService {
	constructor(folder) {
		this._folder = folder;

		if (!fs.existsSync(folder)) {
			fs.mkdirSync(folder, { recursive: true });
		}
	}

	deleteFile() {
		const pathFile = `${this._folder}`;
		// file removed from local storage
		fs.readdir(pathFile, (err, files) => {
			if (err) throw err;

			// eslint-disable-next-line no-restricted-syntax
			for (const file of files) {
				fs.unlink(path.join(pathFile, file), (error) => {
					if (err) throw error;
				});
			}
		});
	}

	writeFile(file, meta) {
		// image extension validation
		const ext = path.extname(meta.filename);
		const validExt = ['.jpg', '.png', '.jpeg'];
		if (validExt.indexOf(ext) === -1) {
			throw new InvariantError('Not allowed file type');
		}

		// custom name from date + filename
		const filename = +new Date() + meta.filename;
		const pathFile = `${this._folder}/${filename}`;

		const fileStream = fs.createWriteStream(pathFile);

		return new Promise((resolve, reject) => {
			fileStream.on('error', (error) => reject(error));
			file.pipe(fileStream);
			file.on('end', () => resolve(filename));
		});
	}

	uploadFile(filename) {
		// eslint-disable-next-line no-shadow
		const path = `${this._folder}/${filename}`;
		// eslint-disable-next-line no-undef
		const storage = storageRef.upload(path, {
			public: true,
			destination: `profile/${filename}`,
			metadata: {
				metadata: {
					firebaseStorageDownloadTokens: uuidv4(),
				},
			},
		});
		const imageLocation = `${process.env.FIREBASE_STORAGE_NAME}/ktpimage/${filename}`;
		return imageLocation;
	}
}

module.exports = StorageService;
