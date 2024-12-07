# Use a lightweight base image
FROM golang:1.22-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the source code and dependencies
COPY go.mod go.sum ./
COPY . .

# Download dependencies
RUN go mod download

# Tiding the deps
RUN go mod tidy

# Build the application
RUN go build -o main ./cmd/app

# Create a new image for the final application
FROM alpine:latest

# Copy the built binary
COPY --from=builder /app/main /app/main

# Expose the port for your application
EXPOSE 8080

# Set the working directory
WORKDIR /app

# Start the application
CMD ["./main"]