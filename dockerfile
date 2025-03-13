# File: Dockerfile
# Build stage
FROM golang:1.19-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blockchain-client .

# Final stage
FROM alpine:3.17

# Install certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/blockchain-client .

# Copy the .env file if it exists (optional in production)
COPY --from=builder /app/.env* ./ 2>/dev/null || true

# Expose the application port
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release
ENV PORT=8080
ENV RPC_URL=https://polygon-rpc.com/
ENV TIMEOUT_SECONDS=10

# Run the application
CMD ["./blockchain-client"]