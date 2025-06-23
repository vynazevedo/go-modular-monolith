# Go Modular Monolith

## Architecture

```
internal/modules/{domain}/
├── domain/     # Business entities, interfaces, domain logic
├── app/        # Application services, commands, queries
├── http/       # HTTP handlers, DTOs, routing
├── infra/      # Infrastructure implementations (repositories)
└── module.go   # Module registration and dependency wiring
```

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- MySQL 8.0+ (or use Docker)

### Development Setup

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd go-modular-monolith
   make deps
   ```

2. **Start with Docker** (recommended):
   ```bash
   make dev
   ```

3. **Start locally** (requires MySQL running):
   ```bash
   make db-up        # Start only database
   make dev-local    # Run app locally
   ```

### Available Commands

| Command | Description |
|---------|-------------|
| `make dev` | Run application with Docker Compose |
| `make dev-local` | Run application locally |
| `make db-up` | Start only the database container |
| `make build` | Build application binary |
| `make test` | Run all tests |
| `make test-coverage` | Run tests with coverage report |
| `make lint` | Run Go linting and formatting |
| `make clean` | Remove build artifacts |

## Configuration

Configure via environment variables or `.env` file:

```env
PORT=8080
DATABASE_URL=root:password@tcp(localhost:3306)/modular_monolith
```


## Adding New Modules

1. Create module directory structure:
   ```bash
   mkdir -p internal/modules/product/{domain,app,http,infra}
   ```

2. Implement the `Module` interface in `module.go`

3. Register module in `cmd/api/main.go`:
   ```go
   productModuleSetup := func(db *gorm.DB) module.Module {
       return product.NewModule(db)
   }
   modules := module.SetupAllModules(db, userModuleSetup, productModuleSetup)
   ```

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific module tests
go test ./internal/modules/user/...
```
