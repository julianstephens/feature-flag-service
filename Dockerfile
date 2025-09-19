# syntax=docker/dockerfile:1
FROM golang:1.24-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 GOOS=linux

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum, download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source
COPY . .

# Build the binary
RUN go build -o feature-flag-api ./cmd/api

# Final minimal image
FROM alpine:latest

# Certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/feature-flag-api .

# Expose API port
EXPOSE 8080

# Set environment variables (can be overwritten at runtime)
ENV STORAGE_ENDPOINT="etcd:2379"
ENV DB_URL=""
ENV JWT_SECRET=""

# Run the service
ENTRYPOINT ["./feature-flag-api"]
