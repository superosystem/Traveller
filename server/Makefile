container-db:
	docker compose -f docker-compose-db.yml up

install:
	npm install

# npm run migrate create "create table users"
migrations:
	npm run migrate create $(name)

# npm run migrate up for db
migrate:
	npm run migrate up

lint-fix:
	npx eslint ./src --fix

test-watch:
	npm run test:watch

run-prod:
	pm2 start npm --name "auth-api" -- run "start"