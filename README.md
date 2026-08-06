# Microdrive Distributed Storage

Microdrive is a distributed storage system built with Go and microservice architecture.

The project consists of several services that are responsible for different parts of the system. Services communicate with each other using gRPC.

## Services

### Auth Service

Responsible for authentication and user management.

Implemented features:

* User registration
* Login
* JWT access tokens
* Refresh tokens
* Password hashing with bcrypt
* User role checking

Auth service uses PostgreSQL to store users and refresh tokens.

### Storage Service

Responsible for working with stored files.

Main responsibilities:

* File operations
* Storage logic
* Communication with other services

### Payment Gateway

Responsible for payment-related operations.

Currently contains payment logic and integration layer.

## Technologies

* Go
* gRPC
* Protocol Buffers
* PostgreSQL
* JWT
* bcrypt
* Docker

## Communication between services

Services communicate through gRPC.

I chose gRPC because it provides:

* clear contracts between services using `.proto` files;
* generated client and server code;
* structured communication between microservices.

Example:

```
Client
  |
  |
Auth Service
  |
  |
Storage Service
```

## Project structure

```
microdrive-distributed-storage

├── microdrive_auth
├── storage
├── payment_gateway
└── proto
```

## How to run

### 1. Start dependencies

Run Docker containers:

```bash
docker compose up -d
```

### 2. Apply migrations

Example:

```bash
migrate -path migrations -database "postgres://user:password@localhost:5432/dbname" up
```

### 3. Start service

Example for auth service:

```bash
CONFIG_PATH=config/config_local.yaml go run ./cmd/sso
```

Service starts on:

```
localhost:44044
```

## Configuration

Example configuration:

```yaml
env: "local"

storage_path: "host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"

grpc:
  port: 44044
```

## Testing

Run tests:

```bash
go test ./...
```

## Notes

This project is still in development. Some parts of the system can be improved and extended in the future.
