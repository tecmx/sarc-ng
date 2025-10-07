---
sidebar_position: 3
tags:
  - development
---

# Development Guide

## Setup

```bash
# Clone and install
git clone https://github.com/tecmx/sarc-ng.git
cd sarc-ng
go mod download

# Start services
make docker-up

# Verify
curl http://localhost:8080/health
```

## Running the Application

### Docker (Recommended)
```bash
make docker-up      # Start all services
make docker-logs    # View logs
make docker-down    # Stop services
```

### Local Go
```bash
make run            # Run server
make debug          # Run with hot reload (requires air)
```

### Hybrid
```bash
# Database in Docker, app locally
cd infrastructure/docker && docker compose up -d db
cd ../.. && make run
```

## Development Commands

```bash
# Build & Run
make build          # Build binaries
make run            # Run server
make debug          # Hot reload dev mode

# Code Quality
make lint           # Run linters
make format         # Format code
make test           # Run tests
make coverage       # Test coverage

# Documentation
make swagger        # Generate API docs
make wire           # Generate DI code

# Workflow
make pre-commit     # Run all checks before commit
```

## Testing

### API Testing
```bash
# Health check
curl http://localhost:8080/health

# Get all buildings (with auth)
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/buildings

# Create building
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"name": "Library", "code": "LIB"}' \
  http://localhost:8080/api/v1/buildings

# Swagger UI
open http://localhost:8080/swagger/index.html
```

### Unit Tests
```bash
make test                              # All tests
go test -v ./internal/domain/...       # Specific package
go test -v -race ./...                 # With race detection
```

## Database

### Access
```bash
# Docker CLI
docker exec -it sarc-ng-db-dev mysql -u root -p
# Password: example

# Adminer UI
open http://localhost:8081
```

### Migrations
Auto-run on startup. Manual control:

```bash
go run cmd/cli/main.go migrate up
go run cmd/cli/main.go migrate status
go run cmd/cli/main.go migrate down
```

## CLI Commands

```bash
go run cmd/cli/main.go --help

# Examples
go run cmd/cli/main.go buildings create --name "Main" --code "MAIN"
go run cmd/cli/main.go resources list
go run cmd/cli/main.go auth generate-token --user-id 1
```

## Configuration

Environment variables override YAML config:

```bash
export DB_HOST=localhost
export DB_PASSWORD=mypass
export JWT_SECRET=mysecret
make run
```

Config files:
- `configs/default.yaml` - Base settings
- `configs/development.yaml` - Dev overrides

## Troubleshooting

### Database Connection
```bash
# Test connectivity
nc -zv localhost 3306

# Check container
docker ps
docker logs sarc-ng-db-dev
```

### Port Conflicts
```bash
# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>
```

### Clean Start
```bash
make docker-clean    # Remove all data
make docker-up       # Fresh start
```
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
