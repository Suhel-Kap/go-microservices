FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerServiceApp
MAIL_BINARY=mailServiceApp
LISTENER_BINARY=listenerApp
FRONT_BINARY=frontEndApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_broker build_auth build_logger build_mailer build_listener
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_front_linux: builds the front end binary as a linux executable
build_front_linux:
	@echo "Building front end linyx binary..."
	cd ./front-end && env GOOS=linux CGO_ENABLED=0 go build -o ${FRONT_BINARY} ./cmd/web
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api
	@echo "Done!"

## build_auth: builds the auth binary as a linux executable
build_auth:
		@echo "Building auth binary..."
		cd ./auth-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
		@echo "Done!"

## build_logger: builds the logger binary as a linux executable
build_logger:
		@echo "Building logger binary..."
		cd ./logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
		@echo "Done!"

## build_mailer: builds the mailer binary as a linux executable
build_mailer:
		@echo "Building mailer binary..."
		cd ./mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
		@echo "Done!"

## build_listener: builds the listener binary as a linux executable
build_listener:
		@echo "Building listener binary..."
		cd ./listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} .
		@echo "Done!"

## build_front: builds the frone end binary
build_front:
	@echo "Building front end binary..."
	cd ./front-end && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ./front-end && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"

## proto_log: generates the proto file for the logger service
proto_log:
	@echo "Generating proto file for logger service..."
	cd ./logger-service/logs && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative logs.proto
	@echo "Done!"

## proto_log_broker: generates the proto file for the broker service
proto_log_broker:
	@echo "Generating proto file for broker service..."
	cd ./broker-service/logs && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative logs.proto
	@echo "Done!"
