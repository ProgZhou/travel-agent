.PHONY: help build run test clean

help:
	@echo "Available commands:"
	@echo "  make build    - Build the server binary"
	@echo "  make run      - Run the server"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"

build:
	@echo "Building server..."
	go build -o bin/server cmd/server/main.go

run:
	@echo "Running server..."
	go run cmd/server/main.go

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning..."
	rm -rf bin/
