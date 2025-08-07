---
sidebar_position: 3
---

# Docker Environment

Docker configuration for SARC-NG development and deployment.

## Directory Structure

```
docker/
├── docker-compose.yml      # Development services
└── Dockerfile              # Multi-stage application build
```

## Services

### Development Stack

- **app**: Application server (port 8080)
- **db**: MySQL 8.0 database (port 3306)
- **adminer**: Database admin interface (port 8081)

### Monitoring Stack (Optional)

- **prometheus**: Metrics collection (port 8082)
- **grafana**: Visualization dashboards (port 8083)
- **node-exporter**: System metrics

## Quick Start

```bash
# Start development environment
docker compose up -d

# Stop and cleanup
docker compose down -v --remove-orphans
```

## Service Configuration

### Application (app)

```yaml
Ports:
  - 8080: API server
  - 2345: Debugger (development only)

Endpoints:
  - API: http://localhost:8080/api/v1
  - Swagger: http://localhost:8080/swagger/index.html
  - Health: http://localhost:8080/health
  - Metrics: http://localhost:8080/metrics
```

### Database (db)

```yaml
Port: 3306
Credentials:
  - User: root
  - Password: example (configurable via DB_PASSWORD)
  - Database: sarcng (configurable via DB_NAME)

Volumes:
  - Persistent data storage
  - Configuration files
```

### Database Admin (adminer)

```yaml
Port: 8081
Access:
  - URL: http://localhost:8081
  - System: MySQL
  - Server: db
  - Username: root
  - Password: example
```

## Environment Variables

```bash
# Database configuration
DB_PASSWORD=example     # Database password
DB_NAME=sarcng         # Database name
DB_USER=root           # Database user

# User mapping (Linux/macOS)
USER_ID=$(id -u)       # Host user ID
GROUP_ID=$(id -g)      # Host group ID

# Environment
ENVIRONMENT=development
```

## File Structure

### Dockerfile

Multi-stage build with:
- **Development stage**: Hot reloading with Air
- **Production stage**: Minimal distroless image

### docker-compose.yml

Development configuration:
- Volume mounts for live editing
- Development dependencies
- Debug ports exposed
- Host user mapping

## Development Workflow

### Hot Reloading

Code changes are automatically detected and the application reloads:

1. Edit source files on host
2. Changes sync to container via volumes
3. Air detects changes and rebuilds
4. Application restarts automatically

### Testing

```bash
# Integration tests
docker compose up -d && docker compose run --rm app go test -tags=integration ./test/integration/... && docker compose down -v --remove-orphans
```

### Database Access

Database admin interface available at: http://localhost:8081
- System: MySQL
- Server: db
- Username: root
- Password: example

## Production Build

```bash
# Build production image
docker build --target production -t sarc-ng:prod -f Dockerfile ../
```

## Troubleshooting

### Common Issues

**Permission problems (Linux/macOS):**
```bash
USER_ID=$(id -u) GROUP_ID=$(id -g) docker compose up
```

**Database connection failed:**
```bash
# Check database is running
docker compose ps db

# Check logs
docker compose logs db
```

**Port already in use:**
```bash
# Stop and cleanup
docker compose down -v --remove-orphans
```

### Cleanup

```bash
# Complete cleanup (removes containers, volumes, networks, and orphans)
docker compose down -v --remove-orphans
```

## Monitoring

### Metrics Collection

Application exposes Prometheus metrics at `/metrics`:
- HTTP request metrics
- Database operation metrics
- Application performance metrics

### Monitoring Stack

Start with monitoring enabled:
```bash
# Development with monitoring
docker compose --profile monitoring up -d

# Access dashboards
open http://localhost:8082  # Prometheus
open http://localhost:8083  # Grafana (admin/admin)
```
