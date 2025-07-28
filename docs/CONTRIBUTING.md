# Contributing to SARC-NG

This guide covers development setup, workflows, and contribution guidelines for SARC-NG.

## Development Setup

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- Make

### Initial Setup

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd sarc-ng
   make setup
   ```

2. **Choose development approach**:

   **Native Go (fastest for development):**
   ```bash
   make run        # Run application
   make debug      # Run with hot reloading
   ```

   **Docker environment (with database):**
   ```bash
   cd docker
   docker compose up -d    # Start environment
   ```

## Development Workflow

### Running the Application

```bash
# Native Go
make run                    # Direct execution
make debug                  # With hot reloading

# Docker
cd docker && docker compose up -d    # Start with database
```

### Testing

```bash
# Unit tests
make test

# Integration tests (requires Docker)
cd docker && docker compose up -d && docker compose run --rm app go test -tags=integration ./test/integration/... && docker compose down -v --remove-orphans

# Coverage
make coverage
```

### Code Quality

```bash
make format     # Format code
make lint       # Run linters
make test       # Run tests
```

### Building

```bash
make build      # Development build
make release    # Production build
```

### API Documentation

```bash
make swagger    # Generate documentation
# View at: http://localhost:8080/swagger/index.html
```

## Project Architecture

### Directory Structure

```
sarc-ng/
├── cmd/                    # Application entry points
│   ├── cli/               # CLI application
│   ├── lambda/            # AWS Lambda handler
│   └── server/            # HTTP server
├── internal/              # Private application code
│   ├── adapter/          # External service adapters
│   ├── config/           # Configuration management
│   ├── domain/           # Business entities and logic
│   ├── service/          # Application services
│   └── transport/        # HTTP handlers and routing
├── pkg/                   # Public libraries
├── configs/               # Configuration files
├── docker/                # Docker environment
├── infrastructure/        # Terraform infrastructure
└── test/                  # Test files
```

### Clean Architecture Layers

1. **Domain Layer** (`internal/domain/`): Core business entities and logic
2. **Service Layer** (`internal/service/`): Application services and use cases
3. **Transport Layer** (`internal/transport/`): HTTP handlers and routing
4. **Adapter Layer** (`internal/adapter/`): External service integration

### Key Technologies

- **Language**: Go 1.24+
- **Web Framework**: Gin
- **Database**: MySQL with GORM
- **Dependency Injection**: Google Wire
- **Configuration**: Viper with YAML
- **Documentation**: Swagger/OpenAPI
- **Testing**: Testify
- **Hot Reloading**: Air

## Development Tools

The following tools are automatically installed via `make setup`:

- **Air**: Hot reloading during development
- **Wire**: Dependency injection code generation
- **Swag**: API documentation generation
- **golangci-lint**: Code linting and static analysis
- **Delve**: Go debugger

### Wire Dependency Injection

After modifying dependency injection in `cmd/*/wire.go`:

```bash
make wire       # Regenerate injection code
make build      # Build includes wire generation
```

## Configuration

### Environment Variables

```bash
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=sarcng

# Application
PORT=8080
ENVIRONMENT=development
```

### Configuration Files

- `configs/default.yaml`: Base configuration
- `configs/development.yaml`: Development overrides

The application automatically loads the appropriate configuration based on the `ENVIRONMENT` variable.

## Docker Development

### Development Environment

```bash
cd docker
docker compose up -d                     # Start environment
docker compose down -v --remove-orphans  # Stop and cleanup
```

### Database Access

Database admin interface: http://localhost:8081 (admin/example)

### Container Operations

```bash
cd docker
docker compose down -v --remove-orphans  # Stop and cleanup all
```

## Testing Guidelines

### Unit Tests

- Place tests next to the code they test
- Use table-driven tests where appropriate
- Mock external dependencies
- Test business logic thoroughly

### Integration Tests

- Use build tags: `//go:build integration`
- Test against real database in Docker
- Test complete request/response flows
- Place in `test/integration/`

### Test Execution

```bash
# Unit tests only
make test

# Integration tests (requires Docker environment)
cd docker && docker compose up -d && docker compose run --rm app go test -tags=integration ./test/integration/... && docker compose down -v --remove-orphans
```

## API Development

### Adding New Endpoints

1. **Define entity** in `internal/domain/{entity}/`
2. **Create service** in `internal/service/{entity}/`
3. **Add handlers** in `internal/transport/rest/{entity}/`
4. **Register routes** in router
5. **Update documentation** with Swagger comments

### Swagger Documentation

Add Swagger comments to handlers:

```go
// @Summary Get building by ID
// @Description Retrieve a building by its ID
// @Tags buildings
// @Accept json
// @Produce json
// @Param id path int true "Building ID"
// @Success 200 {object} BuildingDTO
// @Router /api/v1/buildings/{id} [get]
func (h *Handler) GetBuilding(c *gin.Context) {
    // implementation
}
```

Generate documentation:

```bash
make swagger
```

## Database Development

### Migrations

Database schema is managed through GORM AutoMigrate. Add new fields to domain entities and restart the application.

### Database Operations

Database admin interface: http://localhost:8081 (root/example)

Reset database:
```bash
cd docker
docker compose down -v --remove-orphans && docker compose up -d
```

## Infrastructure

Infrastructure is managed with Terraform and Terragrunt:

```bash
cd infrastructure/terraform
# See README.md for detailed commands
```

## Code Standards

### Formatting and Linting

```bash
make format     # Format all Go code
make lint       # Run golangci-lint
```

### Commit Messages

Use conventional commit format:

```
feat: add building reservation endpoint
fix: resolve database connection timeout
docs: update API documentation
test: add integration tests for classes
```

### Code Review

- All changes require pull request review
- Ensure tests pass and coverage is maintained
- Update documentation for API changes
- Follow existing code patterns

## Contribution Process

1. **Fork** the repository
2. **Create** feature branch from `develop`
3. **Implement** changes with tests
4. **Run** quality checks: `make format lint test`
5. **Update** documentation if needed
6. **Submit** pull request to `develop` branch

## Getting Help

- Check existing documentation in `docs/`
- Review Docker setup in `docker/README.md`
- Check infrastructure setup in `infrastructure/terraform/README.md`
- Open issue for bugs or feature requests
