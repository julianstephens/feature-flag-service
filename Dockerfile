FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

# Certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

COPY . .

CMD ["air", "-c", ".air.toml"]
