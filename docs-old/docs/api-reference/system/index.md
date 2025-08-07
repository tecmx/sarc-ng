---
sidebar_position: 6
---

# System API

System-level API endpoints for monitoring and managing the SARC-NG application.

## Health Check

Check the health status of the API.

```http
GET /health
```

### Response

```json
{
  "status": "healthy",
  "service": "sarc-ng",
  "timestamp": "2024-01-01T00:00:00Z",
  "version": "1.0.0",
  "uptime": "10h32m15s"
}
```

## Metrics

Prometheus-compatible metrics endpoint.

```http
GET /metrics
```

### Response

Plain text Prometheus metrics:

```
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",path="/api/v1/buildings",status="200"} 42
http_requests_total{method="POST",path="/api/v1/buildings",status="201"} 5
http_requests_total{method="GET",path="/api/v1/resources",status="200"} 37

# HELP http_request_duration_seconds HTTP request duration in seconds
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{method="GET",path="/api/v1/buildings",le="0.1"} 35
http_request_duration_seconds_bucket{method="GET",path="/api/v1/buildings",le="0.5"} 40
http_request_duration_seconds_bucket{method="GET",path="/api/v1/buildings",le="1"} 41
http_request_duration_seconds_bucket{method="GET",path="/api/v1/buildings",le="5"} 42
http_request_duration_seconds_bucket{method="GET",path="/api/v1/buildings",le="+Inf"} 42

# HELP database_queries_total Total number of database queries
# TYPE database_queries_total counter
database_queries_total{operation="select",table="buildings"} 105
database_queries_total{operation="insert",table="buildings"} 5
database_queries_total{operation="update",table="buildings"} 2
database_queries_total{operation="delete",table="buildings"} 1
```

## API Version

Get information about the current API version.

```http
GET /api/version
```

### Response

```json
{
  "version": "1.0.0",
  "buildDate": "2024-01-01T00:00:00Z",
  "gitCommit": "a1b2c3d4",
  "goVersion": "go1.24"
}
```

## Authentication Status

Get information about the current authentication configuration.

```http
GET /api/auth/status
```

### Response

```json
{
  "enabled": true,
  "provider": "jwt",
  "requiredPaths": [
    "/api/v1/buildings*",
    "/api/v1/resources*",
    "/api/v1/classes*",
    "/api/v1/lessons*",
    "/api/v1/reservations*"
  ],
  "publicPaths": [
    "/health",
    "/metrics",
    "/api/version",
    "/api/auth/status",
    "/api/auth/login"
  ]
}
```
