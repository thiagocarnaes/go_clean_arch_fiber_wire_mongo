# Makefile for User Management API

# Variables
GO_VERSION := 1.24
BINARY_NAME := user-management

# Colors for terminal output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
RESET := \033[0m

.PHONY: help
help: ## Show this help message
	@echo "$(BLUE)User Management API - Available Commands$(RESET)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make $(GREEN)<command>$(RESET)\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(GREEN)%-15s$(RESET) %s\n", $$1, $$2 } /^##@/ { printf "\n$(YELLOW)%s$(RESET)\n", substr($$0, 5) }' $(MAKEFILE_LIST)

##@ Development Commands

.PHONY: setup
setup: ## Setup complete development environment
	@echo "$(BLUE)Setting up development environment...$(RESET)"
	go mod download
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/air-verse/air@latest
	$(MAKE) wire
	@echo "$(GREEN)Setup completed!$(RESET)"

.PHONY: dev
dev: ## Start development server with hot reload using Air
	@echo "$(BLUE)Starting development server with hot reload...$(RESET)"
	air

.PHONY: run
run: ## Run the application in production mode
	@echo "$(BLUE)Starting application...$(RESET)"
	go run main.go initApiServer

.PHONY: wire
wire: ## Regenerate Wire dependency injection code
	@echo "$(BLUE)Regenerating Wire code...$(RESET)"
	cd cmd && wire && cd ..
	@echo "$(GREEN)Wire code generated successfully!$(RESET)"

.PHONY: deps-update
deps-update: ## Update Go dependencies
	@echo "$(BLUE)Updating dependencies...$(RESET)"
	go get -u ./...
	go mod tidy

##@ Testing Commands

.PHONY: test
test: ## Run all tests
	@echo "$(BLUE)Running all tests...$(RESET)"
	go test -v -race ./...

.PHONY: test-integration
test-integration: ## Run integration tests with coverage
	@echo "$(BLUE)Running integration tests with coverage...$(RESET)"
	go test -v -race -coverprofile=coverage-integration.out --coverpkg=./... ./tests/...
	go tool cover -func=coverage-integration.out

.PHONY: test-coverage
test-coverage: ## Generate test coverage report
	@echo "$(BLUE)Generating coverage report...$(RESET)"
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report generated: coverage.html$(RESET)"

.PHONY: coverage-html
coverage-html: ## Open coverage report in browser
	go tool cover -html=coverage-integration.out

##@ Code Quality Commands

.PHONY: check
check: fmt lint vet ## Run all code quality checks
	@echo "$(GREEN)All quality checks passed!$(RESET)"

.PHONY: fmt
fmt: ## Format Go code
	@echo "$(BLUE)Formatting code...$(RESET)"
	go fmt ./...

.PHONY: lint
lint: ## Run linter
	@echo "$(BLUE)Running linter...$(RESET)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint not installed. Installing...$(RESET)"; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run; \
	fi

.PHONY: vet
vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(RESET)"
	go vet ./...

##@ Build Commands

.PHONY: build
build: ## Build the application binary
	@echo "$(BLUE)Building application...$(RESET)"
	go build -o bin/$(BINARY_NAME) main.go
	@echo "$(GREEN)Binary built: bin/$(BINARY_NAME)$(RESET)"

.PHONY: build-linux
build-linux: ## Build for Linux
	@echo "$(BLUE)Building for Linux...$(RESET)"
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux main.go

.PHONY: build-windows
build-windows: ## Build for Windows
	@echo "$(BLUE)Building for Windows...$(RESET)"
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows.exe main.go

.PHONY: build-mac
build-mac: ## Build for macOS
	@echo "$(BLUE)Building for macOS...$(RESET)"
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-mac main.go

##@ Docker Commands

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(BLUE)Building Docker image...$(RESET)"
	docker build -t $(BINARY_NAME):latest .
	@echo "$(GREEN)Docker image built successfully!$(RESET)"

.PHONY: docker-up
docker-up: ## Start all services with Docker Compose
	@echo "$(BLUE)Starting services with Docker Compose...$(RESET)"
	@if [ ! -f .env ]; then \
		echo "$(YELLOW)Warning: .env file not found. Creating example .env file...$(RESET)"; \
		echo "DD_API_KEY=your_datadog_api_key_here" > .env; \
		echo "DD_SITE=datadoghq.com" >> .env; \
		echo "DD_SOURCE=go" >> .env; \
		echo "DD_SERVICE=user-management" >> .env; \
		echo "DD_TAGS=env:docker,app:fiber" >> .env; \
		echo "$(RED)Please update .env file with your Datadog API key before running again$(RESET)"; \
		exit 1; \
	fi
	docker-compose up -d
	@echo "$(GREEN)All services started!$(RESET)"
	@echo "$(BLUE)Services available at:$(RESET)"
	@echo "  - API: http://localhost:8080"
	@echo "  - MongoDB Express: http://localhost:8081 (admin/admin)"

.PHONY: docker-down
docker-down: ## Stop all Docker Compose services
	@echo "$(BLUE)Stopping Docker Compose services...$(RESET)"
	docker-compose down
	@echo "$(GREEN)All services stopped!$(RESET)"

.PHONY: docker-logs
docker-logs: ## Show logs from all services
	@echo "$(BLUE)Showing logs from all services...$(RESET)"
	docker-compose logs -f

.PHONY: docker-logs-app
docker-logs-app: ## Show logs from app service only
	@echo "$(BLUE)Showing logs from app service...$(RESET)"
	docker-compose logs -f app

.PHONY: docker-restart
docker-restart: ## Restart all Docker Compose services
	@echo "$(BLUE)Restarting Docker Compose services...$(RESET)"
	docker-compose restart
	@echo "$(GREEN)All services restarted!$(RESET)"

.PHONY: docker-clean
docker-clean: ## Clean Docker containers, images and volumes
	@echo "$(BLUE)Cleaning Docker resources...$(RESET)"
	docker-compose down -v --remove-orphans
	docker system prune -f
	@echo "$(GREEN)Docker cleanup completed!$(RESET)"

##@ Cleanup Commands


.PHONY: clean
clean: ## Clean build artifacts and Go cache  
	@echo "$(BLUE)Cleaning all artifacts...$(RESET)"
	rm -rf bin/
	rm -f coverage*.out coverage*.html
	rm -f $(BINARY_NAME)
	go clean -cache
	go clean -modcache
	@echo "$(GREEN)Deep cleanup completed!$(RESET)"

##@ Release Commands

.PHONY: release
release: clean test build ## Prepare a release (clean, test, build)
	@echo "$(GREEN)Release prepared successfully!$(RESET)"

.PHONY: pre-commit
pre-commit: check test ## Run pre-commit checks
	@echo "$(GREEN)Pre-commit checks passed!$(RESET)"

##@ Utility Commands

.PHONY: deps
deps: ## Show project dependencies
	@echo "$(BLUE)Project dependencies:$(RESET)"
	go list -m all

.PHONY: deps-graph
deps-graph: ## Generate dependency graph
	@echo "$(BLUE)Generating dependency graph...$(RESET)"
	go mod graph

.PHONY: version
version: ## Show Go and tool versions
	@echo "$(BLUE)Version Information:$(RESET)"
	@echo "Go version: $$(go version)"
	@echo "Wire version: $$(wire --version 2>/dev/null || echo 'not installed')"
	@echo "Air version: $$(air -v 2>/dev/null || echo 'not installed')"

.PHONY: env
env: ## Show environment variables
	@echo "$(BLUE)Environment Variables:$(RESET)"
	@cat .env 2>/dev/null || echo "No .env file found"

# Default target
.DEFAULT_GOAL := help
