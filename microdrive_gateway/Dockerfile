FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o gateway-handler ./cmd/main.go

FROM debian:bookworm-slim 

WORKDIR /app


COPY --from=builder /app/gateway-handler /app/gateway-handler

EXPOSE 55056

CMD [ "/app/gateway-handler" ]