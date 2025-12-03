.PHONY: build run fmt lint test

build:
	go build -o bin/dhund ./cmd/dhund

run:
	go run ./cmd/dhund

fmt:
	go fmt ./...

lint:
	golangci-lint run

test:
	go test ./...
