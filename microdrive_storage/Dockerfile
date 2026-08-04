
FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY . .

RUN go build -ldflags="-s -w" -o storage ./cmd/main.go


FROM debian:bookworm-slim 

WORKDIR /app


RUN apt-get update && apt-get install -y \
    postgresql-client \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/storage /app/storage

COPY --from=builder /app/internal/app/watermark.png ./watermark.png


RUN mkdir -p /app/download && chmod 777 /app/download

EXPOSE 55055

CMD [ "/app/storage" ]