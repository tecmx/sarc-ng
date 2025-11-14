# Docker Development Guide

Complete guide for Docker-based local development.

## Services

- API: <http://localhost:8080/api/v1>
- Swagger: <http://localhost:8080/swagger/index.html>
- DB Admin: <http://localhost:8081>
- Metrics: <http://localhost:8080/metrics>

## Commands

```bash
make docker-up       # Start services
make docker-down     # Stop services
make docker-logs     # View logs
make docker-clean    # Remove all data
```

## Manual Commands

```bash
cd infrastructure/docker

# Start
docker compose up -d --wait

# Stop
docker compose down

# Stop and remove data
docker compose down -v

# Logs
docker compose logs -f app

# Database access
docker exec -it sarc-ng-db-dev mysql -u root -p
```

## Troubleshooting

### Check Status

```bash
docker ps -a
```

### View Logs

```bash
docker logs sarc-ng-server-dev
```

### Rebuild

```bash
docker compose down -v
docker compose build --no-cache
docker compose up -d --wait
```

### Database Access

```bash
docker exec -it sarc-ng-db-dev mysql -u root -p
```

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

