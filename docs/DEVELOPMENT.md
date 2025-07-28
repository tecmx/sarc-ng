# Development

Local development setup, workflows, and practices for SARC-NG.

## Prerequisites

Required tools and versions:

- **Go**: 1.24+ with module support
- **Docker**: Latest stable with Docker Compose
- **Make**: For development commands
- **Git**: Version control

Optional but recommended:
- **LocalStack**: For AWS service simulation
- **Terraform**: 1.0+ for infrastructure development
- **Terragrunt**: 0.45+ for infrastructure management

## Quick Setup

### 1. Initial Setup

```bash
# Clone repository
git clone <repository-url>
cd sarc-ng

# Install dependencies and development tools
make setup
```

The `make setup` command installs:
- Air (hot reloading)
- Wire (dependency injection)
- Swag (API documentation)
- golangci-lint (code quality)
- Delve (debugging)

### 2. Choose Development Mode

**Option A: Native Go (Fastest)**
```bash
make run        # Direct execution
make debug      # With hot reloading via Air
```

**Option B: Docker Environment (Full Stack)**
```bash
docker compose up -d    # Start all services
```

## Development Workflows

### Code Development

**Hot Reloading Development:**
```bash
make debug
# Application restarts automatically on file changes
# Debugger available on port 2345
```

**Standard Development:**
```bash
make run
# Fast startup for quick testing
```

### Testing Workflows

**Unit Tests:**
```bash
make test           # Run all unit tests
make coverage       # Generate coverage report
```

**Integration Tests:**
```bash
# Start environment, run tests, cleanup
docker compose up -d && \
docker compose run --rm app go test -tags=integration ./test/integration/... && \
docker compose down -v --remove-orphans
```

**Test Categories:**
- Unit tests: Fast, isolated, no external dependencies
- Integration tests: Real database, full request lifecycle
- Build tag `integration` separates test types

### Code Quality

**Pre-commit Checks:**
```bash
make format     # Format code with gofmt and goimports
make lint       # Run golangci-lint
make test       # Verify all tests pass
```

**Complete Quality Check:**
```bash
make format && make lint && make test
```

### API Development

**Generate Swagger Documentation:**
```bash
make swagger
# View at: http://localhost:8080/swagger/index.html
```

**API Testing:**
```bash
# Manual testing
curl http://localhost:8080/api/v1/buildings

# View API docs
open http://localhost:8080/swagger/index.html
```

### Database Development

**Schema Changes:**
1. Modify domain entities in `internal/domain/*/entity.go`
2. Restart application (GORM AutoMigrate handles schema)
3. Verify changes in database admin panel

**Database Access:**
- Admin panel: http://localhost:8081 (admin/example)
- Direct connection: MySQL port 3306
- Credentials: root/example (configurable)

### Building and Deployment

**Development Builds:**
```bash
make build      # Build server and CLI binaries
```

**Production Builds:**
```bash
make release    # Optimized builds with version info
```

**Docker Images:**
```bash
# Development image
docker build --target development -t sarc-ng:dev .

# Production image  
docker build --target production -t sarc-ng:prod .
```

## Project Structure

### Code Organization

```
sarc-ng/
├── cmd/                    # Application entry points
│   ├── cli/               # CLI commands and formatters
│   ├── server/            # HTTP server main
│   └── lambda/            # AWS Lambda handler
├── internal/              # Private application code
│   ├── adapter/          # External integrations
│   │   ├── db/           # Database configuration
│   │   └── gorm/         # GORM repository implementations
│   ├── config/           # Configuration management
│   ├── domain/           # Business entities and interfaces
│   │   ├── building/     # Building domain logic
│   │   ├── class/        # Class domain logic
│   │   ├── lesson/       # Lesson domain logic
│   │   ├── reservation/  # Reservation domain logic
│   │   └── resource/     # Resource domain logic
│   ├── service/          # Application services
│   └── transport/        # HTTP handlers and routing
├── pkg/                   # Public libraries
│   ├── rest/             # REST client and middleware
│   └── metrics/          # Metrics collection
├── configs/               # Configuration files
├── docker/                # Docker development setup
└── test/                  # Test files and utilities
```

### Key Patterns

**Domain Structure:**
Each domain package contains:
- `entity.go`: Business entity with validation rules
- `usecase.go`: Business operation interface
- `repository.go`: Data access interface

**Service Layer:**
- Implements use case interfaces
- Contains business logic and validation
- Coordinates between repositories

**Transport Layer:**
- HTTP handlers with DTO mapping
- Input validation and error handling
- API documentation annotations

## Development Practices

### Clean Architecture

**Dependency Flow:**
```
Transport → Service → Domain ← Adapter
```

**Key Rules:**
- Domain layer has no external dependencies
- Services depend only on domain interfaces
- Adapters implement domain interfaces
- Transport layer depends on services

### Code Style

**Naming Conventions:**
- Interfaces: Descriptive nouns (e.g., `Repository`, `Usecase`)
- Implementations: Context + Interface (e.g., `GormAdapter`, `Service`)
- Files: Purpose-based (e.g., `entity.go`, `handler.go`)

**Package Organization:**
- One concept per package
- Small, focused interfaces
- Clear public/private boundaries

### Testing Strategy

**Unit Testing:**
- Test business logic in isolation
- Mock external dependencies
- Table-driven tests for multiple scenarios
- High coverage for domain and service layers

**Integration Testing:**
- Test complete request flows
- Use real database via Docker
- Verify API contracts and responses
- Test error scenarios and edge cases

### Configuration Management

**Configuration Hierarchy:**
1. Default values in `configs/default.yaml`
2. Environment overrides in `configs/{env}.yaml`
3. Environment variables for sensitive data

**Environment Variables:**
```bash
# Database configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=sarcng

# Application settings
PORT=8080
ENVIRONMENT=development
GIN_MODE=debug
```

### Dependency Injection

**Wire Setup:**
1. Define providers in `wire.go`
2. Generate injection code: `make wire`
3. Use injected dependencies in main functions

**Adding New Dependencies:**
1. Add provider function to `wire.go`
2. Add to provider set
3. Regenerate with `make wire`
4. Update main function to use new dependency

### Error Handling

**Patterns:**
- Return errors, don't panic in library code
- Use custom error types for business logic
- Log errors at the boundary (transport layer)
- Provide meaningful error messages to users

**HTTP Error Responses:**
- Consistent error format across all endpoints
- Appropriate HTTP status codes
- No sensitive information in error messages

## Development Tools

### Make Commands

```bash
# Setup and dependencies
make setup          # Install all development tools
make install        # Install Go dependencies only

# Development
make run            # Run server directly
make debug          # Run with hot reloading
make build          # Build applications
make release        # Production build with version info

# Code quality
make test           # Run unit tests
make coverage       # Generate test coverage report
make lint           # Run linters
make format         # Format code

# Documentation
make swagger        # Generate API documentation

# Infrastructure
make wire           # Generate dependency injection
make clean          # Remove build artifacts
```

### IDE Setup

**VS Code Configuration:**
- Install Go extension
- Configure debug launch configurations
- Set up test runner integration
- Enable format on save

**Debugging:**
- Delve debugger integration
- Docker debugger support on port 2345
- Breakpoint debugging in development mode

### Local Services

**Development Stack:**
- **App Server**: http://localhost:8080
- **API Documentation**: http://localhost:8080/swagger/index.html
- **Database Admin**: http://localhost:8081
- **Metrics**: http://localhost:8080/metrics (Prometheus format)
- **Health Check**: http://localhost:8080/health

**Service Management:**
```bash
# Start development environment
docker compose up -d

# View logs
docker compose logs -f app

# Stop and cleanup
docker compose down -v --remove-orphans
```

## Troubleshooting

### Common Issues

**Port Already in Use:**
```bash
# Find process using port 8080
lsof -i :8080
# Kill process or use different port
```

**Database Connection Errors:**
```bash
# Ensure Docker services are running
docker compose ps
# Check database logs
docker compose logs db
```

**Hot Reload Not Working:**
```bash
# Ensure Air is installed
make setup
# Check .air.toml configuration
# Verify file permissions
```

**Wire Generation Fails:**
```bash
# Install Wire tool
go install github.com/google/wire/cmd/wire@latest
# Regenerate injection code
make wire
```

### Performance Optimization

**Development Speed:**
- Use native Go for fastest iteration
- Docker only when full environment needed
- Optimize imports and dependencies

**Memory Usage:**
- Monitor with `go tool pprof`
- Check database connection pools
- Profile with integrated tools

### Best Practices

**Before Committing:**
1. Run `make format lint test`
2. Verify API documentation updates
3. Test integration scenarios
4. Update relevant documentation

**Adding New Features:**
1. Start with domain entity and interfaces
2. Implement service layer with business logic
3. Add repository implementation
4. Create transport layer handlers
5. Add comprehensive tests
6. Update API documentation 