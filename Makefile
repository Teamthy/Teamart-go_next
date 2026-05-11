.PHONY: help build run dev docker-up docker-down docker-logs test lint fmt clean install-deps

help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make dev            - Run in development mode with hot reload"
	@echo "  make docker-up      - Start Docker containers (PostgreSQL, Redis)"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make docker-logs    - View Docker logs"
	@echo "  make test           - Run tests"
	@echo "  make lint           - Run linter"
	@echo "  make fmt            - Format code"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make install-deps   - Install Go dependencies"

install-deps:
	@echo "Installing Go dependencies..."
	cd backend && go mod download && go mod tidy

build:
	@echo "Building application..."
	cd backend && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api ./cmd/api

run: build
	@echo "Running application..."
	cd backend && ./bin/api

dev:
	@echo "Running in development mode..."
	cd backend && go run ./cmd/api/main.go

docker-up:
	@echo "Starting Docker containers..."
	cd backend && docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "Services are ready!"
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"

docker-down:
	@echo "Stopping Docker containers..."
	cd backend && docker-compose down

docker-logs:
	cd backend && docker-compose logs -f

docker-clean:
	@echo "Removing Docker containers and volumes..."
	cd backend && docker-compose down -v

test:
	@echo "Running tests..."
	cd backend && go test -v -race -coverprofile=coverage.out ./...

lint:
	@echo "Running linter..."
	cd backend && go fmt ./...
	cd backend && go vet ./...

fmt:
	@echo "Formatting code..."
	cd backend && go fmt ./...

clean:
	@echo "Cleaning build artifacts..."
	cd backend && rm -f bin/api
	cd backend && go clean

migrate-up:
	@echo "Running migrations (placeholder)"
	@echo "Implement migration tooling with sql-migrate or similar"

migrate-down:
	@echo "Reverting migrations (placeholder)"
	@echo "Implement migration tooling with sql-migrate or similar"

.DEFAULT_GOAL := help
