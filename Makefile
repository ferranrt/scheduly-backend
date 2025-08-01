.PHONY: help build run test clean deps migrate

# Default target
help:
	@echo "Available commands:"
	@echo "  deps     - Install dependencies"
	@echo "  build    - Build the application"
	@echo "  run      - Run the application"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  migrate  - Run database migrations"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Build the application
build:
	@echo "Building application..."
	go build -o bin/scheduly-backend main.go

# Run the application
run:
	@echo "Running application..."
	go run main.go

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean

# Run database migrations
migrate:
	@echo "Running database migrations..."
	go run cmd/dbmigrate/migrate.go

# Development mode with hot reload (requires air)
dev:
	@echo "Starting development server with hot reload..."
	air

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t scheduly-backend .

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 scheduly-backend

# Docker compose up
docker-up:
	@echo "Starting services with Docker Compose..."
	docker-compose up -d

# Docker compose down
docker-down:
	@echo "Stopping services with Docker Compose..."
	docker-compose down

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Generate API documentation
docs:
	@echo "Generating API documentation..."
	swag init -g main.go
