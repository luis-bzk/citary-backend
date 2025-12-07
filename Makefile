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
	@echo Available commands:
	@echo   build                Build the application locally
	@echo   run                  Run the application locally
	@echo   fmt                  Format code
	@echo   lint                 Run linter
	@echo   clean                Clean build artifacts
	@echo   docker-build         Build Docker image
	@echo   docker-build-no-cache Build Docker image without cache
	@echo   docker-push          Build and push Docker image to registry
	@echo   docker-run           Run Docker container locally
	@echo   docker-test          Build and test Docker image locally
	@echo   deps-download        Download dependencies
	@echo   deps-tidy            Tidy dependencies
	@echo   deps-verify          Verify dependencies
	@echo   deps-update          Update all dependencies
	@echo   dev-setup            Setup development environment
	@echo   dev-db               Start PostgreSQL with Docker
	@echo   dev-db-stop          Stop PostgreSQL container
	@echo   version              Show version
	@echo   info                 Show build information

# Local Development
build: ## Build the application locally
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME) ./cmd/api

run: ## Run the application locally
	@echo "Running $(APP_NAME)..."
	go run ./cmd/api


fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

lint: ## Run linter
	@echo "Running linter..."
	go vet ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	-del /Q $(APP_NAME) $(APP_NAME).exe coverage.out coverage.html 2>nul
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
	docker run --rm -p 3001:3001 \
		-e DATABASE_URL="$(DATABASE_URL)" \
		-e PORT=3001 \
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
