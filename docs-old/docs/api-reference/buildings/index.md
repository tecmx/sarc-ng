---
sidebar_position: 1
---

# Buildings API

API endpoints for managing building resources.

## Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/buildings` | List all buildings |
| GET | `/api/v1/buildings/:id` | Get building by ID |
| POST | `/api/v1/buildings` | Create a new building |
| PUT | `/api/v1/buildings/:id` | Update a building |
| DELETE | `/api/v1/buildings/:id` | Delete a building |

## List Buildings

Retrieves a list of all buildings with pagination support.

```http
GET /api/v1/buildings?page=1&limit=10&sort=name&direction=asc
```

### Query Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number |
| `limit` | integer | 50 | Results per page (max 100) |
| `sort` | string | `id` | Sort field (`id`, `name`, `code`, `created_at`, `updated_at`) |
| `direction` | string | `asc` | Sort direction (`asc`, `desc`) |

### Response

```json
{
  "data": [
    {
      "id": 1,
      "name": "Main Campus Building",
      "code": "MCB-A",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    },
    {
      "id": 2,
      "name": "Science Building",
      "code": "SCI-B",
      "createdAt": "2024-01-01T00:00:00Z",
      "updatedAt": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 2,
    "pages": 1
  }
}
```

## Get Building by ID

Retrieves a specific building by its ID.

```http
GET /api/v1/buildings/:id
```

### URL Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Building ID |

### Response

```json
{
  "data": {
    "id": 1,
    "name": "Main Campus Building",
    "code": "MCB-A",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
}
```

### Error Responses

**Building not found (404):**

```json
{
  "code": 404,
  "error": "ResourceNotFound",
  "message": "Building not found",
  "details": {
    "id": 999
  }
}
```

## Create Building

Creates a new building.

```http
POST /api/v1/buildings
Content-Type: application/json

{
  "name": "Engineering Building",
  "code": "ENG-A"
}
```

### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Building name |
| `code` | string | Yes | Unique building code |

### Response

```json
{
  "data": {
    "id": 3,
    "name": "Engineering Building",
    "code": "ENG-A",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "message": "Building created successfully"
}
```

### Error Responses

**Validation error (400):**

```json
{
  "code": 400,
  "error": "ValidationError",
  "message": "Invalid building data",
  "details": [
    {
      "field": "name",
      "message": "Name is required"
    },
    {
      "field": "code",
      "message": "Code must be unique"
    }
  ]
}
```

## Update Building

Updates an existing building.

```http
PUT /api/v1/buildings/:id
Content-Type: application/json

{
  "name": "Engineering Building Updated",
  "code": "ENG-B"
}
```

### URL Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Building ID |

### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Building name |
| `code` | string | Yes | Unique building code |

### Response

```json
{
  "data": {
    "id": 3,
    "name": "Engineering Building Updated",
    "code": "ENG-B",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  },
  "message": "Building updated successfully"
}
```

### Error Responses

**Building not found (404):**

```json
{
  "code": 404,
  "error": "ResourceNotFound",
  "message": "Building not found",
  "details": {
    "id": 999
  }
}
```

**Validation error (400):**

```json
{
  "code": 400,
  "error": "ValidationError",
  "message": "Invalid building data",
  "details": [
    {
      "field": "code",
      "message": "Code must be unique"
    }
  ]
}
```

## Delete Building

Deletes a building.

```http
DELETE /api/v1/buildings/:id
```

### URL Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Building ID |

### Response

```json
{
  "message": "Building deleted successfully"
}
```

### Error Responses

**Building not found (404):**

```json
{
  "code": 404,
  "error": "ResourceNotFound",
  "message": "Building not found",
  "details": {
    "id": 999
  }
}
```

**Conflict (409):**

```json
{
  "code": 409,
  "error": "ResourceConflict",
  "message": "Cannot delete building with associated resources",
  "details": {
    "associatedResourcesCount": 5
  }
}
```

## Status Codes

| Status Code | Description |
|-------------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Building created successfully |
| 204 | No Content - Building deleted successfully |
| 400 | Bad Request - Invalid request data |
| 404 | Not Found - Building not found |
| 409 | Conflict - Resource conflict |
| 500 | Internal Server Error - Server error |

## Examples

### cURL Examples

```bash
# List buildings
curl -X GET http://localhost:8080/api/v1/buildings

# Get building by ID
curl -X GET http://localhost:8080/api/v1/buildings/1

# Create building
curl -X POST http://localhost:8080/api/v1/buildings \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Engineering Building",
    "code": "ENG-A"
  }'

# Update building
curl -X PUT http://localhost:8080/api/v1/buildings/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Engineering Building Updated",
    "code": "ENG-B"
  }'

# Delete building
curl -X DELETE http://localhost:8080/api/v1/buildings/1
```

### JavaScript Example

```javascript
// Create a building
fetch('http://localhost:8080/api/v1/buildings', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    name: 'Engineering Building',
    code: 'ENG-A'
  }),
})
.then(response => response.json())
.then(data => console.log('Building created:', data))
.catch(error => console.error('Error:', error));
```
