const {
  postAuthenticationHandler,
  putAuthenticationHandler,
  deleteAuthenticationHandler,
  getGoogleAuthenticationHandler,
} = require('./handler');

const routes = [
  {
    method: 'POST',
    path: '/authentications',
    handler: postAuthenticationHandler,
  },
  {
    method: 'PUT',
    path: '/authentications',
    handler: putAuthenticationHandler,
  },
  {
    method: 'DELETE',
    path: '/authentications',
    handler: deleteAuthenticationHandler,
  },
  {
    method: 'GET',
    path: '/auth/google',
	handler: getGoogleAuthenticationHandler,
  },
];

module.exports = routes;
