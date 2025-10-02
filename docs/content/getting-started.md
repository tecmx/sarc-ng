---
sidebar_position: 2
tags:
  - getting-started
  - backend
---

# Getting Started

This guide will help you get up and running with the SARC-NG backend development environment.

## Prerequisites

Before you begin, ensure you have the following installed:

### Required

- **Go 1.24+** - [Download & Install Go](https://golang.org/dl/)
- **Docker & Docker Compose** - [Get Docker](https://docs.docker.com/get-docker/)
- **Git** - Version control
- **Make** - Build automation (usually pre-installed on Unix systems)

### Optional but Recommended

- **MySQL 8.0** - For local database development without Docker
- **AWS CLI** - For cloud deployment (Lambda/RDS)
- **AWS SAM CLI** - For serverless deployment
- **Terraform/Terragrunt** - For infrastructure management

## Quick Start

1. **Clone the repository**

   ```bash
   git clone https://github.com/tecmx/sarc-ng.git
   cd sarc-ng
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Start with Docker Compose** (Recommended for first-time setup)

   ```bash
   make docker-up
   ```

   This starts:
   - MySQL 8.0 database
   - SARC-NG API server
   - Adminer (database management UI)

4. **Verify the setup**

   ```bash
   # Check if API is running
   curl http://localhost:8080/health

   # Should return: {"service":"sarc-ng","status":"healthy"}

   # Access Swagger documentation
   open http://localhost:8080/swagger/index.html
   ```

## Development Environment

### Using Docker (Recommended)

The easiest way to get started is using Docker Compose via Make commands:

```bash
# Start all services
make docker-up

# View logs
make docker-logs

# View specific service logs
make docker-logs service=app

# Stop services (keeps data)
make docker-down

# Remove all data (WARNING: deletes database)
make docker-clean
```

**Services Available:**

- API: `http://localhost:8080/api/v1`
- Swagger: `http://localhost:8080/swagger/index.html`
- DB Admin (Adminer): `http://localhost:8081`
- Metrics: `http://localhost:8080/metrics`

### Manual Setup

If you prefer to run services individually:

1. **Start MySQL**

   ```bash
   # Using Docker
   docker run --name sarc-mysql \
     -e MYSQL_ROOT_PASSWORD=example \
     -e MYSQL_DATABASE=sarcng \
     -p 3306:3306 \
     -d mysql:8.0

   # Or install locally and create database
   mysql -u root -p -e "CREATE DATABASE sarcng;"
   ```

2. **Run the API server**

   ```bash
   make run
   # or
   go run cmd/server/main.go
   ```

   **Configuration:**
   - Configuration is managed via `configs/default.yaml` and `configs/development.yaml`
   - Override settings using environment variables: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
   - Example: `DB_HOST=localhost DB_PASSWORD=mypass make run`

## Available Commands

The project includes several useful commands:

```bash
# Development
make run            # Run API server locally
make debug          # Run with hot reloading (requires air)
make wire           # Generate dependency injection code

# Build & Release
make build          # Build server and CLI applications
make release        # Build production binaries

# Testing & Quality
make test           # Run tests
make coverage       # Generate test coverage report
make lint           # Run linters
make format         # Format Go code

# Documentation
make swagger        # Generate Swagger API documentation
make docs-build     # Build documentation site
make docs-serve     # Start documentation dev server

# Docker
make docker-up      # Start Docker services
make docker-down    # Stop Docker services
make docker-logs    # View logs
make docker-clean   # Remove all data

# Workflow
make pre-commit     # Run all pre-commit checks
make check          # Quick validation (lint, test)
make ci             # CI pipeline (wire, lint, test, build)
```

## Project Structure

```
sarc-ng/
├── cmd/                    # Application entrypoints
│   ├── cli/               # CLI application
│   ├── lambda/            # AWS Lambda function
│   └── server/            # HTTP API server
├── internal/              # Internal application code
│   ├── adapter/           # External service adapters (GORM, etc.)
│   ├── config/            # Configuration loader
│   ├── domain/            # Business logic (entities, repositories, use cases)
│   ├── service/           # Application services
│   └── transport/         # HTTP handlers (REST)
├── pkg/                   # Shared packages
│   ├── metrics/           # Metrics collection
│   └── rest/              # REST client utilities
├── api/                   # API specifications
│   └── swagger/           # Generated Swagger docs
├── configs/               # Configuration files
├── infrastructure/        # Infrastructure as Code
│   ├── docker/            # Docker Compose configurations
│   ├── sam/               # AWS SAM (Serverless) deployment
│   └── terraform/         # Terraform/Terragrunt IaC
└── test/                  # Integration tests
```

## Configuration

The application uses a hierarchical configuration system with the following priority (highest to lowest):

1. **Environment Variables** - Override any config setting
2. **Environment-Specific Config** - `configs/development.yaml` (based on `ENVIRONMENT` or `ENV` variable)
3. **Base Config** - `configs/default.yaml`
4. **Hardcoded Defaults** - In `internal/config/loader.go`

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
