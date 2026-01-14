### --- Builder stage ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies (git often needed for Go modules)
RUN apk add --no-cache ca-certificates git

# Pre-cache module files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the API binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/api ./cmd/api


### --- Runtime stage ---
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates && update-ca-certificates

# Copy binary from builder
COPY --from=builder /app/bin/api /app/api

# Default environment configuration (override in prod as needed)
ENV SERVER_PORT=":8080" \
    LOCAL_FRONTEND_PORT="3000"

EXPOSE 8080

# Run the API server
CMD ["/app/api"]

