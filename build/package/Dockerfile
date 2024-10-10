# Stage 1: Build the Go application
FROM golang:1.22.0 AS build

WORKDIR /go/src/goapp

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go application (build the entire directory, not just main.go)
RUN go build -o /go/src/goapp/bin/goapp ./cmd/server

# Stage 2: Create a minimal runtime environment
FROM debian:12-slim

# Install necessary system dependencies (ca-certificates, net-tools, etc.)
RUN apt-get -y update && \
    apt-get -y upgrade && \
    apt-get -y install ca-certificates iputils-ping net-tools netcat-traditional procps tzdata && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Copy the Go binary from the build stage
COPY --from=build /go/src/goapp/bin/goapp /goapp/goapp

# Ensure the binary has executable permissions
RUN chmod +x /goapp/goapp

# Set the entry point to execute the Go binary
ENTRYPOINT ["/goapp/goapp"]
