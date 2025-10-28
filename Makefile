.PHONY: start stop build test proto lint smoke

DOCKER_COMPOSE := docker-compose

start:
	$(DOCKER_COMPOSE) up --build -d

stop:
	$(DOCKER_COMPOSE) down

build:
	go build ./cmd/api

test:
	CGO_ENABLED=0 go test ./...

proto:
	@echo "Generating gRPC stubs (requires protoc and protoc-gen-go)..."
	protoc --go_out=. --go-grpc_out=. proto/user.proto

lint:
	go fmt ./...

smoke:
	./scripts/api-smoke.sh
