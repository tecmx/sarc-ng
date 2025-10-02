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

#### Option A: Full Docker Environment (Recommended)

```bash
# Start all services
make docker-up

# Follow logs
make docker-logs

# Follow specific service logs
make docker-logs service=app

# Restart services
make docker-down && make docker-up
```

#### Option B: Hybrid Development

```bash
# Start only database
cd infrastructure/docker
docker compose up -d db

# Run API server locally (from project root)
make run
# or
go run cmd/server/main.go

# In another terminal, run CLI commands
go run cmd/cli/main.go --help
```

#### Option C: Full Local Setup

```bash
# Start MySQL (adjust connection settings in config)
mysql.server start

# Run the server
make run
```

## Development Commands

### Server Management

```bash
# Start HTTP server
make run

# Start server with hot reload (requires air)
make debug

# Generate Wire dependency injection code
make wire

# Build server binary
make build
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

# Run tests with coverage report
make coverage

# Run specific test package
go test ./internal/domain/building/...

# Run tests with verbose output
go test -v -race ./...
```

### Code Quality

```bash
# Run linter
make lint

# Format code
make format

# Generate Swagger documentation
make swagger

# Run pre-commit checks (format, wire, lint, test)
make pre-commit

# Run CI pipeline (wire, lint, test, build)
make ci
```

## API Development

### Testing API Endpoints

```bash
# Health check
curl http://localhost:8080/health

# Get buildings
curl -H "Authorization: Bearer <token>" http://localhost:8080/api/v1/buildings

# Create a building
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name": "Library", "code": "LIB"}' \
  http://localhost:8080/api/v1/buildings

# View Swagger documentation
open http://localhost:8080/swagger/index.html

# Access API documentation JSON
curl http://localhost:8080/swagger/doc.json
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
  http://localhost:8080/api/v1/auth/login
```

## Database Operations

### Direct Database Access

```bash
# Connect to development database
mysql -h localhost -u root -p sarcng

# Or using Docker
docker exec -it sarc-ng-db-dev mysql -u root -pexample sarcng
```

### Common Queries

```sql
-- List all tables
SHOW TABLES;

-- Check buildings
SELECT * FROM buildings LIMIT 5;

-- Check database structure
DESCRIBE buildings;

-- Reset data (development only)
SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE buildings;
TRUNCATE TABLE classes;
TRUNCATE TABLE resources;
TRUNCATE TABLE reservations;
TRUNCATE TABLE lessons;
SET FOREIGN_KEY_CHECKS = 1;
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
# Monitor slow queries (MySQL)
docker exec -it sarc-ng-db-dev mysql -u root -pexample -e "
SELECT * FROM information_schema.processlist
WHERE command != 'Sleep'
ORDER BY time DESC;"

# Check database size
docker exec -it sarc-ng-db-dev mysql -u root -pexample -e "
SELECT
    table_schema AS 'Database',
    ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) AS 'Size (MB)'
FROM information_schema.tables
WHERE table_schema = 'sarcng'
GROUP BY table_schema;"
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
docker ps | grep mysql

# Test connection
mysql -h localhost -u root -p -e "SELECT 1;"

# Check database logs
docker logs sarc-ng-db-dev
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

- **Postman/Insomnia/Thunder Client** - API testing
- **Adminer** - MySQL web GUI (included in Docker stack at <http://localhost:8081>)
- **MySQL Workbench** - Desktop MySQL GUI
- **Swagger UI** - Built-in API documentation (<http://localhost:8080/swagger/index.html>)
- **Air** - Hot reloading for Go applications
- **Delve** - Go debugger

## Contributing Workflow

1. Create feature branch: `git checkout -b feature/your-feature`
2. Make changes and test locally
3. Run tests: `make test`
4. Run linting: `make lint`
5. Commit changes: `git commit -m "feat: your feature"`
6. Push branch: `git push origin feature/your-feature`
7. Create Pull Request

For more details, see the [Contributing Guidelines](https://github.com/tecmx/sarc-ng/blob/main/CONTRIBUTING.md).
