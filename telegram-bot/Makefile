# Makefile for Telegram Bot

.PHONY: help build run clean test deps lint

# Default target
help:
	@echo "Available commands:"
	@echo "  build     - Build the telegram bot"
	@echo "  run       - Run the telegram bot"
	@echo "  clean     - Clean build artifacts"
	@echo "  test      - Run tests"
	@echo "  deps      - Download dependencies"
	@echo "  lint      - Run linter"
	@echo "  docker    - Build Docker image"

# Build the application
build:
	@echo "Building telegram bot..."
	go build -o bin/telegram-bot main.go

# Run the application
run:
	@echo "Running telegram bot..."§
	go run main.go

# Run with custom config
run-config:
	@echo "Running telegram bot with custom config..."
	go run main.go -config $(CONFIG)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf data/
	go clean

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Build for different platforms
build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o bin/telegram-bot-linux main.go

build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o bin/telegram-bot-windows.exe main.go

build-mac:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build -o bin/telegram-bot-mac main.go

# Build all platforms
build-all: build-linux build-windows build-mac

# Create data directory
setup:
	@echo "Setting up directories..."
	mkdir -p data
	mkdir -p bin

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t telegram-bot .

docker-run:
	@echo "Running Docker container..."
	docker run --rm -v $(PWD)/config.yaml:/app/config.yaml telegram-bot

# Development commands
dev: deps setup
	@echo "Setting up development environment..."

# Production build
prod: clean deps build-linux
	@echo "Production build complete" 