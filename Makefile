# Tunisian Bank API Makefile

.PHONY: build run test clean deps dev help docker-build docker-run docker-stop docker-clean docker-logs docker-dev

.DEFAULT_GOAL := help

# Build the application
build:
	@echo "Building bank-api..."
	@go build -o bin/bank-api cmd/server/main.go

# Build and run
run: build
	@echo "Starting bank-api..."
	@./bin/bank-api

# Run in development mode
dev:
	@echo "Running in development mode..."
	@go run cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./tests/...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t bank-api-tunisia .

docker-run:
	@echo "Starting application with Docker Compose..."
	@docker-compose up -d

docker-dev:
	@echo "Starting application in development mode with Docker Compose..."
	@docker-compose up

docker-stop:
	@echo "Stopping Docker containers..."
	@docker-compose down

docker-clean:
	@echo "Cleaning Docker containers and volumes..."
	@docker-compose down -v --remove-orphans
	@docker image rm bank-api-tunisia 2>/dev/null || true

docker-logs:
	@echo "Showing application logs..."
	@docker-compose logs -f bank-api

# Show available commands
help:
	@echo "Tunisian Bank API - Available Commands:"
	@echo ""
	@echo "Local Development:"
	@echo "  build        Build the application"
	@echo "  run          Build and run the application"  
	@echo "  dev          Run in development mode"
	@echo "  test         Run tests"
	@echo "  deps         Install dependencies"
	@echo "  clean        Clean build artifacts"
	@echo ""
	@echo "Docker Commands:"
	@echo "  docker-build Build Docker image"
	@echo "  docker-run   Start with Docker Compose (background)"
	@echo "  docker-dev   Start with Docker Compose (foreground)"
	@echo "  docker-stop  Stop Docker containers"
	@echo "  docker-clean Clean Docker containers and volumes"
	@echo "  docker-logs  Show application logs"
	@echo ""
	@echo "  help         Show this help"
	@echo ""