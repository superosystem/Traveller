### DEVELOPMENT AND TESTING
run-dev:
	go run ./src/cmd/main.go

mockery-repo:
	mockery --dir=src/domain --name=$(name) --filename=$(filename).go --output=src/domain/mocks/repository --outpkg=mocks

mockery-usecase:
	mockery --dir=src/domain --name=$(name) --filename=$(filename).go --output=src/domain/mocks/usecase --outpkg=mocks

run-test:
	go test -v ./test/...

### DOCKER
run-docker-dev:
	docker compose -f docker-compose-dev.yml up

run-docker:
	docker compose up
