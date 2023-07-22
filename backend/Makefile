################################################################################
### GO TOOLS
################################################################################
run-dev:
	go run ./src/cmd/main.go

mockery-repo:
	mockery --dir=src/domain --name=$(name) --filename=$(filename).go --output=src/domain/mocks/repository --outpkg=mocks

mockery-usecase:
	mockery --dir=src/domain --name=$(name) --filename=$(filename).go --output=src/domain/mocks/usecase --outpkg=mocks

run-test:
	go test -v ./...

################################################################################
### DOCKER
################################################################################
container-db:
	docker compose -f docker-compose-db.yml up

image:
	docker bui -t training-system-backend .

container:
	docker compose up
