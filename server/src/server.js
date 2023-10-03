require('dotenv').config();

const Hapi = require('@hapi/hapi');
const Jwt = require('@hapi/jwt');
const Bell = require('@hapi/bell');
// eslint-disable-next-line import/no-extraneous-dependencies
const Inert = require('@hapi/inert');

const authentications = require('./authentications/routes');
const users = require('./users/routes');
const ktps = require('./ktps/routes');
const ktpresults = require('./ktpresults/routes');
const flights = require('./flights/routes');

const init = async () => {
    const server = Hapi.server({
        port: process.env.PORT || 8080,
		host: process.env.HOST || 'localhost',
		routes: {
			cors: {
				origin: ['*'],
			},
		},
    });

    await server.register([
        {
            plugin: Jwt,
        },
        {
            plugin: Bell,
        },
        {
            plugin: Inert,
        },
    ]);

    server.auth.strategy('traveller_jwt', 'jwt', {
		keys: process.env.ACCESS_TOKEN_KEY,
		verify: {
			aud: false,
			iss: false,
			sub: false,
			maxAgeSec: process.env.ACCESS_TOKEN_AGE,
		},
		validate: (artifacts) => ({
			isValid: true,
			credentials: {
				id: artifacts.decoded.payload.id,
			},
		}),
	});

    server.route(authentications);
    server.route(users);
    server.route(ktps);
    server.route(ktpresults);
    server.route(flights);

    await server.start();
    // eslint-disable-next-line no-console
    console.log(`Server is running on ${server.info.uri}`);
};

init();
