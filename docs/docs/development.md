---
sidebar_position: 3
tags:
  - development
  - local
  - workflow
---

# Local Development

This guide covers the development workflow for building and testing SARC-NG locally.

## Development Workflow

### 1. Environment Setup

```bash
# Set up your development environment
export SARC_ENV=development
export SARC_CONFIG_PATH=./configs/development.yaml

# Or create a .env file
echo "SARC_ENV=development" > .env
echo "SARC_CONFIG_PATH=./configs/development.yaml" >> .env
```

### 2. Database Migrations

The application automatically handles database migrations on startup. For manual control:

```bash
# Run migrations manually
go run cmd/cli/main.go migrate up

# Check migration status
go run cmd/cli/main.go migrate status

# Rollback migrations
go run cmd/cli/main.go migrate down
```

### 3. Running Services

#### Option A: Full Docker Environment
**Note:** All Docker commands should be run from the `docker/` directory.

```bash
# Navigate to docker directory
cd docker

# Start all services (recommended)
docker compose up -d

# Follow logs
docker compose logs -f sarc-ng-server-dev

# Restart just the API server
docker compose restart sarc-ng-server-dev
```

#### Option B: Hybrid Development
```bash
# Start only infrastructure
docker compose up -d postgres redis

# Run API server locally
go run cmd/server/main.go

# In another terminal, run CLI commands
go run cmd/cli/main.go --help
```

#### Option C: Full Local Setup
```bash
# Start PostgreSQL (adjust connection settings in config)
postgres -D /usr/local/var/postgres

# Start Redis
redis-server

# Run the server
make run-server
```

## Development Commands

### Server Management
```bash
# Start HTTP server
make run-server

# Start server with hot reload (requires air)
air

# Start server with custom config
SARC_CONFIG_PATH=./configs/custom.yaml go run cmd/server/main.go
```

### CLI Operations
```bash
# View all CLI commands
go run cmd/cli/main.go --help

# Create a building
go run cmd/cli/main.go buildings create --name "Main Building" --code "MAIN"

# List resources
go run cmd/cli/main.go resources list

# Import data from CSV
go run cmd/cli/main.go import --file data.csv --type buildings
```

### Testing
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration

# Run specific test package
go test ./internal/domain/building/...

# Run tests with verbose output
go test -v ./...
```

### Code Quality
```bash
# Run linter
make lint

# Format code
make fmt

# Run security checks
make security

# Generate mocks (if using mockery)
make generate-mocks
```

## API Development

### Testing API Endpoints

```bash
# Health check
curl http://localhost:8080/v1/health

# Get buildings
curl -H "Authorization: Bearer <token>" http://localhost:8080/v1/buildings

# Create a building
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name": "Library", "code": "LIB"}' \
  http://localhost:8080/v1/buildings

# Get API documentation
curl http://localhost:8080/v1/docs
```

### Authentication

For development, you can generate a test JWT token:

```bash
# Generate a development token
go run cmd/cli/main.go auth generate-token --user-id 1 --role admin

# Or use the development token endpoint (if enabled)
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin"}' \
  http://localhost:8080/v1/auth/login
```

## Database Operations

### Direct Database Access

```bash
# Connect to development database
psql -h localhost -U sarc -d sarc_ng

# Or using Docker
docker exec -it sarc-postgres psql -U sarc -d sarc_ng
```

### Common Queries

```sql
-- List all tables
\dt

-- Check buildings
SELECT * FROM buildings LIMIT 5;

-- Check migrations
SELECT * FROM schema_migrations;

-- Reset data (development only)
TRUNCATE buildings, classes, resources, reservations, lessons RESTART IDENTITY CASCADE;
```

### Seeding Development Data

```bash
# Load sample data
go run cmd/cli/main.go seed --env development

# Load specific data sets
go run cmd/cli/main.go seed --type buildings
go run cmd/cli/main.go seed --type users
```

## Debugging

### Server Debugging
```bash
# Run with debug logging
SARC_LOG_LEVEL=debug go run cmd/server/main.go

# Run with profiling enabled
SARC_PPROF_ENABLED=true go run cmd/server/main.go

# Access profiling endpoints
curl http://localhost:8080/debug/pprof/
```

### Using Delve Debugger
```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug server
dlv debug cmd/server/main.go

# Debug with arguments
dlv debug cmd/cli/main.go -- buildings list
```

### VS Code Configuration

Create `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Server",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/server/main.go",
            "env": {
                "SARC_ENV": "development"
            }
        },
        {
            "name": "Launch CLI",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/cli/main.go",
            "args": ["--help"]
        }
    ]
}
```

## Hot Reload Development

### Using Air
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Create .air.toml config
air init

# Start with hot reload
air
```

### Using Watchman (Alternative)
```bash
# Install watchman
brew install watchman  # macOS
# or use your system's package manager

# Watch for changes and restart
watchman-make -p '**/*.go' -r 'make run-server'
```

## Performance Monitoring

### Local Metrics
```bash
# Access metrics endpoint (if enabled)
curl http://localhost:8080/metrics

# View health metrics
curl http://localhost:8080/v1/health/detailed
```

### Database Performance
```bash
# Monitor slow queries (PostgreSQL)
docker exec -it sarc-postgres psql -U sarc -d sarc_ng -c "
SELECT query, calls, total_time, mean_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;"
```

## Troubleshooting

### Common Development Issues

**"Port already in use" Error**
```bash
# Find and kill process using port 8080
lsof -ti:8080 | xargs kill -9

# Or use a different port
SARC_SERVER_PORT=8081 go run cmd/server/main.go
```

**Database Connection Issues**
```bash
# Verify database is running
docker ps | grep postgres

# Test connection
pg_isready -h localhost -p 5432

# Check database logs
docker logs sarc-postgres
```

**Module Issues**
```bash
# Clean and reinstall modules
go clean -modcache
go mod download
go mod tidy
```

**Build Issues**
```bash
# Clean build cache
go clean -cache
make clean
make build
```

### Useful Development Tools

- **Postman/Insomnia** - API testing
- **pgAdmin** - PostgreSQL GUI
- **Redis Insight** - Redis GUI
- **Grafana** - Metrics visualization
- **Jaeger** - Distributed tracing (if implemented)

## Contributing Workflow

1. Create feature branch: `git checkout -b feature/your-feature`
2. Make changes and test locally
3. Run tests: `make test`
4. Run linting: `make lint`
5. Commit changes: `git commit -m "feat: your feature"`
6. Push branch: `git push origin feature/your-feature`
7. Create Pull Request

For more details, see the [Contributing Guidelines](https://github.com/tecmx/sarc-ng/blob/main/CONTRIBUTING.md).
