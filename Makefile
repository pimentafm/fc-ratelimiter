-include .env

build:
	GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server ./cmd/api

init:
	go mod tidy

run:
	docker compose -f docker-compose.yaml up -d --build

build-cli:
	GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o cli-test ./cmd/cli

.PHONY: build init run build-cli