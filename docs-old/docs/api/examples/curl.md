---
sidebar_position: 3
---

# cURL Examples

Examples of using cURL to interact with the SARC-NG API from the command line.

## Basic Usage

cURL is a command-line tool for making HTTP requests. It's useful for testing APIs and automation.

### Request Format

```bash
curl [options] <url>
```

Common options:
- `-X, --request`: HTTP method (GET, POST, PUT, DELETE)
- `-H, --header`: Add header to the request
- `-d, --data`: Send data in the request body
- `-v, --verbose`: Verbose output

## Buildings API

### List Buildings

```bash
# Get all buildings
curl -X GET http://localhost:8080/api/v1/buildings \
  -H "Content-Type: application/json"

# With pagination
curl -X GET "http://localhost:8080/api/v1/buildings?page=1&limit=10" \
  -H "Content-Type: application/json"

# With sorting
curl -X GET "http://localhost:8080/api/v1/buildings?sort=name&direction=asc" \
  -H "Content-Type: application/json"
```

### Get Building by ID

```bash
curl -X GET http://localhost:8080/api/v1/buildings/1 \
  -H "Content-Type: application/json"
```

### Create Building

```bash
curl -X POST http://localhost:8080/api/v1/buildings \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Engineering Building",
    "code": "ENG-A"
  }'
```

### Update Building

```bash
curl -X PUT http://localhost:8080/api/v1/buildings/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Engineering Building Updated",
    "code": "ENG-B"
  }'
```

### Delete Building

```bash
curl -X DELETE http://localhost:8080/api/v1/buildings/1 \
  -H "Content-Type: application/json"
```

## Resources API

### List Resources

```bash
# Get all resources
curl -X GET http://localhost:8080/api/v1/resources \
  -H "Content-Type: application/json"

# Filter by type
curl -X GET "http://localhost:8080/api/v1/resources?type=equipment" \
  -H "Content-Type: application/json"

# Filter by availability
curl -X GET "http://localhost:8080/api/v1/resources?is_available=true" \
  -H "Content-Type: application/json"
```

### Get Resource by ID

```bash
curl -X GET http://localhost:8080/api/v1/resources/1 \
  -H "Content-Type: application/json"
```

### Create Resource

```bash
curl -X POST http://localhost:8080/api/v1/resources \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Projector Room A",
    "type": "equipment",
    "isAvailable": true,
    "description": "High-resolution projector"
  }'
```

### Update Resource

```bash
curl -X PUT http://localhost:8080/api/v1/resources/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Projector Room B",
    "type": "equipment",
    "isAvailable": false,
    "description": "Under maintenance"
  }'
```

## Reservations API

### List Reservations

```bash
# Get all reservations
curl -X GET http://localhost:8080/api/v1/reservations \
  -H "Content-Type: application/json"

# Filter by date range
curl -X GET "http://localhost:8080/api/v1/reservations?fromDate=2023-01-01T00:00:00Z&toDate=2023-01-31T23:59:59Z" \
  -H "Content-Type: application/json"

# Filter by resource
curl -X GET "http://localhost:8080/api/v1/reservations?resourceId=1" \
  -H "Content-Type: application/json"
```

### Get Reservation by ID

```bash
curl -X GET http://localhost:8080/api/v1/reservations/1 \
  -H "Content-Type: application/json"
```

### Create Reservation

```bash
curl -X POST http://localhost:8080/api/v1/reservations \
  -H "Content-Type: application/json" \
  -d '{
    "resourceId": 1,
    "userId": 123,
    "startTime": "2023-06-01T10:00:00Z",
    "endTime": "2023-06-01T12:00:00Z",
    "purpose": "Team Meeting",
    "description": "Weekly planning session"
  }'
```

### Update Reservation Status

```bash
curl -X PUT http://localhost:8080/api/v1/reservations/1/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "confirmed"
  }'
```

### Delete Reservation

```bash
curl -X DELETE http://localhost:8080/api/v1/reservations/1 \
  -H "Content-Type: application/json"
```

## Classes API

### List Classes

```bash
curl -X GET http://localhost:8080/api/v1/classes \
  -H "Content-Type: application/json"
```

### Create Class

```bash
curl -X POST http://localhost:8080/api/v1/classes \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Room 101",
    "capacity": 30
  }'
```

## Lessons API

### List Lessons

```bash
curl -X GET http://localhost:8080/api/v1/lessons \
  -H "Content-Type: application/json"

# Filter by date range
curl -X GET "http://localhost:8080/api/v1/lessons?fromDate=2023-01-01T00:00:00Z&toDate=2023-01-31T23:59:59Z" \
  -H "Content-Type: application/json"
```

### Create Lesson

```bash
curl -X POST http://localhost:8080/api/v1/lessons \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Introduction to Programming",
    "duration": 90,
    "startTime": "2023-06-01T14:00:00Z",
    "endTime": "2023-06-01T15:30:00Z",
    "description": "Basic programming concepts"
  }'
```

## System Endpoints

### Health Check

```bash
curl -X GET http://localhost:8080/health
```

### Metrics

```bash
curl -X GET http://localhost:8080/metrics
```

## Authentication (Future Implementation)

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user@example.com",
    "password": "password123"
  }'
```

### Using JWT Token

```bash
# Replace TOKEN with your actual JWT token
curl -X GET http://localhost:8080/api/v1/buildings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN"
```

### Refresh Token

```bash
# Replace REFRESH_TOKEN with your actual refresh token
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "REFRESH_TOKEN"
  }'
```

## Advanced Usage

### Saving Response to File

```bash
curl -X GET http://localhost:8080/api/v1/buildings \
  -H "Content-Type: application/json" \
  > buildings.json
```

### Using jq to Parse JSON

[jq](https://stedolan.github.io/jq/) is a command-line JSON processor that's useful for working with API responses.

```bash
# Install jq first (apt-get install jq, brew install jq, etc.)

# Get all building names
curl -s http://localhost:8080/api/v1/buildings \
  -H "Content-Type: application/json" \
  | jq '.data[].name'

# Find buildings with a specific name pattern
curl -s http://localhost:8080/api/v1/buildings \
  -H "Content-Type: application/json" \
  | jq '.data[] | select(.name | contains("Science"))'
```

### Scripting Example

Here's a bash script that creates multiple buildings:

```bash
#!/bin/bash

API_URL="http://localhost:8080/api/v1"

# Array of buildings to create
buildings=(
  '{"name":"Physics Building","code":"PHY-A"}'
  '{"name":"Chemistry Building","code":"CHM-B"}'
  '{"name":"Biology Building","code":"BIO-C"}'
)

# Create each building
for building in "${buildings[@]}"; do
  echo "Creating building: $building"
  response=$(curl -s -X POST "$API_URL/buildings" \
    -H "Content-Type: application/json" \
    -d "$building")

  # Check for success
  if echo "$response" | grep -q "id"; then
    echo "Success! Created building with response: $response"
  else
    echo "Error creating building: $response"
  fi
done
```

## Error Handling Examples

### Handling 404 Not Found

```bash
response=$(curl -s -w "%{http_code}" -X GET http://localhost:8080/api/v1/buildings/999 \
  -H "Content-Type: application/json")

status=${response: -3}
body=${response:0:${#response}-3}

if [ "$status" -eq "404" ]; then
  echo "Building not found!"
else
  echo "Building found: $body"
fi
```

### Checking Server Health

```bash
health=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/health)

if [ "$health" -eq "200" ]; then
  echo "Server is healthy"
else
  echo "Server is unhealthy, returned status $health"
  # Send alert or take other action
fi
```

## Performance Testing

### Basic Load Testing with ApacheBench

```bash
# Install Apache Bench (ab) first
# Example: apt-get install apache2-utils

# Send 100 requests with 10 concurrent connections
ab -n 100 -c 10 http://localhost:8080/api/v1/buildings
```
