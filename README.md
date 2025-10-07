# SARC-NG

[![CI](https://github.com/tecmx/sarc-ng/workflows/CI/badge.svg)](https://github.com/tecmx/sarc-ng/actions/workflows/ci.yml)
[![Deploy](https://github.com/tecmx/sarc-ng/workflows/Deploy/badge.svg)](https://github.com/tecmx/sarc-ng/actions/workflows/deploy.yml)
[![Release](https://github.com/tecmx/sarc-ng/workflows/Release/badge.svg)](https://github.com/tecmx/sarc-ng/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tecmx/sarc-ng)](https://goreportcard.com/report/github.com/tecmx/sarc-ng)

Resource Management and Scheduling API - A modern Go-based system for managing buildings, classes, lessons, and resource reservations.

## Quick Start

```bash
git clone <repository-url>
cd sarc-ng
make docker-up
```

**Access:**
- API: http://localhost:8080/api/v1
- Swagger: http://localhost:8080/swagger/index.html

## Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Make

## Commands

```bash
# Development
make run            # Run locally
make debug          # Hot reload mode
make test           # Run tests

# Build
make build          # Build binaries
make release        # Production build

# Docker
make docker-up      # Start services
make docker-down    # Stop services
make docker-logs    # View logs

# Code Quality
make lint           # Run linters
make format         # Format code
make swagger        # Generate API docs
```

## API Endpoints

Base: `/api/v1/`

**Entities:** `buildings`, `classes`, `lessons`, `resources`, `reservations`

**Operations:**
```
GET    /api/v1/{entity}        # List
POST   /api/v1/{entity}        # Create
GET    /api/v1/{entity}/:id    # Get
PUT    /api/v1/{entity}/:id    # Update
DELETE /api/v1/{entity}/:id    # Delete
```

## Configuration

Environment variables:
```bash
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=sarcng
PORT=8080
```

Config files:
- `configs/default.yaml`
- `configs/development.yaml`

## Project Structure

```
cmd/            # Entry points (cli, lambda, server)
internal/       # Application code
  ├── domain/   # Business logic
  ├── service/  # Services
  ├── adapter/  # External adapters
  └── transport/# HTTP handlers
infrastructure/ # Docker, SAM, Terraform
```

## Documentation

See `docs/content/` for detailed guides:

- [Getting Started](docs/content/getting-started.md)
- [Development](docs/content/development.md)
- [Architecture](docs/content/architecture.md)
- [Deployment](docs/content/deployment.md)
- [CI/CD Pipeline](.github/CICD.md)

## CI/CD

Automated workflows for testing, building, and deploying:

- **CI:** Runs on every push/PR (lint, test, build, security scan)
- **Deploy:** Auto-deploys to AWS Lambda on push to `main`/`develop`
- **Release:** Creates multi-platform binaries on version tags

See [CI/CD Documentation](.github/CICD.md) for details.
