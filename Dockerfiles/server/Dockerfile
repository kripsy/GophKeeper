# syntax=docker/dockerfile:1

FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
COPY gen/ gen/
COPY cert/ cert/
COPY *.md ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server/main.go

FROM scratch

COPY --from=builder /bin/server /bin/server
COPY --from=builder /app/cert /cert
EXPOSE 8080 50051
CMD ["/bin/server"]