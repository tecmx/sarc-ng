# SARC-NG

SARC Next Generation - Resource Reservation and Management API

A modern Go-based API for managing resource reservations with support for buildings, classes, lessons, and resource allocation.

## Quick Start

### Prerequisites
- Go 1.24+
- Make
- Docker & Docker Compose

### Development Setup

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd sarc-ng
   make setup
   ```

2. **Start development**:
   ```bash
   # Native Go (fastest)
   make run

   # With hot reloading
   make debug

   # Full Docker environment (with database)
   cd docker && docker compose up -d
   ```

## Development Commands

### Go Development
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

### Docker Development
```bash
cd docker
docker compose up -d                        # Start development environment
docker compose down -v --remove-orphans     # Stop and cleanup all
```

### Infrastructure
```bash
cd infrastructure/terraform
# See README.md for Terragrunt commands
```

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
│   └── terraform/        # Terraform configurations
├── internal/              # Private application code
│   ├── adapter/          # External service adapters
│   ├── config/           # Configuration management
│   ├── domain/           # Business logic & entities
│   ├── service/          # Application services
│   └── transport/        # HTTP handlers & middleware
├── pkg/                   # Public libraries & client packages
└── test/                  # Test files
```

## Configuration

Environment variables:
```bash
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=sarcng

# Server
PORT=8080
ENVIRONMENT=development
```

Configuration files:
- `configs/default.yaml` - Base configuration
- `configs/development.yaml` - Development overrides

## API Endpoints

Base URL: `/api/v1/`

**Core Entities:**
- `/buildings` - Physical structures and locations
- `/classes` - Academic or training classes
- `/lessons` - Individual lesson instances
- `/resources` - Bookable resources (rooms, equipment)
- `/reservations` - Resource booking records

**Standard Operations:**
```
GET    /api/v1/{entity}      # List all
POST   /api/v1/{entity}      # Create new
GET    /api/v1/{entity}/:id  # Get by ID
PUT    /api/v1/{entity}/:id  # Update by ID
DELETE /api/v1/{entity}/:id  # Delete by ID
```

**API Documentation:** http://localhost:8080/swagger/index.html

## Testing

```bash
# Unit tests
make test

# Integration tests
cd docker && docker compose up -d && docker compose run --rm app go test -tags=integration ./test/integration/... && docker compose down -v --remove-orphans

# Coverage report
make coverage
```

## Building

```bash
# Development build
make build

# Production release
make release
```

## Documentation

- [CONTRIBUTING.md](docs/CONTRIBUTING.md) - Development guidelines and setup
- [docker/README.md](docker/README.md) - Docker environment documentation
- [infrastructure/terraform/README.md](infrastructure/terraform/README.md) - Infrastructure setup
