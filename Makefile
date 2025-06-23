# Go Modular Monolith Makefile

APP_NAME=modular-monolith
MAIN_PATH=cmd/api/main.go
BUILD_DIR=bin

.PHONY: help dev-local dev db-up build build-windows build-linux build-all clean test test-verbose test-coverage lint deps down docker-logs

help: ## Mostra esta mensagem de ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

dev-local: ## Executa a aplicação localmente
	go run $(MAIN_PATH)

dev: ## Executa aplicação com Docker Compose (app + banco de dados)
	docker compose up --build

db-up: ## Inicia apenas o container do banco de dados
	docker compose up -d db

build: ## Compila o binário da aplicação para o sistema atual
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "Aplicação compilada com sucesso: $(BUILD_DIR)/$(APP_NAME)"

build-windows: ## Compila o binário para Windows (amd64)
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows.exe $(MAIN_PATH)
	@echo "Aplicação Windows compilada: $(BUILD_DIR)/$(APP_NAME)-windows.exe"

build-linux: ## Compila o binário para Linux (amd64)
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux $(MAIN_PATH)
	@echo "Aplicação Linux compilada: $(BUILD_DIR)/$(APP_NAME)-linux"

build-all: build build-windows build-linux ## Compila para todos os sistemas operacionais
	@echo "Todas as versões compiladas com sucesso!"

clean: ## Remove artefatos de build
	rm -rf $(BUILD_DIR)
	go clean
	@echo "Arquivos de build removidos!"

test: ## Executa todos os testes
	go test ./...

test-verbose: ## Executa testes com saída detalhada
	go test -v ./...

test-coverage: ## Executa testes com relatório de cobertura
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Relatório de cobertura gerado: coverage.html"

lint: ## Executa linting do Go (formatação, vet e modernização)
	gofmt -l .
	go vet ./...
	go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix -test ./...

deps: ## Baixa e organiza dependências
	go mod download
	go mod tidy
	@echo "Dependências atualizadas!"

down: ## Para e remove containers Docker
	docker compose down

docker-logs: ## Mostra logs dos containers Docker
	docker compose logs -f