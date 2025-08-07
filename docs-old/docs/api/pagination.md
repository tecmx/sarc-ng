---
sidebar_position: 3
---

# Pagination

SARC-NG API implements standardized pagination for list endpoints to efficiently handle large datasets.

## Overview

For endpoints that can return large collections of items, the API supports pagination through query parameters and includes pagination metadata in responses.

## Pagination Parameters

All list endpoints support the following query parameters:

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | 1 | Page number to retrieve |
| `limit` | integer | 50 | Number of items per page (max 100) |
| `sort` | string | `id` | Field to sort by |
| `direction` | string | `asc` | Sort direction (`asc` or `desc`) |

## Example Request

```http
GET /api/v1/buildings?page=2&limit=10&sort=name&direction=asc
```

This request retrieves:
- The second page of buildings
- 10 buildings per page
- Sorted by name in ascending order

## Pagination Response Format

When pagination is applied, responses include a `pagination` object with metadata:

```json
{
  "data": [
    // Array of items
  ],
  "pagination": {
    "page": 2,
    "limit": 10,
    "total": 35,
    "pages": 4
  }
}
```

Where:
- `page`: Current page number
- `limit`: Number of items per page
- `total`: Total number of items across all pages
- `pages`: Total number of pages

## Navigating Pages

To navigate through paginated results, adjust the `page` parameter:

```http
# First page
GET /api/v1/buildings?page=1&limit=10

# Next page
GET /api/v1/buildings?page=2&limit=10

# Last page (if total pages is known)
GET /api/v1/buildings?page=4&limit=10
```

## Sorting Options

The `sort` parameter accepts different fields depending on the endpoint:

### Buildings
- `id` (default)
- `name`
- `code`
- `created_at`
- `updated_at`

### Resources
- `id` (default)
- `name`
- `type`
- `is_available`
- `created_at`
- `updated_at`

### Reservations
- `id` (default)
- `resource_id`
- `user_id`
- `start_time`
- `end_time`
- `status`
- `created_at`
- `updated_at`

## Filtering with Pagination

Pagination works together with filtering parameters. For example:

```http
GET /api/v1/resources?type=equipment&is_available=true&page=1&limit=20
```

This retrieves:
- The first page of resources
- Only those of type "equipment" that are available
- 20 items per page

## Error Handling

If an invalid pagination parameter is provided, the API will respond with a 400 Bad Request:

```json
{
  "code": 400,
  "error": "Invalid pagination parameters",
  "message": "Page must be a positive integer"
}
```

Common pagination errors:
- Page number less than 1
- Limit greater than 100
- Invalid sort field
- Invalid sort direction

## Implementation Notes

1. **Performance Considerations**:
   - Use database indexes on sorted fields
   - Use SQL LIMIT and OFFSET clauses efficiently
   - Avoid extremely large offsets for better performance

2. **Consistency**:
   - All list endpoints follow the same pagination pattern
   - Responses always include pagination metadata, even for small result sets

3. **Cursor-Based Pagination**:
   - For very large datasets, consider using the `since_id` parameter instead of page-based pagination
   - Example: `GET /api/v1/buildings?since_id=123&limit=50`

## Examples

### Retrieving Buildings with Pagination

**Request:**
```http
GET /api/v1/buildings?page=1&limit=2
```

**Response:**
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
    "limit": 2,
    "total": 5,
    "pages": 3
  }
}
```
