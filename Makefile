# Makefile for Citary Backend
# Simple commands for local development and Docker image building

.PHONY: help build run test clean docker-build docker-push

# Variables
APP_NAME := citary-backend
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
DOCKER_REGISTRY := your-registry.com
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(APP_NAME)
DOCKER_TAG := $(VERSION)

help: ## Show this help
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Local Development
build: ## Build the application locally
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME) ./cmd/api

run: ## Run the application locally
	@echo "Running $(APP_NAME)..."
	go run ./cmd/api

test: ## Run tests
	@echo "Running tests..."
	go test -v -race -cover ./...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

lint: ## Run linter
	@echo "Running linter..."
	go vet ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f $(APP_NAME) $(APP_NAME).exe coverage.out coverage.html
	go clean -cache -testcache

# Docker Commands
docker-build: ## Build Docker image
	@echo "Building Docker image: $(DOCKER_IMAGE):$(DOCKER_TAG)"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -t $(DOCKER_IMAGE):latest .

docker-build-no-cache: ## Build Docker image without cache
	@echo "Building Docker image (no cache): $(DOCKER_IMAGE):$(DOCKER_TAG)"
	docker build --no-cache -t $(DOCKER_IMAGE):$(DOCKER_TAG) -t $(DOCKER_IMAGE):latest .

docker-push: docker-build ## Build and push Docker image to registry
	@echo "Pushing Docker image: $(DOCKER_IMAGE):$(DOCKER_TAG)"
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_IMAGE):latest

docker-run: ## Run Docker container locally
	@echo "Running Docker container..."
	docker run --rm -p 3005:3005 \
		-e DATABASE_URL="$(DATABASE_URL)" \
		-e PORT=3005 \
		$(DOCKER_IMAGE):latest

docker-test: docker-build ## Build and test Docker image locally
	@echo "Testing Docker image..."
	docker run --rm $(DOCKER_IMAGE):$(DOCKER_TAG) --version || true

# Dependencies
deps-download: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download

deps-tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	go mod tidy

deps-verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	go mod verify

deps-update: ## Update all dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

# Development helpers
dev-setup: ## Setup development environment
	@echo "Setting up development environment..."
	cp .env.example .env
	@echo "Please edit .env with your configuration"
	go mod download

dev-db: ## Start PostgreSQL with Docker (for local dev only)
	@echo "Starting PostgreSQL container..."
	docker run --name citary-postgres-dev -d \
		-e POSTGRES_USER=dev_user \
		-e POSTGRES_PASSWORD=dev_password \
		-e POSTGRES_DB=citary_dev \
		-p 5432:5432 \
		postgres:16-alpine

dev-db-stop: ## Stop PostgreSQL container
	@echo "Stopping PostgreSQL container..."
	docker stop citary-postgres-dev || true
	docker rm citary-postgres-dev || true

# Info
version: ## Show version
	@echo "Version: $(VERSION)"

info: ## Show build information
	@echo "App Name:       $(APP_NAME)"
	@echo "Version:        $(VERSION)"
	@echo "Docker Image:   $(DOCKER_IMAGE):$(DOCKER_TAG)"
	@echo "Go Version:     $(shell go version)"
