# Go Modular Monolith

A production-ready Go application template following modular monolith architecture principles with Clean Architecture patterns. Built with Fiber web framework, GORM ORM, and MySQL database.

## 🏗️ Architecture

### Modular Structure
Each business domain is organized as a self-contained module with clear separation of concerns:

```
internal/modules/{domain}/
├── domain/     # Business entities, interfaces, domain logic
├── app/        # Application services, commands, queries
├── http/       # HTTP handlers, DTOs, routing
├── infra/      # Infrastructure implementations (repositories)
└── module.go   # Module registration and dependency wiring
```

### Shared Components
- **Config**: Centralized configuration management with Viper
- **Database**: GORM connection management and utilities
- **HTTP**: Common HTTP utilities and health checks
- **Module**: Core module interface and registration system

## 🚀 Quick Start

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

## 📡 API Endpoints

### Health Check
- `GET /health-check/alive` - Application health status

### User Management
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `PUT /api/v1/users/:id/activate` - Activate user
- `PUT /api/v1/users/:id/deactivate` - Deactivate user
- `GET /api/v1/users` - List users (with pagination)

## 🔧 Configuration

Configure via environment variables or `.env` file:

```env
PORT=8080
DATABASE_URL=root:password@tcp(localhost:3306)/modular_monolith
```

## 🏛️ Key Design Patterns

### Module Interface
```go
type Module interface {
    RegisterRoutes(router fiber.Router)
    GetModels() []any
}
```

### Clean Architecture Layers
1. **Domain**: Pure business logic, no external dependencies
2. **Application**: Use cases and business workflows
3. **Infrastructure**: External concerns (database, HTTP)
4. **HTTP**: Web layer (handlers, DTOs, routing)

### Database Management
- Auto-migration on startup
- Repository pattern for data access
- GORM with MySQL driver

## 📦 Adding New Modules

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

## 🔍 Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific module tests
go test ./internal/modules/user/...
```

## 🐳 Docker

- **Application**: Multi-stage build with Alpine final image
- **Database**: MySQL 8.0 with health checks
- **Compose**: Orchestrates app and database with proper dependencies for development

## 📋 Development Guidelines

### Code Organization
- Follow Clean Architecture principles
- Keep modules independent and self-contained
- Use dependency injection at module boundaries
- Implement proper error handling and logging

### Database
- Use GORM struct tags for schema definition
- Implement repository interfaces in domain layer
- Handle migrations automatically via AutoMigrate

### HTTP Layer
- Use DTOs for request/response serialization
- Implement proper error handling middleware
- Follow RESTful conventions

## 🛠️ Tech Stack

- **Framework**: Fiber v2 (Express-like Go web framework)
- **ORM**: GORM (Go Object-Relational Mapping)
- **Database**: MySQL 8.0
- **Config**: Viper (configuration management)
- **Testing**: Go standard testing package
- **Containerization**: Docker & Docker Compose
