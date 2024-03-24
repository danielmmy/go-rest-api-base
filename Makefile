## build and release for docker
up:
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	DOCKER_SCAN_SUGGEST=false docker compose up --build -d 
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"

## run from local machine
run:
	@echo "compiling and running locally..."
	go run ./cmd/api

## run all tests
test:
	go clean -testcache
	go test -cover ./...

## check cover
cover:
	go test -coverprofile /tmp/gorest.data ./...
	go tool cover -html /tmp/gorest.data