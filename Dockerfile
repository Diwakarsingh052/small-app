# micro-1/Dockerfile

# Stage 1: Build the Go application
FROM golang:1.24.1-alpine3.21 AS builder

WORKDIR /app


# Set the target OS and architecture
ENV GOOS=linux
ENV GOARCH=amd64

# Copy go.mod and go.sum first to leverage layer caching
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Run go mod tidy to ensure that dependencies are consistent
RUN go mod tidy

# Copy the source code into the container
COPY . .


# Build the Go application
RUN go build -o small-app .

# Stage 2: Create a lightweight runtime image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/small-app .

COPY private.pem .
COPY pubkey.pem .

EXPOSE 80

CMD ["./small-app"]



