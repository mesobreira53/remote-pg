# Stage 1: Build the Go application
FROM golang:1.23-alpine as builder

# Set the working directory in the container
WORKDIR /app

# Copy go mod and sum files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the source code
COPY . .

# Build the Go application
RUN go build -o pg-operator-proxy .

# Stage 2: Create the minimal runtime image
FROM alpine:latest

# Install necessary libraries (if needed for your Go app)
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the built Go binary from the builder image
COPY --from=builder /app/pg-operator-proxy .

# Expose the port that your Go app uses (e.g., 8080)
EXPOSE 8080

# Command to run the Go app
CMD ["./pg-operator-proxy"]