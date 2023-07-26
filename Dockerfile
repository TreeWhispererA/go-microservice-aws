# Use an official Golang runtime as a parent image
FROM golang:1.16-alpine AS builder

# Set the working directory to /app
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the application code to the container
COPY . .

# Build the application binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o metaparser .

# Use a minimal Alpine Linux image as a runtime base
FROM alpine:3.14

# Set the working directory to /app
WORKDIR /app

# Copy the application binary from the builder image
COPY --from=builder /app/metaparser .

ENV PORT = 9000
# Expose port 9000 for the application to listen on
EXPOSE 9000

# Run the application binary when the container starts
CMD ["./metaparser"]