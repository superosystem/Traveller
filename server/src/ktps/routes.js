const { addImageKtp } = require('./handler');

const routes = [
  {
    method: 'POST',
    path: '/ktp',
    options: {
      payload: {
        output: 'stream',
        parse: true,
        multipart: true,
        maxBytes: 1000 * 1000 * 3,
      },
      auth: {
        strategies: ['traveller_jwt'],
      },
    },
    handler: addImageKtp,
  },
  {
    method: 'PUT',
    path: '/ktp',
    options: {
      payload: {
        output: 'stream',
        parse: true,
        multipart: true,
        maxBytes: 1000 * 1000 * 3,
      },
      auth: {
        strategies: ['traveller_jwt'],
      },
    },
    handler: addImageKtp,
  },
];

module.exports = routes;
