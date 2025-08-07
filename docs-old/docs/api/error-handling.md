---
sidebar_position: 4
---

# Error Handling

SARC-NG API implements consistent error handling and reporting to help clients understand and resolve issues.

## Error Response Format

All API errors follow a standard JSON format:

```json
{
  "code": 400,
  "error": "Invalid input",
  "message": "The provided data is invalid",
  "details": [
    {
      "field": "name",
      "message": "Name is required"
    }
  ]
}
```

Where:
- `code`: HTTP status code
- `error`: Error type identifier
- `message`: Human-readable error message
- `details`: (Optional) Array of specific validation errors or additional information

## HTTP Status Codes

The API uses standard HTTP status codes to indicate the result of requests:

### Success Codes

- **200 OK**: Successful GET or PUT request
- **201 Created**: Successful POST request that created a new resource
- **204 No Content**: Successful DELETE request

### Client Error Codes

- **400 Bad Request**: Invalid request format or parameters
- **401 Unauthorized**: Authentication required
- **403 Forbidden**: Authenticated but insufficient permissions
- **404 Not Found**: Resource not found
- **409 Conflict**: Resource conflict (e.g., reservation overlap)
- **422 Unprocessable Entity**: Request syntax correct but semantically invalid

### Server Error Codes

- **500 Internal Server Error**: Unexpected server error
- **503 Service Unavailable**: Temporary server unavailability

## Common Error Types

### Validation Errors

```json
{
  "code": 400,
  "error": "ValidationError",
  "message": "The request contains invalid data",
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

### Resource Not Found

```json
{
  "code": 404,
  "error": "ResourceNotFound",
  "message": "The requested resource was not found",
  "details": {
    "type": "building",
    "id": 123
  }
}
```

### Authentication Errors

```json
{
  "code": 401,
  "error": "Unauthorized",
  "message": "Authentication required"
}
```

### Permission Errors

```json
{
  "code": 403,
  "error": "Forbidden",
  "message": "Insufficient permissions to perform this action"
}
```

### Reservation Conflicts

```json
{
  "code": 409,
  "error": "ResourceConflict",
  "message": "Resource is not available for the requested time",
  "details": {
    "resourceId": 5,
    "conflictingReservationId": 42,
    "requestedStartTime": "2024-06-15T14:00:00Z",
    "requestedEndTime": "2024-06-15T16:00:00Z"
  }
}
```

### Database Errors

```json
{
  "code": 500,
  "error": "DatabaseError",
  "message": "An error occurred while accessing the database"
}
```

## Error Handling Best Practices

When working with the SARC-NG API, follow these best practices for handling errors:

### 1. Check HTTP Status Codes

Always check the HTTP status code first to determine the category of error.

```javascript
if (response.status >= 400) {
  // Handle error based on status code
}
```

### 2. Parse Error Messages

Extract meaningful error information from the response body.

```javascript
async function handleErrorResponse(response) {
  const errorData = await response.json();
  console.error(`Error ${errorData.code}: ${errorData.message}`);

  // Handle specific error types
  if (errorData.code === 409) {
    // Handle conflict
  }
}
```

### 3. Handle Validation Errors

For validation errors, display specific messages for each field.

```javascript
function displayValidationErrors(errorDetails) {
  errorDetails.forEach(error => {
    const field = document.getElementById(`field-${error.field}`);
    if (field) {
      field.classList.add('error');
      field.setAttribute('data-error', error.message);
    }
  });
}
```

### 4. Implement Retry Logic

For certain errors (e.g., 503 Service Unavailable), implement retry logic with exponential backoff.

```javascript
async function fetchWithRetry(url, options, maxRetries = 3) {
  let retries = 0;

  while (retries < maxRetries) {
    try {
      const response = await fetch(url, options);
      if (response.status !== 503) return response;

      retries++;
      const delay = Math.pow(2, retries) * 1000;
      await new Promise(resolve => setTimeout(resolve, delay));
    } catch (error) {
      if (retries === maxRetries - 1) throw error;
      retries++;
    }
  }
}
```

### 5. Log Errors

Log errors for debugging and monitoring purposes.

```javascript
function logApiError(error, request) {
  console.error('API Error', {
    url: request.url,
    method: request.method,
    status: error.code,
    message: error.message,
    timestamp: new Date().toISOString(),
  });

  // Could also send to error tracking service
}
```

## Common Error Scenarios

| Scenario | Status | Error Type | Solution |
|----------|--------|------------|----------|
| Missing required field | 400 | ValidationError | Check required fields before submitting |
| Invalid format | 400 | ValidationError | Validate data formats client-side |
| Resource not found | 404 | ResourceNotFound | Verify IDs are correct and resources exist |
| Duplicate resource | 409 | ResourceConflict | Check for existing resources before creating |
| Reservation conflict | 409 | ResourceConflict | Check availability before booking |
| Invalid credentials | 401 | Unauthorized | Prompt user to re-authenticate |
| Expired token | 401 | Unauthorized | Refresh token or re-authenticate |
| Insufficient permissions | 403 | Forbidden | Request appropriate permissions |
| Server error | 500 | InternalError | Retry or contact system administrator |

## Error Codes Reference

| Error Code | Error Type | Description |
|------------|------------|-------------|
| `INVALID_INPUT` | ValidationError | Request contains invalid data |
| `RESOURCE_NOT_FOUND` | ResourceNotFound | Requested resource does not exist |
| `RESOURCE_CONFLICT` | ResourceConflict | Resource state conflicts with request |
| `UNAUTHORIZED` | Unauthorized | Authentication required |
| `FORBIDDEN` | Forbidden | Insufficient permissions |
| `INTERNAL_ERROR` | InternalError | Unexpected server error |
| `DATABASE_ERROR` | DatabaseError | Error communicating with database |
| `SERVICE_UNAVAILABLE` | ServiceUnavailable | Service temporarily unavailable |
