---
sidebar_position: 1
slug: /api-reference
---

# API Reference

Welcome to the SARC-NG API Reference. This section provides detailed documentation for all API endpoints, generated from our OpenAPI specification.

## Available Endpoints

The SARC-NG API is organized into the following resource categories:

- **Buildings** - Manage physical structures where classrooms and resources are located
- **Classes** - Manage rooms or spaces within buildings
- **Resources** - Manage equipment and other resources
- **Reservations** - Manage bookings of resources and spaces
- **Lessons** - Manage scheduled teaching sessions
- **Authentication** - User authentication and authorization

## Authentication

Most API endpoints require authentication. See the [Authentication](/api/authentication) section for details on how to authenticate your requests.

## Response Format

All API responses are returned in JSON format with consistent error handling. See the [Error Handling](/api/error-handling) section for details on error responses.

## Rate Limiting

API requests are subject to rate limiting to ensure fair usage. Rate limit information is provided in the response headers.

## Pagination

List endpoints support pagination. See the [Pagination](/api/pagination) section for details on how to paginate results.
