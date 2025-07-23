# Makefile for User Management API

# Variables
GO_VERSION := 1.24
BINARY_NAME := user-management
DOCKER_COMPOSE_FILE := docker-compose.test.yml

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
.PHONY: build
build: ## Build the application
	go build -o bin/$(BINARY_NAME) .

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/
	go clean

# Development targets
.PHONY: run
run: ## Run the application
	go run .

.PHONY: dev
dev: ## Run the application with hot reload (requires air)
	air

# Dependency targets
.PHONY: deps
deps: ## Download dependencies
	go mod download
	go mod tidy

.PHONY: deps-update
deps-update: ## Update dependencies
	go get -u ./...
	go mod tidy

# Wire generation
.PHONY: wire
wire: ## Generate wire dependencies
	cd cmd && wire

# Test targets
.PHONY: test-all
test-all: test test-integration ## Run all tests (unit and integration)

.PHONY: test
test: ## Run unit tests with coverage
	@echo "Running unit tests..."
	go test -v -race -coverprofile=coverage.out ./internal/...
	go tool cover -func=coverage.out

.PHONY: test-integration
test-integration: ## Run integration tests with coverage
	@echo "Starting integration tests..."
	@echo "Make sure MongoDB is running on localhost:27017"
	go test -v -race -coverprofile=coverage-integration.out ./tests/...
	go tool cover -func=coverage-integration.out

.PHONY: test-integration-docker
test-integration-docker: ## Run integration tests with Docker MongoDB
	@echo "Starting MongoDB for integration tests..."
	docker run --name mongo-test -p 27017:27017 -d mongo:7.0 || docker start mongo-test
	@echo "Waiting for MongoDB to be ready..."
	@sleep 5
	@echo "Running integration tests..."
	go test -v -race -coverprofile=coverage-integration.out -tags=integration ./tests/... || (docker stop mongo-test && exit 1)
	go tool cover -func=coverage-integration.out
	@echo "Stopping test MongoDB..."
	docker stop mongo-test
	docker rm mongo-test

.PHONY: test-coverage
test-coverage: test test-integration ## Generate test coverage report
	@echo "Merging coverage reports..."
	@echo "mode: atomic" > coverage-total.out
	@tail -n +2 coverage.out >> coverage-total.out
	@tail -n +2 coverage-integration.out >> coverage-total.out
	go tool cover -html=coverage-total.out -o coverage.html
	go tool cover -func=coverage-total.out
	@echo "Coverage report generated: coverage.html"

.PHONY: test-clean
test-clean: ## Clean test artifacts
	rm -f coverage*.out coverage.html

.PHONY: sonar
sonar: test-coverage ## Run SonarQube analysis with coverage
	@echo "Running SonarQube analysis..."
	sonar-scanner \
		-Dsonar.login=${SONAR_TOKEN} \
		-Dsonar.host.url=${SONAR_HOST_URL} \
		-Dsonar.projectVersion=$(shell git rev-parse --short HEAD) \
		-Dsonar.go.coverage.reportPaths=coverage-total.out \
		-Dsonar.go.tests.reportPaths=test-report.json \
		-Dsonar.qualitygate.wait=true

.PHONY: sonar-local
sonar-local: test-coverage ## Run SonarQube analysis locally
	@echo "Running local SonarQube analysis..."
	sonar-scanner \
		-Dsonar.host.url=http://localhost:9000 \
		-Dsonar.projectVersion=$(shell git rev-parse --short HEAD) \
		-Dsonar.go.coverage.reportPaths=coverage-total.out \
		-Dsonar.go.tests.reportPaths=test-report.json \
		-Dsonar.qualitygate.wait=true

# Database targets
.PHONY: mongo-start
mongo-start: ## Start MongoDB using Docker
	docker run --name mongo-dev -p 27017:27017 -d mongo:7.0

.PHONY: mongo-stop
mongo-stop: ## Stop MongoDB Docker container
	docker stop mongo-dev || true
	docker rm mongo-dev || true

.PHONY: mongo-logs
mongo-logs: ## Show MongoDB logs
	docker logs -f mongo-dev

# Linting and formatting
.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: fmt
fmt: ## Format code
	go fmt ./...
	goimports -w .

.PHONY: vet
vet: ## Run go vet
	go vet ./...

# Quality targets
.PHONY: check
check: fmt vet lint test ## Run all quality checks

# Docker targets
.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(BINARY_NAME):latest .

.PHONY: docker-run
docker-run: ## Run application in Docker
	docker run --rm -p 8080:8080 $(BINARY_NAME):latest

# Docker Compose targets
.PHONY: up
up: ## Start all services with Docker Compose
	docker-compose up -d

.PHONY: down
down: ## Stop all services with Docker Compose
	docker-compose down

.PHONY: logs
logs: ## Show logs from all services
	docker-compose logs -f

.PHONY: logs-app
logs-app: ## Show logs from application service
	docker-compose logs -f app

.PHONY: logs-mongo
logs-mongo: ## Show logs from MongoDB service
	docker-compose logs -f mongodb

.PHONY: restart
restart: ## Restart all services
	docker-compose restart

.PHONY: rebuild
rebuild: ## Rebuild and restart all services
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

# Environment setup
.PHONY: env-setup
env-setup: ## Setup development environment
	@echo "Setting up development environment..."
	@if ! command -v go >/dev/null 2>&1; then \
		echo "Go is not installed. Please install Go $(GO_VERSION) or later."; \
		exit 1; \
	fi
	@if ! command -v air >/dev/null 2>&1; then \
		echo "Installing air for hot reload..."; \
		go install github.com/air-verse/air@latest; \
	fi
	@if ! command -v wire >/dev/null 2>&1; then \
		echo "Installing wire for dependency injection..."; \
		go install github.com/google/wire/cmd/wire@latest; \
	fi
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.55.2; \
	fi
	@echo "Development environment setup complete!"

# Example .env file
.PHONY: env-example
env-example: ## Create example .env file
	@echo "Creating .env.example file..."
	@cat > .env.example << 'EOF'
# MongoDB Configuration
MONGO_URI=mongodb://localhost:27017
MONGO_DB=user_management

# Server Configuration  
PORT=:8080

# Test Configuration (optional)
TEST_MONGO_URI=mongodb://localhost:27017
TEST_MONGO_DB=user_management_test
TEST_PORT=:3001
EOF
	@echo ".env.example created! Copy to .env and adjust values as needed."

# All-in-one development setup
.PHONY: setup
setup: env-setup deps wire env-example ## Complete development setup
	@echo "✅ Development setup complete!"
	@echo "📝 Next steps:"
	@echo "   1. Copy .env.example to .env and configure your settings"
	@echo "   2. Start MongoDB: make mongo-start"
	@echo "   3. Run the application: make run"
	@echo "   4. Run tests: make test-integration"
