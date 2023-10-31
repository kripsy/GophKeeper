# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.21

WORKDIR /app

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download

RUN mkdir -p -m 0600 ./bin
COPY ./cmd/server ./cmd/server
COPY ./internal/server ./internal/server
COPY ./internal/utils ./internal/utils
COPY ./internal/logger ./internal/logger
COPY ./internal/models ./internal/models
COPY ./pkg ./pkg
COPY ./gen ./gen
COPY ./cert ./cert
COPY ./*.md ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/server ./cmd/server/main.go
CMD ["./bin/server"]

EXPOSE 8080 50051