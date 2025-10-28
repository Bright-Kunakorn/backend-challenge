FROM golang:1.21-alpine AS builder

WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/api

FROM alpine:3.19

RUN addgroup -S app && adduser -S app -G app
WORKDIR /home/app

COPY --from=builder /app/server /usr/local/bin/server

EXPOSE 8080 50051
USER app

ENTRYPOINT ["/usr/local/bin/server"]
