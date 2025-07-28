# API Documentation

REST API reference and integration guide for SARC-NG.

## API Overview

SARC-NG provides a RESTful API for managing resource reservations and scheduling. The API follows standard HTTP conventions and returns JSON-formatted responses.

**Base URL:** `http://localhost:8080/api/v1`

**Version:** 1.0.0

### Core Features

- **Resource Management**: Manage bookable resources (rooms, equipment)
- **Building Management**: Organize physical structures and locations
- **Class Management**: Handle classroom and space definitions
- **Lesson Scheduling**: Schedule and manage educational sessions
- **Reservation System**: Book resources with conflict detection

## Authentication

Currently under development. Authentication will be implemented using JWT tokens.

**Planned Authentication Flow:**
```http
POST /api/v1/auth/login
Authorization: Bearer <jwt_token>
```

**Future Headers:**
```http
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

## API Endpoints

### Buildings

Manage physical buildings and locations.

#### List All Buildings

```http
GET /api/v1/buildings
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Main Campus Building",
    "code": "MCB-A",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
]
```

#### Get Building by ID

```http
GET /api/v1/buildings/{id}
```

**Parameters:**
- `id` (path, integer): Building ID

**Response:**
```json
{
  "id": 1,
  "name": "Main Campus Building",
  "code": "MCB-A",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### Create Building

```http
POST /api/v1/buildings
```

**Request Body:**
```json
{
  "name": "Science Building",
  "code": "SCI-B"
}
```

**Response:**
```json
{
  "id": 2,
  "name": "Science Building",
  "code": "SCI-B",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### Update Building

```http
PUT /api/v1/buildings/{id}
```

**Request Body:**
```json
{
  "id": 1,
  "name": "Main Academic Building",
  "code": "MAB-A"
}
```

#### Delete Building

```http
DELETE /api/v1/buildings/{id}
```

**Response:** `204 No Content`

### Classes

Manage classrooms and physical spaces.

#### List All Classes

```http
GET /api/v1/classes
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Room 101",
    "capacity": 30,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
]
```

#### Create Class

```http
POST /api/v1/classes
```

**Request Body:**
```json
{
  "name": "Conference Room A",
  "capacity": 50
}
```

**Response:**
```json
{
  "id": 2,
  "name": "Conference Room A",
  "capacity": 50,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### Resources

Manage bookable resources like rooms and equipment.

#### List All Resources

```http
GET /api/v1/resources
```

**Query Parameters:**
- `type` (optional): Filter resources by type

**Response:**
```json
[
  {
    "id": 1,
    "name": "Projector Room A",
    "type": "equipment",
    "isAvailable": true,
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
]
```

#### Get Resource by ID

```http
GET /api/v1/resources/{id}
```

**Response:**
```json
{
  "id": 1,
  "name": "Projector Room A",
  "type": "equipment",
  "isAvailable": true,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### Create Resource

```http
POST /api/v1/resources
```

**Request Body:**
```json
{
  "name": "Laptop Set 1",
  "type": "equipment",
  "isAvailable": true
}
```

**Response:**
```json
{
  "id": 2,
  "name": "Laptop Set 1",
  "type": "equipment",
  "isAvailable": true,
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### Update Resource

```http
PUT /api/v1/resources/{id}
```

**Request Body:**
```json
{
  "id": 1,
  "name": "Projector Room A - Updated",
  "type": "equipment",
  "isAvailable": false
}
```

#### Delete Resource

```http
DELETE /api/v1/resources/{id}
```

**Response:** `204 No Content`

### Lessons

Schedule and manage educational sessions.

#### List All Lessons

```http
GET /api/v1/lessons
```

**Query Parameters:**
- `classId` (optional): Filter by class ID
- `fromDate` (optional): Start date for filtering (YYYY-MM-DD)
- `toDate` (optional): End date for filtering (YYYY-MM-DD)

**Response:**
```json
[
  {
    "id": 1,
    "title": "Introduction to Programming",
    "description": "Basic programming concepts",
    "duration": 90,
    "startTime": "2024-01-01T14:00:00Z",
    "endTime": "2024-01-01T15:30:00Z",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
]
```

#### Create Lesson

```http
POST /api/v1/lessons
```

**Request Body:**
```json
{
  "title": "Advanced Database Design",
  "description": "Database normalization and optimization",
  "duration": 120,
  "startTime": "2024-01-02T10:00:00Z",
  "endTime": "2024-01-02T12:00:00Z"
}
```

**Response:**
```json
{
  "id": 2,
  "title": "Advanced Database Design",
  "description": "Database normalization and optimization", 
  "duration": 120,
  "startTime": "2024-01-02T10:00:00Z",
  "endTime": "2024-01-02T12:00:00Z",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### Reservations

Book and manage resource reservations.

#### List All Reservations

```http
GET /api/v1/reservations
```

**Query Parameters:**
- `fromDate` (optional): Start date for filtering (YYYY-MM-DD)
- `toDate` (optional): End date for filtering (YYYY-MM-DD)
- `resourceId` (optional): Filter by resource ID

**Response:**
```json
[
  {
    "id": 1,
    "resourceId": 1,
    "userId": 123,
    "startTime": "2024-01-01T14:00:00Z",
    "endTime": "2024-01-01T16:00:00Z",
    "purpose": "Team Meeting",
    "status": "confirmed",
    "description": "Weekly team standup",
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z"
  }
]
```

#### Create Reservation

```http
POST /api/v1/reservations
```

**Request Body:**
```json
{
  "resourceId": 1,
  "userId": 123,
  "startTime": "2024-01-03T10:00:00Z",
  "endTime": "2024-01-03T12:00:00Z",
  "purpose": "Training Session",
  "description": "New employee onboarding"
}
```

**Response:**
```json
{
  "id": 2,
  "resourceId": 1,
  "userId": 123,
  "startTime": "2024-01-03T10:00:00Z",
  "endTime": "2024-01-03T12:00:00Z",
  "purpose": "Training Session",
  "status": "pending",
  "description": "New employee onboarding",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

#### Update Reservation

```http
PUT /api/v1/reservations/{id}
```

**Request Body:**
```json
{
  "id": 1,
  "resourceId": 1,
  "userId": 123,
  "startTime": "2024-01-03T09:00:00Z",
  "endTime": "2024-01-03T11:00:00Z",
  "purpose": "Training Session - Updated",
  "status": "confirmed",
  "description": "Updated training session details"
}
```

#### Delete Reservation

```http
DELETE /api/v1/reservations/{id}
```

**Response:** `204 No Content`

## System Endpoints

### Health Check

```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "sarc-ng"
}
```

### Metrics

```http
GET /metrics
```

**Response:** Prometheus-formatted metrics

### API Documentation

```http
GET /swagger/index.html
```

Interactive Swagger UI for API exploration and testing.

## Response Formats

### Success Response

Standard successful response format:

```json
{
  "data": { },
  "message": "Operation completed successfully"
}
```

### Error Response

Standard error response format:

```json
{
  "code": 400,
  "error": "Invalid input",
  "message": "The provided data is invalid"
}
```

### HTTP Status Codes

- **200 OK**: Successful GET, PUT requests
- **201 Created**: Successful POST requests
- **204 No Content**: Successful DELETE requests
- **400 Bad Request**: Invalid request data
- **401 Unauthorized**: Authentication required
- **403 Forbidden**: Access denied
- **404 Not Found**: Resource not found
- **409 Conflict**: Resource conflict (reservation overlap)
- **500 Internal Server Error**: Server error

## Data Types and Formats

### Date/Time Format

All date/time fields use ISO 8601 format with UTC timezone:

```
"2024-01-01T14:00:00Z"
```

### Validation Rules

**Building:**
- `name`: Required, string
- `code`: Required, string, unique

**Class:**
- `name`: Required, string
- `capacity`: Required, integer, minimum 1

**Resource:**
- `name`: Required, string
- `type`: Required, string
- `isAvailable`: Boolean, defaults to true

**Lesson:**
- `title`: Required, string
- `duration`: Required, integer (minutes)
- `startTime`: Required, ISO 8601 datetime
- `endTime`: Required, ISO 8601 datetime

**Reservation:**
- `resourceId`: Required, integer
- `userId`: Required, integer
- `startTime`: Required, ISO 8601 datetime
- `endTime`: Required, ISO 8601 datetime
- `purpose`: Required, string

## Client Integration

### JavaScript/TypeScript Example

```javascript
// API Client Configuration
const API_BASE_URL = 'http://localhost:8080/api/v1';

class SarcApiClient {
  constructor(baseUrl = API_BASE_URL) {
    this.baseUrl = baseUrl;
    this.headers = {
      'Content-Type': 'application/json',
    };
  }

  // Buildings
  async getBuildings() {
    const response = await fetch(`${this.baseUrl}/buildings`, {
      headers: this.headers
    });
    return response.json();
  }

  async createBuilding(building) {
    const response = await fetch(`${this.baseUrl}/buildings`, {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify(building)
    });
    return response.json();
  }

  // Resources
  async getResources(type = null) {
    const url = new URL(`${this.baseUrl}/resources`);
    if (type) url.searchParams.append('type', type);
    
    const response = await fetch(url, {
      headers: this.headers
    });
    return response.json();
  }

  async createReservation(reservation) {
    const response = await fetch(`${this.baseUrl}/reservations`, {
      method: 'POST',
      headers: this.headers,
      body: JSON.stringify(reservation)
    });
    
    if (response.status === 409) {
      throw new Error('Reservation conflict: Resource is not available for the requested time');
    }
    
    return response.json();
  }
}

// Usage Example
const client = new SarcApiClient();

// Get all buildings
const buildings = await client.getBuildings();
console.log('Buildings:', buildings);

// Create a new reservation
try {
  const reservation = await client.createReservation({
    resourceId: 1,
    userId: 123,
    startTime: '2024-01-03T10:00:00Z',
    endTime: '2024-01-03T12:00:00Z',
    purpose: 'Meeting',
    description: 'Team planning session'
  });
  console.log('Reservation created:', reservation);
} catch (error) {
  console.error('Failed to create reservation:', error.message);
}
```

### Python Example

```python
import requests
from datetime import datetime
from typing import Optional, Dict, List

class SarcApiClient:
    def __init__(self, base_url: str = "http://localhost:8080/api/v1"):
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json'
        })

    def get_buildings(self) -> List[Dict]:
        """Get all buildings"""
        response = self.session.get(f"{self.base_url}/buildings")
        response.raise_for_status()
        return response.json()

    def create_building(self, name: str, code: str) -> Dict:
        """Create a new building"""
        data = {"name": name, "code": code}
        response = self.session.post(f"{self.base_url}/buildings", json=data)
        response.raise_for_status()
        return response.json()

    def get_resources(self, resource_type: Optional[str] = None) -> List[Dict]:
        """Get all resources, optionally filtered by type"""
        params = {}
        if resource_type:
            params['type'] = resource_type
            
        response = self.session.get(f"{self.base_url}/resources", params=params)
        response.raise_for_status()
        return response.json()

    def create_reservation(self, resource_id: int, user_id: int, 
                         start_time: str, end_time: str, 
                         purpose: str, description: str = "") -> Dict:
        """Create a new reservation"""
        data = {
            "resourceId": resource_id,
            "userId": user_id,
            "startTime": start_time,
            "endTime": end_time,
            "purpose": purpose,
            "description": description
        }
        
        response = self.session.post(f"{self.base_url}/reservations", json=data)
        
        if response.status_code == 409:
            raise ValueError("Reservation conflict: Resource is not available for the requested time")
        
        response.raise_for_status()
        return response.json()

# Usage Example
client = SarcApiClient()

# Get all resources
resources = client.get_resources()
print(f"Found {len(resources)} resources")

# Create a reservation
try:
    reservation = client.create_reservation(
        resource_id=1,
        user_id=123,
        start_time="2024-01-03T10:00:00Z",
        end_time="2024-01-03T12:00:00Z",
        purpose="Team Meeting",
        description="Weekly planning session"
    )
    print(f"Reservation created with ID: {reservation['id']}")
except ValueError as e:
    print(f"Reservation failed: {e}")
```

### cURL Examples

**Create a Building:**
```bash
curl -X POST http://localhost:8080/api/v1/buildings \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Engineering Building",
    "code": "ENG-A"
  }'
```

**Get Resources by Type:**
```bash
curl "http://localhost:8080/api/v1/resources?type=equipment" \
  -H "Content-Type: application/json"
```

**Create a Reservation:**
```bash
curl -X POST http://localhost:8080/api/v1/reservations \
  -H "Content-Type: application/json" \
  -d '{
    "resourceId": 1,
    "userId": 123,
    "startTime": "2024-01-03T10:00:00Z",
    "endTime": "2024-01-03T12:00:00Z",
    "purpose": "Workshop",
    "description": "Technical training workshop"
  }'
```

## Error Handling

### Common Error Scenarios

**Validation Errors (400):**
```json
{
  "code": 400,
  "error": "Validation failed",
  "message": "Name field is required"
}
```

**Resource Not Found (404):**
```json
{
  "code": 404,
  "error": "Resource not found",
  "message": "No building exists with the specified ID"
}
```

**Reservation Conflict (409):**
```json
{
  "code": 409,
  "error": "Reservation conflict",
  "message": "Resource is not available for the requested time"
}
```

### Error Handling Best Practices

1. **Check HTTP Status Codes**: Always verify the response status
2. **Parse Error Messages**: Extract meaningful error information
3. **Handle Conflicts**: Implement retry logic for reservation conflicts
4. **Validate Input**: Validate data before sending requests
5. **Implement Timeouts**: Set appropriate request timeouts

## Rate Limiting

Currently not implemented. Future versions will include:

- Rate limiting by IP address
- Authenticated user quotas
- Burst protection mechanisms

## Pagination

For endpoints returning large datasets, pagination will be implemented:

**Future Pagination Format:**
```http
GET /api/v1/reservations?page=1&limit=50
```

**Response with Pagination:**
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 200,
    "pages": 4
  }
}
```

## WebSocket Support

Planned for real-time updates:

**Future WebSocket Endpoints:**
- `/ws/reservations` - Real-time reservation updates
- `/ws/resources` - Resource availability changes

## Testing the API

### Using Swagger UI

1. Start the application
2. Navigate to `http://localhost:8080/swagger/index.html`
3. Explore endpoints and test requests interactively

### Integration Testing

```bash
# Start the application with Docker
docker compose up -d

# Run integration tests
go test -tags=integration ./test/integration/...

# Cleanup
docker compose down -v --remove-orphans
```

### API Health Verification

```bash
# Check API health
curl http://localhost:8080/health

# Verify all endpoints are accessible
curl http://localhost:8080/api/v1/buildings
curl http://localhost:8080/api/v1/resources
curl http://localhost:8080/api/v1/reservations
curl http://localhost:8080/api/v1/classes
curl http://localhost:8080/api/v1/lessons
``` 