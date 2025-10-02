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
- **Go 1.19+** - [Download & Install Go](https://golang.org/dl/)
- **Docker & Docker Compose** - [Get Docker](https://docs.docker.com/get-docker/)
- **Git** - Version control
- **Make** - Build automation (usually pre-installed on Unix systems)

### Optional but Recommended
- **PostgreSQL** - For local database development
- **Redis** - For caching and sessions
- **AWS CLI** - For cloud deployment
- **Terraform** - For infrastructure management

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
   cd infrastructure/docker && docker compose up -d
   ```
   This starts:
   - PostgreSQL database
   - Redis cache
   - SARC-NG API server
   - All required services

4. **Verify the setup**
   ```bash
   # Check if API is running
   curl http://localhost:8080/v1/health

   # Should return: {"status": "ok"}
   ```

## Development Environment

### Using Docker (Recommended)

The easiest way to get started is using Docker Compose:

**Note:** All Docker commands should be run from the `infrastructure/docker/` directory.

```bash
# Navigate to docker directory
cd infrastructure/docker

# Start all services
docker compose up -d

# View logs
docker compose logs -f sarc-ng-server-dev

# Stop services
docker compose down
```

### Manual Setup

If you prefer to run services individually:

1. **Start PostgreSQL**
   ```bash
   # Using Docker
   docker run --name sarc-postgres -e POSTGRES_DB=sarc_ng -e POSTGRES_USER=sarc -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:15

   # Or install locally and create database
   createdb sarc_ng
   ```

2. **Start Redis**
   ```bash
   # Using Docker
   docker run --name sarc-redis -p 6379:6379 -d redis:7

   # Or install locally
   redis-server
   ```

3. **Configure environment**
   ```bash
   cp configs/development.yaml.example configs/development.yaml
   # Edit the configuration file with your database settings
   ```

4. **Run the API server**
   ```bash
   go run cmd/server/main.go
   ```

## Available Commands

The project includes several useful commands:

```bash
# Start API server
make run-server

# Start CLI tool
make run-cli

# Run tests
make test

# Build binaries
make build

# Generate API documentation
make docs

# Clean build artifacts
make clean

# Run linter
make lint
```

## Project Structure

```
sarc-ng/
├── cmd/                    # Application entrypoints
│   ├── cli/               # CLI application
│   ├── lambda/            # AWS Lambda function
│   └── server/            # HTTP API server
├── internal/              # Internal application code
│   ├── domain/            # Business logic
│   ├── service/           # Application services
│   └── transport/         # HTTP handlers
├── pkg/                   # Shared packages
├── configs/               # Configuration files
├── docker/                # Docker configurations
└── infrastructure/        # Terraform IaC
```

## Configuration

The application uses YAML configuration files in the `configs/` directory:

- `default.yaml` - Base configuration
- `development.yaml` - Development overrides
- `production.yaml` - Production settings (not in repo)

Key configuration sections:
- Database connection settings
- Redis configuration
- JWT authentication settings
- Server port and timeouts
- Logging configuration

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
# Verify PostgreSQL is running
docker ps | grep postgres

# Check connection
psql -h localhost -U sarc -d sarc_ng
```

**Go modules issues**
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download
```

Need help? Check our [GitHub Issues](https://github.com/tecmx/sarc-ng/issues) or start a [Discussion](https://github.com/tecmx/sarc-ng/discussions).
