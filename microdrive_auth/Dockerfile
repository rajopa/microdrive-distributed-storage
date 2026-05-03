FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o auth-service ./cmd/sso/main.go

FROM debian:bookworm-slim 

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates postgresql-client && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/auth-service /app/auth-service
COPY --from=builder /app/config /app/config
COPY --from=builder /app/migrations /app/migrations


EXPOSE 44044


CMD [ "/app/auth-service" ]