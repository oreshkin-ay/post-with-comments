# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code to the working directory
COPY . .

# Build the Go binary
RUN go build -o /go-app ./cmd/server

# Stage 2: Create a small, final image
FROM alpine:latest

# Set the working directory in the final image
WORKDIR /root/

# Copy the Go binary from the builder stage
COPY --from=builder /go-app .

# Copy migrations to the final image
COPY --from=builder /app/internal/pkg/db/migrations /root/migrations

# Expose the application port
EXPOSE 8080

# Command to run the Go application
CMD ["./go-app"]
