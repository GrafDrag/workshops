# syntax=docker/dockerfile:1

FROM golang:1.17

# Set the Current Working Directory inside the container
WORKDIR /build

ENV GO111MODULE=on

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

COPY . /build

# Build the application
RUN go build -o ./ ./cmd/apiserver
RUN go build -o ./ ./cmd/grpcserver

COPY ./configs/server_docker.toml ./configs/server.toml

# Expose ports to the outside world
RUN chmod +x ./apiserver
RUN chmod +x ./grpcserver

CMD ["./apiserver","./grpcserver"]
