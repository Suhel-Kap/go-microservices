# Microservices with Golang

This a learning project to understand how to build microservices using Golang.

## Services

- [x] **Broker Service**: This service is responsible for handling the incoming requests and routing them to the appropriate service.
- [x] **Auth Service**: This service is responsible for handling the authentication and authorization of the users. It uses PostgreSQL as the database.
- [x] **Logger Service**: This service is responsible for handling the logging of the requests and responses in a MongoDB database.
- [x] **Mail Service**: This service is responsible for managing an email server built using MailHog.
- [x] **Listener Service**: This service is responsible for listening to the RabbitMQ queue and processing the messages.

## Technologies

- [x] **Golang**: The primary programming language used to build the services.
- [x] **Docker**: The services are containerized using Docker.
- [x] **Docker Compose**: The services are orchestrated using Docker Compose.
- [x] **PostgreSQL**: The database used for the Auth Service.
- [x] **MongoDB**: The database used for the Logger Service.
- [x] **RabbitMQ**: The message broker used for the Listener Service.
- [x] **MailHog**: The email server used for the Mail Service.
- [x] **RPC**: The communication protocol used between the services.
- [x] **gRPC**: The framework used for optimizing the RPC communication.
- [x] **Protobuf**: The serialization format used for the gRPC communication.

## Setup

1. Clone the repository:

```bash
git clone github.com/suhel-kap/go-microservices
```

2. Everything is containerized using Docker. So, you need to have Docker installed on your machine.

3. Use the `make` command to build the images and start the services:

```bash
make up_build
```

4. To start the frontend, use the following command:

```bash
make start
```

5. To stop the services, use the following command:

```bash
make down
```
