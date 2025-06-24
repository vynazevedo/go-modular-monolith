![Go](https://img.shields.io/badge/Go-1.24.4-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?logo=go&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-ORM-blue)
![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?logo=docker&logoColor=white)
![Logrus](https://img.shields.io/badge/Logrus-Logger-green)

# Go Template DDD Module

> Template pronto para uso de um monólito modular seguindo Domain-Driven Design

Essa é uma implementação prática de como estruturar um projeto Go seguindo padrões de DDD dentro de uma arquitetura de monólito modular. Ideal para quem quer ter um ponto de partida sólido sem a complexidade inicial de microsserviços.

## Por que monólito modular?

- **Deploy simplificado** - Um único binário
- **Módulos independentes** - Cada domínio vive isolado
- **Evolução gradual** - Migre para microsserviços quando fizer sentido
- **Manutenção facilitada** - Menos complexidade operacional

## Estrutura dos módulos

```
internal/modules/{dominio}/
├── domain/     # Entidades, Value Objects, regras de negócio
├── app/        # Casos de uso, Commands, Queries
├── http/       # Controllers, DTOs, validações
├── infra/      # Implementações de repositórios
└── module.go   # Configuração e injeção de dependências
```

Cada módulo é completamente isolado e só se comunica através de interfaces bem definidas.

## Stack & Ferramentas

- **Web Framework**: Gin (migrado do Fiber)
- **ORM**: GORM com MySQL
- **Logging**: Logrus estruturado
- **Config**: Viper com variáveis de ambiente
- **Middlewares**: Exemplos de autenticação e validação
- **Testes**: Go testing nativo
- **Build**: Makefile com cross-compilation

## Setup rápido

```bash
# Clone o template
git clone https://github.com/vynazevedo/go-modular-monolith.git go-template-ddd-module

# Entre no diretório
cd go-template-ddd-module

# Instale dependências  
make deps

# Rode com Docker (banco incluso)
make dev

# Ou rode local (precisa do MySQL)
make db-up && make dev-local
```

## Variáveis de ambiente

Copie `.env.example` para `.env` e ajuste conforme necessário:

```bash
# Server
PORT=8080
GIN_MODE=debug

# Database (componentes separados)
DB_HOST=localhost
DB_PORT=3306  
DB_USER=root
DB_PASSWORD=password
DB_NAME=modular_monolith

# Logging
LOG_LEVEL=info
LOG_FORMAT=text

# CORS (para frontend React)
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Origin,Content-Type,Accept,Authorization,X-API-Key
CORS_MAX_AGE=86400
```

## Build para produção

```bash
# Build local
make build

# Cross-compilation
make build-windows    # Windows .exe
make build-linux      # Linux binary  
make build-all        # Todas as plataformas
```

## Comandos úteis

### Desenvolvimento
| Comando | O que faz |
|---------|-----------|
| `make help` | Lista todos os comandos disponíveis |
| `make dev` | Sobe app + banco com Docker |
| `make dev-local` | Roda só a aplicação local |
| `make test` | Executa todos os testes |
| `make test-coverage` | Gera relatório de cobertura |
| `make lint` | Formatação e análise de código |

### Migrações de Banco
| Comando | O que faz |
|---------|-----------|
| `make migrate-up` | Executa migrações pendentes |
| `make migrate-down` | Desfaz última migração |
| `make migrate-status` | Status das migrações |
| `make migrate-create NAME=exemplo` | Cria nova migração |

Veja `docs/migrations.md` para guia completo de migrações.

## Como adicionar um novo módulo

1. **Crie a estrutura**:
```bash
mkdir -p internal/modules/produto/{domain,app,http,infra}
```

2. **Implemente as camadas** seguindo o padrão do módulo `user`

3. **Registre no main.go**:
```go
produtoModule := func(db *gorm.DB) module.Module {
    return produto.NewModule(db)
}
modules := module.SetupAllModules(db, userModule, produtoModule)
```

## Testando

```bash
# Todos os testes
make test

# Com cobertura detalhada  
make test-coverage

# Testes de um módulo específico
go test ./internal/modules/user/...
```

## Middleware de exemplo

O projeto inclui um middleware de validação de API Key como exemplo:

```bash
# Requisição protegida (precisa do header)
curl -X POST http://localhost:8080/api/v1/users \
  -H "X-API-Key: api-key-exemplo" \
  -H "Content-Type: application/json" \
  -d '{"name": "Vinicius", "email": "vinicius@teste.com"}'

# Requisição pública (header opcional)  
curl http://localhost:8080/api/v1/users
```
---

**Dica**: Este template foi pensado para crescer com seu projeto. Comece simples e evolua conforme a necessidade.

## Contribuindo

Encontrou algo que pode melhorar? PRs são bem-vindos! 
