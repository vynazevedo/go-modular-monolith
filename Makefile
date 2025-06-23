# Go Modular Monolith Makefile

APP_NAME=go-modular-monolith
MAIN_PATH=cmd/api/main.go
BUILD_DIR=bin

.PHONY: help dev docker-dev db-up build clean test test-verbose test-coverage lint deps

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

dev-local: ## Run the application locally
	go run $(MAIN_PATH)

dev: ## Run application with Docker Compose (app + database)
	docker compose up --build

db-up: ## Start only the database container
	docker compose up -d db

build: ## Build the application binary
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)
	go clean

test: ## Run all tests
	go test ./...

test-verbose: ## Run tests with verbose output
	go test -v ./...

test-coverage: ## Run tests with coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run Go linting (format check, vet and modernize)
	gofmt -l .
	go vet ./...
	go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix -test ./...

deps: ## Download and tidy dependencies
	go mod download
	go mod tidy

down: ## Stop and remove Docker containers
	docker compose down

docker-logs: ## Show Docker container logs
	docker compose logs -f
