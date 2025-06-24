# Go Modular Monolith Makefile

ifeq ($(OS),Windows_NT)
    BINARY_EXT := .exe
    RM_CMD := del /Q
    MKDIR_CMD := mkdir
else
    BINARY_EXT :=
    RM_CMD := rm -f
    MKDIR_CMD := mkdir -p
endif

APP_NAME=modular-monolith
MAIN_PATH=cmd/api/main.go
BINARY_NAME := main$(BINARY_EXT)
BUILD_DIR=bin

.PHONY: help dev-local dev db-up build build-windows build-linux build-all clean test test-verbose test-coverage lint deps down docker-logs migrate-up migrate-down migrate-status migrate-force migrate-create

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
	go build -o $(BUILD_DIR)/$(BINARY_PATH) $(MAIN_PATH)
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

migrate-up: ## Executa todas as migrações pendentes
	go run cmd/migrate/main.go -action=up
	@echo "Migrações executadas!"

migrate-down: ## Desfaz a última migração
	go run cmd/migrate/main.go -action=down
	@echo "Migração desfeita!"

migrate-status: ## Mostra o status atual das migrações
	go run cmd/migrate/main.go -action=status

migrate-force: ## Força uma versão específica (uso: make migrate-force VERSION=1)
	@if [ -z "$(VERSION)" ]; then \
		echo "Erro: VERSION é obrigatório. Uso: make migrate-force VERSION=1"; \
		exit 1; \
	fi
	go run cmd/migrate/main.go -action=force -version=$(VERSION)
	@echo "Versão $(VERSION) forçada!"

migrate-create: ## Cria uma nova migração (uso: make migrate-create NAME=create_products_table)
	@if [ -z "$(NAME)" ]; then \
		echo "Erro: NAME é obrigatório. Uso: make migrate-create NAME=create_products_table"; \
		exit 1; \
	fi
	@TIMESTAMP=$$(date +%s); \
	PADDED_TIMESTAMP=$$(printf "%06d" $$TIMESTAMP); \
	UP_FILE="migrations/$${PADDED_TIMESTAMP}_$(NAME).up.sql"; \
	DOWN_FILE="migrations/$${PADDED_TIMESTAMP}_$(NAME).down.sql"; \
	echo "-- Migração: $(NAME)" > $$UP_FILE; \
	echo "-- TODO: Adicionar comandos SQL aqui" >> $$UP_FILE; \
	echo "" >> $$UP_FILE; \
	echo "-- Rollback: $(NAME)" > $$DOWN_FILE; \
	echo "-- TODO: Adicionar comandos de rollback aqui" >> $$DOWN_FILE; \
	echo "" >> $$DOWN_FILE; \
	echo "Migração criada:"; \
	echo "  $$UP_FILE"; \
	echo "  $$DOWN_FILE"