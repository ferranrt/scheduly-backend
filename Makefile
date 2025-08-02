.PHONY: help build run test clean deps migrate

# Default target
help:
	@echo "Available commands:"
	@echo "  deps     - Install dependencies"
	@echo "  build    - Build the application"
	@echo "  run      - Run the application"
	@echo "  test     - Run tests
	@echo "  migrate  - Run database migrations"
	@echo "  rebuild  - Rebuild database tables"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Build the application
build:
	@echo "Building application..."
	go build -o bin/scheduly-core-rest cmd/rest/main.go

start:
	@echo "Starting application..."
	./bin/scheduly-core-rest

# Run the application
run:
	@echo "Running application..."
	go run cmd/rest/main.go

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
	go run cmd/dbtools/main.go migrate


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

# Generate API documentation
docs:
	@echo "Generating API documentation..."
	swag init -g main.go
