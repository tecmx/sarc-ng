---
sidebar_position: 4
tags:
  - architecture
  - design
  - structure
---

# Architecture Overview

SARC-NG follows Clean Architecture principles with a clear separation of concerns and dependency inversion.

## High-Level Architecture

```mermaid
graph TB
    subgraph "Client Layer"
        Web[Web UI]
        CLI[CLI Tool]
        API[External APIs]
    end
    
    subgraph "Transport Layer"
        HTTP[HTTP/REST]
        GraphQL[GraphQL]
        gRPC[gRPC]
    end
    
    subgraph "Application Layer"
        Server[HTTP Server]
        Handlers[HTTP Handlers]
        Middleware[Middleware]
    end
    
    subgraph "Business Layer"
        Services[Application Services]
        Domain[Domain Logic]
        Entities[Domain Entities]
    end
    
    subgraph "Data Layer"
        Repositories[Repositories]
        Database[(PostgreSQL)]
        Cache[(Redis)]
        Files[(File Storage)]
    end
    
    Web --> HTTP
    CLI --> Server
    API --> HTTP
    
    HTTP --> Handlers
    GraphQL --> Handlers
    gRPC --> Handlers
    
    Handlers --> Services
    Middleware --> Handlers
    
    Services --> Domain
    Services --> Repositories
    
    Repositories --> Database
    Repositories --> Cache
    Repositories --> Files
```

## Project Structure

```
sarc-ng/
├── cmd/                    # Application entrypoints
│   ├── cli/               # Command-line interface
│   ├── lambda/            # AWS Lambda function
│   └── server/            # HTTP API server
│
├── internal/              # Internal application code
│   ├── adapter/           # External service adapters
│   │   ├── db/           # Database interfaces
│   │   └── gorm/         # GORM implementations
│   │
│   ├── domain/            # Business domain logic
│   │   ├── building/     # Building entity & logic
│   │   ├── class/        # Class entity & logic
│   │   ├── lesson/       # Lesson entity & logic
│   │   ├── reservation/  # Reservation entity & logic
│   │   ├── resource/     # Resource entity & logic
│   │   └── common/       # Shared domain logic
│   │
│   ├── service/           # Application services
│   │   ├── building/     # Building service
│   │   ├── class/        # Class service
│   │   ├── lesson/       # Lesson service
│   │   ├── reservation/  # Reservation service
│   │   └── resource/     # Resource service
│   │
│   └── transport/         # Transport layer
│       ├── rest/         # REST API handlers
│       └── common/       # Shared transport logic
│
├── pkg/                   # Public packages
│   ├── metrics/          # Metrics collection
│   └── rest/             # REST client utilities
│
├── configs/               # Configuration files
├── docker/               # Docker configurations
├── infrastructure/       # Infrastructure as Code
└── test/                 # Integration tests
```

## Core Components

### 1. Domain Layer (`internal/domain/`)

The domain layer contains the business logic and entities:

```go
// Domain Entity
type Building struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Code      string    `json:"code"`
    Address   string    `json:"address,omitempty"`
    Floors    int       `json:"floors,omitempty"`
    CreatedAt time.Time `json:"createdAt"`
    UpdatedAt time.Time `json:"updatedAt"`
}

// Domain Repository Interface
type BuildingRepository interface {
    Create(ctx context.Context, building *Building) error
    GetByID(ctx context.Context, id int) (*Building, error)
    GetAll(ctx context.Context) ([]*Building, error)
    Update(ctx context.Context, building *Building) error
    Delete(ctx context.Context, id int) error
}

// Domain Service
type BuildingService struct {
    repo BuildingRepository
}
```

### 2. Service Layer (`internal/service/`)

Application services orchestrate business operations:

```go
type BuildingService struct {
    buildingRepo domain.BuildingRepository
    logger       *slog.Logger
}

func (s *BuildingService) CreateBuilding(ctx context.Context, req CreateBuildingRequest) (*Building, error) {
    // Validation logic
    if err := s.validateBuilding(req); err != nil {
        return nil, err
    }
    
    // Business logic
    building := &Building{
        Name:    req.Name,
        Code:    strings.ToUpper(req.Code),
        Address: req.Address,
        Floors:  req.Floors,
    }
    
    // Persist
    if err := s.buildingRepo.Create(ctx, building); err != nil {
        return nil, err
    }
    
    return building, nil
}
```

### 3. Transport Layer (`internal/transport/`)

HTTP handlers handle request/response concerns:

```go
type BuildingHandler struct {
    buildingService service.BuildingService
}

func (h *BuildingHandler) CreateBuilding(w http.ResponseWriter, r *http.Request) {
    var req CreateBuildingRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    building, err := h.buildingService.CreateBuilding(r.Context(), req)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteStatus(http.StatusCreated)
    json.NewEncoder(w).Encode(building)
}
```

### 4. Data Layer (`internal/adapter/`)

Repository implementations handle data persistence:

```go
type GormBuildingRepository struct {
    db *gorm.DB
}

func (r *GormBuildingRepository) Create(ctx context.Context, building *Building) error {
    return r.db.WithContext(ctx).Create(building).Error
}

func (r *GormBuildingRepository) GetByID(ctx context.Context, id int) (*Building, error) {
    var building Building
    err := r.db.WithContext(ctx).First(&building, id).Error
    if err != nil {
        return nil, err
    }
    return &building, nil
}
```

## Design Patterns

### 1. Repository Pattern
- Abstracts data access logic
- Enables easy testing with mocks
- Supports multiple data sources

### 2. Dependency Injection
- Uses constructor injection
- Managed by dependency injection container
- Enables loose coupling

### 3. Clean Architecture
- Business logic independent of frameworks
- Testable architecture
- Database and UI agnostic core

### 4. Domain-Driven Design (DDD)
- Organized around business domains
- Rich domain models
- Ubiquitous language

## Data Flow

### Request Flow
```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant Service  
    participant Repository
    participant Database
    
    Client->>Handler: HTTP Request
    Handler->>Handler: Parse & Validate
    Handler->>Service: Business Operation
    Service->>Service: Business Logic
    Service->>Repository: Data Operation
    Repository->>Database: SQL Query
    Database-->>Repository: Result
    Repository-->>Service: Domain Entity
    Service-->>Handler: Business Result
    Handler->>Handler: Format Response
    Handler-->>Client: HTTP Response
```

### Error Handling Flow
```mermaid
graph TD
    A[Request] --> B{Validation}
    B -->|Invalid| C[400 Bad Request]
    B -->|Valid| D[Business Logic]
    D --> E{Business Rules}
    E -->|Violation| F[422 Unprocessable Entity]
    E -->|Valid| G[Data Layer]
    G --> H{Database Operation}
    H -->|Not Found| I[404 Not Found] 
    H -->|Conflict| J[409 Conflict]
    H -->|Error| K[500 Internal Error]
    H -->|Success| L[200/201 Success]
```

## Database Design

### Core Entities

```sql
-- Buildings
CREATE TABLE buildings (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    address TEXT,
    floors INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Resources (Classrooms, Equipment, Labs)
CREATE TABLE resources (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'classroom', 'equipment', 'lab'
    capacity INTEGER,
    location VARCHAR(255),
    building_id INTEGER REFERENCES buildings(id),
    details JSONB, -- Flexible attributes
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Classes  
CREATE TABLE classes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    instructor_name VARCHAR(255) NOT NULL,
    semester VARCHAR(100),
    capacity INTEGER,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Reservations
CREATE TABLE reservations (
    id SERIAL PRIMARY KEY,
    resource_id INTEGER REFERENCES resources(id),
    user_id INTEGER NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    purpose TEXT,
    status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'confirmed', 'cancelled'
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Lessons
CREATE TABLE lessons (
    id SERIAL PRIMARY KEY,
    class_id INTEGER REFERENCES classes(id),
    resource_id INTEGER REFERENCES resources(id),
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    repeat_pattern VARCHAR(50), -- 'once', 'daily', 'weekly', 'biweekly'
    repeat_until DATE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## Configuration

### Environment-Based Configuration

```yaml
# configs/development.yaml
server:
  port: 8080
  timeout: 30s
  
database:
  host: localhost
  port: 5432
  name: sarc_ng
  user: sarc
  password: password
  
redis:
  addr: localhost:6379
  db: 0
  
auth:
  jwt_secret: dev-secret-key
  token_ttl: 24h
  
logging:
  level: debug
  format: json
```

## Security

### Authentication & Authorization
- JWT-based authentication
- Role-based access control (RBAC)
- API key authentication for integrations
- Request rate limiting

### Data Protection
- Input validation and sanitization
- SQL injection prevention (parameterized queries)
- XSS protection
- CORS configuration

### Infrastructure Security
- HTTPS/TLS encryption
- Database connection encryption
- Secrets management
- Network security groups

## Monitoring & Observability

### Metrics
- Application metrics (Prometheus format)
- Database performance metrics
- HTTP request metrics
- Custom business metrics

### Logging
- Structured logging (JSON)
- Correlation IDs for request tracing
- Different log levels per environment
- Centralized log aggregation

### Health Checks
- Liveness probes
- Readiness probes  
- Dependency health checks
- Graceful shutdown

This architecture provides a solid foundation for maintainable, scalable, and testable code while following Go best practices and industry standards.
