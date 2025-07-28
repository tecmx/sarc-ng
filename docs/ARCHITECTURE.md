# Architecture

SARC-NG system architecture and design overview.

## System Purpose

SARC-NG (Schedule and Resource Control - Next Generation) is a resource reservation and management API for educational and organizational environments. The system manages buildings, classes, lessons, resources, and reservations with support for scheduling and availability checking.

## Architecture Overview

SARC-NG follows Clean Architecture principles with clear separation of concerns across four main layers:

```
┌─────────────────────────────────────────────────────────────┐
│                    Transport Layer                          │
│  (HTTP Handlers, REST API, CLI Commands, Lambda)           │
├─────────────────────────────────────────────────────────────┤
│                    Service Layer                            │
│       (Application Services, Business Use Cases)           │
├─────────────────────────────────────────────────────────────┤
│                    Domain Layer                             │
│        (Entities, Business Logic, Interfaces)              │
├─────────────────────────────────────────────────────────────┤
│                    Adapter Layer                            │
│    (Database, External APIs, Infrastructure)               │
└─────────────────────────────────────────────────────────────┘
```

### Key Principles

- **Clean Architecture**: Clear boundaries between layers with dependency inversion
- **Domain-Driven Design**: Business logic encapsulated in domain entities
- **Dependency Injection**: Google Wire for compile-time dependency injection
- **Interface Segregation**: Small, focused interfaces for better testability

## Core Components

### Domain Entities

The system manages five core business entities:

**Building**: Physical structures and locations
- Properties: ID, Name, Code, timestamps
- Purpose: Organize physical spaces for resource allocation

**Class**: Classrooms or physical spaces
- Properties: ID, Name, Capacity, timestamps  
- Purpose: Define available spaces for lessons and activities

**Resource**: Bookable resources (rooms, equipment, etc.)
- Properties: ID, Name, Type, Description, IsAvailable, Location, timestamps
- Purpose: Represent reservable items and spaces

**Lesson**: Individual teaching or training sessions
- Properties: ID, Title, Duration, Description, StartTime, EndTime, timestamps
- Purpose: Scheduled educational activities

**Reservation**: Booking records for resources
- Properties: ID, ResourceID, UserID, StartTime, EndTime, Purpose, Status, Description, timestamps
- Purpose: Track resource allocations and prevent conflicts

### Layer Details

#### 1. Transport Layer (`internal/transport/`)

Handles external communication and protocol adaptation:

**REST API** (`rest/`):
- Gin-based HTTP handlers for each entity
- DTO mapping for request/response transformation
- Swagger documentation generation
- Middleware for CORS, logging, metrics, recovery

**CLI Interface** (`cmd/cli/`):
- Cobra-based command structure
- Entity-specific commands with formatters
- API client integration for remote operations

**Lambda Handler** (`cmd/lambda/`):
- AWS Lambda integration via Gin adapter
- API Gateway event processing
- Cold start optimization

#### 2. Service Layer (`internal/service/`)

Application services implementing business use cases:

- **Building Service**: CRUD operations, validation
- **Class Service**: Space management, capacity tracking
- **Resource Service**: Availability management, type filtering
- **Lesson Service**: Schedule management, time calculations
- **Reservation Service**: Booking logic, conflict detection, availability checking

#### 3. Domain Layer (`internal/domain/`)

Core business logic and contracts:

**Entities**: Domain objects with business rules
**Use Cases**: Interface definitions for business operations
**Repositories**: Data access interface contracts

Each domain package contains:
- `entity.go`: Core business entity
- `usecase.go`: Business operation interface
- `repository.go`: Data access interface

#### 4. Adapter Layer (`internal/adapter/`)

Infrastructure and external service integration:

**Database** (`gorm/`):
- GORM-based repository implementations
- MySQL database integration
- Migration support through domain entities

**Configuration** (`config/`):
- Viper-based configuration management
- Environment-specific settings
- Database connection parameters

## Application Entry Points

### HTTP Server (`cmd/server/`)

Primary application mode for API serving:
- Wire dependency injection initialization
- Gin router setup with middleware
- Database migration on startup
- Graceful shutdown handling

### CLI Application (`cmd/cli/`)

Command-line interface for administrative operations:
- REST API client integration
- Human-readable output formatting
- Support for all CRUD operations across entities

### Lambda Function (`cmd/lambda/`)

Serverless deployment option:
- API Gateway integration
- Container reuse for performance
- Release mode configuration for production

## Data Flow

### Request Processing

```
HTTP Request → Middleware → Handler → Service → Repository → Database
                ↓           ↓          ↓         ↓          ↓
             CORS/Auth    DTO→Entity  Business  Data      MySQL
             Logging      Validation   Logic    Access     
             Metrics      
```

### Response Flow

```
Database → Repository → Service → Handler → Middleware → HTTP Response
   ↓          ↓          ↓         ↓          ↓             ↓
 MySQL    Entity      Business   Entity→DTO  Headers     JSON
         Mapping      Rules      Mapping    Metrics     
```

## Technology Stack

### Core Technologies

- **Language**: Go 1.24+
- **Web Framework**: Gin (HTTP router, middleware)
- **Database**: MySQL 8.0 with GORM ORM
- **Dependency Injection**: Google Wire (compile-time)
- **Configuration**: Viper with YAML configuration files

### Development Tools

- **Documentation**: Swagger/OpenAPI generation via Swag
- **Testing**: Testify framework for assertions
- **Hot Reloading**: Air for development workflows
- **Linting**: golangci-lint for code quality
- **Debugging**: Delve debugger integration

### Infrastructure

- **Containers**: Docker with multi-stage builds
- **Orchestration**: Docker Compose for local development
- **Cloud**: AWS Lambda, API Gateway, RDS
- **Infrastructure as Code**: Terraform with Terragrunt

## Configuration Management

### Configuration Sources

1. **Default Settings** (`configs/default.yaml`): Base configuration
2. **Environment Overrides** (`configs/development.yaml`): Environment-specific settings
3. **Environment Variables**: Runtime overrides for sensitive data

### Key Configuration Areas

- **Server**: Port, timeouts, host binding
- **Database**: Connection parameters, pool settings
- **API**: Versioning, metadata, contact information
- **CORS**: Cross-origin resource sharing policies
- **Logging**: Level, format, output configuration

## Security Considerations

### Current Implementation

- **CORS**: Configurable cross-origin policies
- **Input Validation**: DTO validation and sanitization
- **SQL Injection**: GORM ORM provides protection
- **Error Handling**: Consistent error responses without data leakage

### Future Enhancements

- **Authentication**: JWT token support prepared in configuration
- **Authorization**: Role-based access control framework
- **Rate Limiting**: Request throttling configuration ready
- **Audit Logging**: Comprehensive operation tracking

## Deployment Options

### Local Development

- **Native Go**: Direct execution for fastest development cycle
- **Docker Compose**: Full environment with database and monitoring
- **Hot Reloading**: Air-based automatic recompilation

### Production Deployment

- **Container**: Docker image with distroless base for security
- **Serverless**: AWS Lambda with API Gateway integration
- **Traditional**: Binary deployment with external database

## Monitoring and Observability

### Metrics

- **Prometheus**: Application metrics collection
- **Custom Metrics**: Request counts, response times, error rates
- **Health Checks**: Endpoint for service monitoring

### Logging

- **Structured Logging**: JSON format for machine processing
- **Log Levels**: Configurable from debug to error
- **Request Tracing**: Comprehensive request lifecycle logging

## Scalability Considerations

### Horizontal Scaling

- **Stateless Design**: No server-side session storage
- **Database Connection Pooling**: Configurable pool sizes
- **Load Balancer Ready**: Health check endpoints available

### Performance Optimization

- **Minimal Dependencies**: Focused technology choices
- **Efficient Queries**: GORM optimization for database access
- **Resource Management**: Proper connection cleanup and timeouts 