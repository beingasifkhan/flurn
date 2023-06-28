# Start from the official Go base image with Go version 1.17
FROM golang:1.20.2-alpine

# Set the working directory inside the container
WORKDIR /build

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main main.go

# Set the entry point for the container
ENTRYPOINT ["./main"]
