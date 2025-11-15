.PHONY: build run up test lint tidy docker-build

build:
	go build -o bin/server ./cmd/server

run: build
	./bin/server

up:
	docker-compose up --build

test:
	go test ./... -v

lint:
	golangci-lint run

tidy:
	go mod tidy

docker-build:
	docker build -t pr-reviewer:latest .
