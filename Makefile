-include .env

build:
	GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server ./cmd/api

run:
	docker compose -f docker-compose.yaml up -d --build

test:
	GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o cli-test ./cmd/cli

.PHONY: build run test