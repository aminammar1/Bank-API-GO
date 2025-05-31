# Tunisian Bank API Makefile

.PHONY: build run test clean deps dev help

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

# Show available commands
help:
	@echo "Tunisian Bank API - Available Commands:"
	@echo ""
	@echo "  build    Build the application"
	@echo "  run      Build and run the application"  
	@echo "  dev      Run in development mode"
	@echo "  test     Run tests"
	@echo "  deps     Install dependencies"
	@echo "  clean    Clean build artifacts"
	@echo "  help     Show this help"
	@echo ""