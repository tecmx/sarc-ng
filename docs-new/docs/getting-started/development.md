---
sidebar_position: 2
---

# Development

This guide covers the development workflow and tools for SARC-NG.

## Development Workflow

SARC-NG supports multiple development workflows to suit different preferences.

### Native Go Development

For the fastest development cycle, run the application directly with Go:

```bash
# Run the application
make run

# With hot reloading
make debug
```

Hot reloading with `make debug` uses the [Air](https://github.com/cosmtrek/air) tool to automatically rebuild and restart the application when source files change.

### Docker Development

For a full development environment with database and other services:

```bash
cd docker
docker compose up -d
```

This will:
- Start a MySQL database on port 3306
- Run the application on port 8080
- Provide a database admin interface on port 8081
- Set up development volumes for live code editing

To stop the environment:

```bash
docker compose down -v --remove-orphans
```

## Development Commands

### Core Commands

```bash
make setup        # Install dependencies and tools
make run          # Run the application
make debug        # Run with hot reloading
make test         # Run tests
make build        # Build applications
make lint         # Run linters
make format       # Format code
make coverage     # Generate test coverage
make swagger      # Generate API documentation
make clean        # Remove build artifacts
```

### Wire Dependency Injection

After modifying dependency injection in `cmd/*/wire.go`:

```bash
make wire       # Regenerate injection code
```

### Testing

```bash
# Unit tests
make test

# Coverage report
make coverage

# Integration tests (requires Docker)
cd docker && docker compose up -d && \
docker compose run --rm app go test -tags=integration ./test/integration/... && \
docker compose down -v --remove-orphans
```

## Development Tools

The following tools are automatically installed via `make setup`:

- **Air**: Hot reloading during development
- **Wire**: Dependency injection code generation
- **Swag**: API documentation generation
- **golangci-lint**: Code linting and static analysis
- **Delve**: Go debugger

## Project Structure

```
sarc-ng/
├── api/                    # API documentation (Swagger)
├── cmd/                    # Application entry points
│   ├── cli/               # CLI application
│   ├── lambda/            # AWS Lambda entry point
│   └── server/            # HTTP server
├── configs/               # Configuration files
├── docker/                # Docker development & deployment
├── docs/                  # Project documentation
├── infrastructure/        # Infrastructure as Code
│   └── aws/              # AWS infrastructure configurations
├── internal/              # Private application code
│   ├── adapter/          # External service adapters
│   ├── config/           # Configuration management
│   ├── domain/           # Business logic & entities
│   ├── service/          # Application services
│   └── transport/        # HTTP handlers & middleware
├── pkg/                   # Public libraries & client packages
└── test/                  # Test files
```

## Architecture Layers

SARC-NG follows Clean Architecture principles with clear separation of concerns:

1. **Domain Layer** (`internal/domain/`): Core business entities and logic
2. **Service Layer** (`internal/service/`): Application services and use cases
3. **Transport Layer** (`internal/transport/`): HTTP handlers and routing
4. **Adapter Layer** (`internal/adapter/`): External service integration

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

Database schema is managed through GORM AutoMigrate. Add new fields to domain entities and restart the application.

Database admin interface: http://localhost:8081
- System: MySQL
- Server: db
- Username: root
- Password: example

Reset database:
```bash
cd docker
docker compose down -v --remove-orphans && docker compose up -d
```

## Next Steps

- [Learn about Docker setup](docker)
- [Follow the quick start guide](quick-start)
- [Explore the API documentation](/api/overview)
