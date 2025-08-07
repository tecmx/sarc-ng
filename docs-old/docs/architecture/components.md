---
sidebar_position: 2
---

# Components

This document details the core components of SARC-NG and their responsibilities.

## Domain Components

SARC-NG is organized around the following domain concepts:

### Buildings

Buildings represent physical structures where classrooms and resources are located.

```go
// Building entity
type Building struct {
    ID        uint
    Name      string
    Code      string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**Key Files:**
- `internal/domain/building/entity.go` - Entity definition
- `internal/domain/building/repository.go` - Repository interface
- `internal/domain/building/usecase.go` - Business logic interfaces

### Classes

Classes represent rooms or spaces within buildings.

```go
// Class entity
type Class struct {
    ID        uint
    Name      string
    Capacity  int
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

**Key Files:**
- `internal/domain/class/entity.go` - Entity definition
- `internal/domain/class/repository.go` - Repository interface
- `internal/domain/class/usecase.go` - Business logic interfaces

### Resources

Resources are bookable items such as rooms, equipment, or facilities.

```go
// Resource entity
type Resource struct {
    ID          uint
    Name        string
    Type        string
    IsAvailable bool
    Description string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

**Key Files:**
- `internal/domain/resource/entity.go` - Entity definition
- `internal/domain/resource/repository.go` - Repository interface
- `internal/domain/resource/usecase.go` - Business logic interfaces

### Lessons

Lessons represent scheduled teaching or training sessions.

```go
// Lesson entity
type Lesson struct {
    ID          uint
    Title       string
    Description string
    Duration    int
    StartTime   time.Time
    EndTime     time.Time
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

**Key Files:**
- `internal/domain/lesson/entity.go` - Entity definition
- `internal/domain/lesson/repository.go` - Repository interface
- `internal/domain/lesson/usecase.go` - Business logic interfaces

### Reservations

Reservations track bookings for resources, including who requested them and for what purpose.

```go
// Reservation entity
type Reservation struct {
    ID          uint
    ResourceID  uint
    UserID      uint
    StartTime   time.Time
    EndTime     time.Time
    Purpose     string
    Status      string
    Description string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

**Key Files:**
- `internal/domain/reservation/entity.go` - Entity definition
- `internal/domain/reservation/repository.go` - Repository interface
- `internal/domain/reservation/usecase.go` - Business logic interfaces

## Services

Services implement the business logic and use cases for each domain entity.

### Building Service

Manages building-related operations.

**Responsibilities:**
- Create, update, and delete buildings
- Retrieve building information
- Validate building data

**Key Files:**
- `internal/service/building/service.go`

### Class Service

Manages classroom and space operations.

**Responsibilities:**
- Create, update, and delete classes
- Retrieve class information
- Track class capacity

**Key Files:**
- `internal/service/class/service.go`

### Resource Service

Manages bookable resources.

**Responsibilities:**
- Create, update, and delete resources
- Track resource availability
- Categorize resources by type

**Key Files:**
- `internal/service/resource/service.go`

### Lesson Service

Manages scheduled lessons.

**Responsibilities:**
- Create, update, and delete lessons
- Schedule lessons at specific times
- Manage lesson durations

**Key Files:**
- `internal/service/lesson/service.go`

### Reservation Service

Handles the booking and management of resources.

**Responsibilities:**
- Create, update, and delete reservations
- Check for conflicts between reservations
- Validate reservation time periods
- Enforce business rules for bookings

**Key Files:**
- `internal/service/reservation/service.go`

## Transport Layer

The transport layer handles HTTP communications and API endpoints.

### REST Handlers

Each domain entity has a corresponding REST handler.

**Building Handler:**
- `internal/transport/rest/building/handler.go`
- Endpoints for CRUD operations on buildings

**Class Handler:**
- `internal/transport/rest/class/handler.go`
- Endpoints for CRUD operations on classes

**Resource Handler:**
- `internal/transport/rest/resource/handler.go`
- Endpoints for CRUD operations on resources

**Lesson Handler:**
- `internal/transport/rest/lesson/handler.go`
- Endpoints for CRUD operations on lessons

**Reservation Handler:**
- `internal/transport/rest/reservation/handler.go`
- Endpoints for creating, updating, and managing reservations
- Endpoints for checking resource availability

### Router

The router configures all API endpoints and middleware.

**Key Files:**
- `internal/transport/rest/router.go` - Main router configuration

### Middleware

Middleware components provide cross-cutting concerns.

**Available Middleware:**
- Logger - Request logging
- CORS - Cross-origin resource sharing
- Authentication - JWT token validation
- Metrics - Prometheus metrics collection

**Key Files:**
- `pkg/rest/middleware/*.go` - Middleware implementations

## Infrastructure Layer

The infrastructure layer connects the application to external systems.

### Database Adapters

Database adapters provide the implementation for repositories.

**GORM Adapters:**
- `internal/adapter/gorm/building/adapter.go`
- `internal/adapter/gorm/class/adapter.go`
- `internal/adapter/gorm/resource/adapter.go`
- `internal/adapter/gorm/lesson/adapter.go`
- `internal/adapter/gorm/reservation/adapter.go`

### Database Configuration

Database connection and configuration management.

**Key Files:**
- `internal/adapter/db/config.go` - Database configuration
- `internal/adapter/db/connection.go` - Connection management

## Application Entry Points

SARC-NG has multiple application entry points:

### HTTP Server

The main HTTP API server.

**Key Files:**
- `cmd/server/main.go` - Server entry point
- `cmd/server/wire.go` - Dependency injection

### CLI Tool

Command-line interface for administrative tasks.

**Key Files:**
- `cmd/cli/main.go` - CLI entry point
- `cmd/cli/root.go` - Root command configuration

### Lambda Function

AWS Lambda entry point for serverless deployment.

**Key Files:**
- `cmd/lambda/main.go` - Lambda entry point
- `cmd/lambda/wire.go` - Dependency injection

## Configuration

Configuration management for the application.

**Key Files:**
- `internal/config/config.go` - Configuration structures
- `internal/config/loader.go` - Configuration loading logic
- `configs/*.yaml` - Configuration files

## Common Components

### Error Handling

Centralized error handling and mapping.

**Key Files:**
- `internal/domain/common/errors.go` - Domain error definitions
- `internal/transport/common/error_mapper.go` - Error mapping to HTTP responses

### Response Formatting

Consistent response formatting for API endpoints.

**Key Files:**
- `internal/transport/common/response.go` - Response formatting utilities

### Utility Functions

Common utility functions used across the application.

**Key Files:**
- `internal/adapter/gorm/common/utils.go` - GORM utilities

## External Client Libraries

Client libraries for consuming the SARC-NG API.

**Key Files:**
- `pkg/rest/client/client.go` - Base API client
- `pkg/rest/client/buildings.go` - Buildings API client
- `pkg/rest/client/classes.go` - Classes API client