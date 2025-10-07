---
sidebar_position: 2
tags:
  - getting-started
---

# Getting Started

## Prerequisites

- **Go 1.24+** - [Download](https://golang.org/dl/)
- **Docker & Docker Compose** - [Get Docker](https://docs.docker.com/get-docker/)
- **Make** - Usually pre-installed on Unix systems

## Quick Start

```bash
# Clone repository
git clone https://github.com/tecmx/sarc-ng.git
cd sarc-ng

# Start with Docker
make docker-up

# Verify
curl http://localhost:8080/health
# Response: {"service":"sarc-ng","status":"healthy"}
```

**Access Points:**
- API: http://localhost:8080/api/v1
- Swagger: http://localhost:8080/swagger/index.html
- DB Admin: http://localhost:8081
- Metrics: http://localhost:8080/metrics

## Manual Setup (Optional)

### Start MySQL
```bash
docker run --name sarc-mysql \
  -e MYSQL_ROOT_PASSWORD=example \
  -e MYSQL_DATABASE=sarcng \
  -p 3306:3306 \
  -d mysql:8.0
```

### Run API Server
```bash
make run
# or
go run cmd/server/main.go
```

### Configuration
Override via environment variables:
```bash
DB_HOST=localhost DB_PASSWORD=mypass make run
```

Config files:
- `configs/default.yaml` - Base configuration
- `configs/development.yaml` - Development overrides

## Available Commands

```bash
# Development
make run            # Run API server
make debug          # Hot reload mode
make wire           # Generate dependency injection

# Build
make build          # Build binaries
make release        # Production build

# Testing
make test           # Run tests
make coverage       # Coverage report
make lint           # Run linters
make format         # Format code

# Docker
make docker-up      # Start services
make docker-down    # Stop services
make docker-logs    # View logs
make docker-clean   # Remove data

# Workflow
make pre-commit     # All pre-commit checks
make check          # Quick validation
```

## Project Structure

```
cmd/                # Entry points (cli, lambda, server)
internal/           # Private code
  ├── domain/       # Business logic & entities
  ├── service/      # Application services
  ├── adapter/      # External adapters (GORM, etc)
  └── transport/    # HTTP handlers
pkg/                # Public packages
api/                # OpenAPI specifications
configs/            # Configuration files
infrastructure/     # Docker, SAM, Terraform
test/               # Integration tests
```

## Configuration Priority

1. **Environment Variables** (highest)
2. **Environment YAML** (`development.yaml`)
3. **Base YAML** (`default.yaml`)
4. **Code Defaults** (lowest)

## Next Steps

- Browse [API Reference](http://localhost:8080/swagger/index.html)
- Read [Development Guide](development.md)
- Explore [Architecture](architecture.md)
- Learn about [Deployment](deployment.md)

### Configuration Files

- `configs/default.yaml` - Base configuration for all environments
- `configs/development.yaml` - Development-specific overrides

### Environment Variables

Standard database environment variables are automatically mapped:

```bash
export DB_HOST=localhost        # Database host
export DB_PORT=3306             # Database port
export DB_USER=root             # Database user
export DB_PASSWORD=password     # Database password
export DB_NAME=sarcng           # Database name
export PORT=8080                # Server port
export ENVIRONMENT=development  # Environment (dev/staging/prod)
```

### Key Configuration Sections

- **Server:** Port, host, timeouts
- **Database:** MySQL connection settings and pool configuration
- **JWT:** Authentication secret and token expiration
- **Logging:** Level, format, output
- **CORS:** Cross-origin resource sharing settings
- **Swagger:** API documentation configuration

## Next Steps

- 📖 Read the [Local Development](./development) guide for detailed development workflows
- 🏗️ Check [Architecture Overview](./architecture) to understand the project structure
- 🚀 Learn about [Deployment](./deployment) for cloud deployment

## Troubleshooting

### Common Issues

**Port already in use**

```bash
# Check what's using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

**Database connection failed**

```bash
# Verify MySQL is running
docker ps | grep mysql

# Check connection
mysql -h localhost -u root -p sarcng
```

**Go modules issues**

```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download
```

Need help? Check our [GitHub Issues](https://github.com/tecmx/sarc-ng/issues) or start a [Discussion](https://github.com/tecmx/sarc-ng/discussions).
