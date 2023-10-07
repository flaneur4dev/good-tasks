BIN := "./bin/calendar"
DOCKER_IMG := "calendar:develop"
COMPOSE_CONFIG := ./deployments/compose.yaml
PB_PATH := internal/server/grpc/pb

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

generate:
	protoc \
		--proto_path=api \
		--go_out=$(PB_PATH) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_PATH) --go-grpc_opt=paths=source_relative \
		event_service.proto

up:
	docker compose -f $(COMPOSE_CONFIG) up -d

down:
	docker compose -f $(COMPOSE_CONFIG) down -v

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img generate up down version test lint
