---
sidebar_position: 2
---

# Component Details

This document provides detailed information about the SARC-NG system components and their implementation.

## Domain Components

The domain layer contains the core business logic and entities of the system.

### Building

Buildings represent physical structures and locations.

**Entity Properties**
- `ID`: Unique identifier
- `Name`: Building name
- `Code`: Short reference code
- `CreatedAt`: Creation timestamp
- `UpdatedAt`: Last update timestamp

**Repository Interface**
```go
type BuildingRepository interface {
    FindAll(ctx context.Context) ([]Building, error)
    FindByID(ctx context.Context, id int) (*Building, error)
    Create(ctx context.Context, building *Building) error
    Update(ctx context.Context, building *Building) error
    Delete(ctx context.Context, id int) error
}
```

**Use Case Interface**
```go
type BuildingUseCase interface {
    GetAll(ctx context.Context) ([]Building, error)
    GetByID(ctx context.Context, id int) (*Building, error)
    Create(ctx context.Context, building *Building) error
    Update(ctx context.Context, building *Building) error
    Delete(ctx context.Context, id int) error
}
```

### Class

Classes represent classrooms or physical spaces.

**Entity Properties**
- `ID`: Unique identifier
- `Name`: Class name
- `Capacity`: Maximum capacity of the space
- `CreatedAt`: Creation timestamp
- `UpdatedAt`: Last update timestamp

**Repository Interface**
```go
type ClassRepository interface {
    FindAll(ctx context.Context) ([]Class, error)
    FindByID(ctx context.Context, id int) (*Class, error)
    Create(ctx context.Context, class *Class) error
    Update(ctx context.Context, class *Class) error
    Delete(ctx context.Context, id int) error
}
```

### Resource

Resources are bookable items like rooms, equipment, etc.

**Entity Properties**
- `ID`: Unique identifier
- `Name`: Resource name
- `Type`: Resource type (e.g., room, equipment)
- `IsAvailable`: Availability status
- `Description`: Optional description
- `CreatedAt`: Creation timestamp
- `UpdatedAt`: Last update timestamp

**Repository Interface**
```go
type ResourceRepository interface {
    FindAll(ctx context.Context) ([]Resource, error)
    FindByID(ctx context.Context, id int) (*Resource, error)
    FindByType(ctx context.Context, resourceType string) ([]Resource, error)
    Create(ctx context.Context, resource *Resource) error
    Update(ctx context.Context, resource *Resource) error
    Delete(ctx context.Context, id int) error
}
```

### Lesson

Lessons represent teaching or training sessions.

**Entity Properties**
- `ID`: Unique identifier
- `Title`: Lesson title
- `Duration`: Duration in minutes
- `Description`: Optional lesson description
- `StartTime`: Scheduled start time
- `EndTime`: Scheduled end time
- `CreatedAt`: Creation timestamp
- `UpdatedAt`: Last update timestamp

**Repository Interface**
```go
type LessonRepository interface {
    FindAll(ctx context.Context) ([]Lesson, error)
    FindByID(ctx context.Context, id int) (*Lesson, error)
    FindByTimeRange(ctx context.Context, start, end time.Time) ([]Lesson, error)
    Create(ctx context.Context, lesson *Lesson) error
    Update(ctx context.Context, lesson *Lesson) error
    Delete(ctx context.Context, id int) error
}
```

### Reservation

Reservations track resource bookings.

**Entity Properties**
- `ID`: Unique identifier
- `ResourceID`: Associated resource ID
- `UserID`: User who made the reservation
- `StartTime`: Reservation start time
- `EndTime`: Reservation end time
- `Purpose`: Reason for reservation
- `Status`: Current status (pending, confirmed, cancelled)
- `Description`: Optional additional information
- `CreatedAt`: Creation timestamp
- `UpdatedAt`: Last update timestamp

**Repository Interface**
```go
type ReservationRepository interface {
    FindAll(ctx context.Context) ([]Reservation, error)
    FindByID(ctx context.Context, id int) (*Reservation, error)
    FindByResourceID(ctx context.Context, resourceID int) ([]Reservation, error)
    FindByUserID(ctx context.Context, userID int) ([]Reservation, error)
    FindByTimeRange(ctx context.Context, start, end time.Time) ([]Reservation, error)
    CheckConflicts(ctx context.Context, resourceID int, start, end time.Time, excludeID int) (bool, error)
    Create(ctx context.Context, reservation *Reservation) error
    Update(ctx context.Context, reservation *Reservation) error
    UpdateStatus(ctx context.Context, id int, status string) error
    Delete(ctx context.Context, id int) error
}
```

## Service Components

The service layer implements business use cases and orchestrates domain operations.

### Building Service

```go
type BuildingService struct {
    repo domain.BuildingRepository
}

func NewBuildingService(repo domain.BuildingRepository) *BuildingService {
    return &BuildingService{repo: repo}
}

func (s *BuildingService) GetAll(ctx context.Context) ([]domain.Building, error) {
    return s.repo.FindAll(ctx)
}

// Other methods implementing domain.BuildingUseCase...
```

### Reservation Service

```go
type ReservationService struct {
    repo domain.ReservationRepository
    resourceRepo domain.ResourceRepository
}

func NewReservationService(repo domain.ReservationRepository, resourceRepo domain.ResourceRepository) *ReservationService {
    return &ReservationService{
        repo: repo,
        resourceRepo: resourceRepo,
    }
}

func (s *ReservationService) Create(ctx context.Context, reservation *domain.Reservation) error {
    // Check if resource exists
    resource, err := s.resourceRepo.FindByID(ctx, reservation.ResourceID)
    if err != nil {
        return err
    }
    if resource == nil {
        return domain.ErrResourceNotFound
    }

    // Check for conflicts
    hasConflict, err := s.repo.CheckConflicts(
        ctx,
        reservation.ResourceID,
        reservation.StartTime,
        reservation.EndTime,
        0,
    )
    if err != nil {
        return err
    }
    if hasConflict {
        return domain.ErrReservationConflict
    }

    // Create reservation
    return s.repo.Create(ctx, reservation)
}

// Other methods implementing domain.ReservationUseCase...
```

## Transport Components

The transport layer handles external communication via REST API, CLI, and Lambda.

### REST API Handlers

```go
// BuildingHandler manages HTTP requests for buildings
type BuildingHandler struct {
    baseHandler common.BaseHandler
    buildingUseCase domain.BuildingUseCase
}

func NewBuildingHandler(baseHandler common.BaseHandler, buildingUseCase domain.BuildingUseCase) *BuildingHandler {
    return &BuildingHandler{
        baseHandler: baseHandler,
        buildingUseCase: buildingUseCase,
    }
}

// GetBuildings godoc
// @Summary Get all buildings
// @Description Retrieve all buildings
// @Tags buildings
// @Accept json
// @Produce json
// @Success 200 {array} dto.BuildingResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /api/v1/buildings [get]
func (h *BuildingHandler) GetBuildings(c *gin.Context) {
    ctx := c.Request.Context()

    buildings, err := h.buildingUseCase.GetAll(ctx)
    if err != nil {
        h.baseHandler.HandleError(c, err)
        return
    }

    response := mapper.MapBuildingsToResponse(buildings)
    h.baseHandler.HandleSuccess(c, response)
}

// Other handler methods for CRUD operations...
```

### CLI Commands

```go
// BuildingCommand represents the building CLI commands
type BuildingCommand struct {
    client *client.Client
}

func NewBuildingCommand(client *client.Client) *BuildingCommand {
    return &BuildingCommand{client: client}
}

// CreateBuildingCommand returns a command to create a building
func (b *BuildingCommand) CreateBuildingCommand() *cobra.Command {
    var name, code string

    cmd := &cobra.Command{
        Use:   "create",
        Short: "Create a new building",
        RunE: func(cmd *cobra.Command, args []string) error {
            building, err := b.client.Buildings.Create(name, code)
            if err != nil {
                return err
            }

            return b.printBuilding(cmd.OutOrStdout(), building)
        },
    }

    cmd.Flags().StringVarP(&name, "name", "n", "", "Building name")
    cmd.Flags().StringVarP(&code, "code", "c", "", "Building code")
    cmd.MarkFlagRequired("name")
    cmd.MarkFlagRequired("code")

    return cmd
}

// Other CLI commands...
```

## Adapter Components

The adapter layer handles external dependencies and infrastructure concerns.

### Database Adapter (GORM)

```go
// BuildingAdapter implements the domain.BuildingRepository interface
type BuildingAdapter struct {
    db *gorm.DB
}

func NewBuildingAdapter(db *gorm.DB) *BuildingAdapter {
    return &BuildingAdapter{db: db}
}

// Model represents the database model for buildings
type BuildingModel struct {
    ID        uint      `gorm:"primarykey"`
    Name      string    `gorm:"size:255;not null"`
    Code      string    `gorm:"size:50;not null;uniqueIndex"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (a *BuildingAdapter) FindAll(ctx context.Context) ([]domain.Building, error) {
    var models []BuildingModel
    result := a.db.WithContext(ctx).Find(&models)
    if result.Error != nil {
        return nil, result.Error
    }

    return mapToDomainBuildings(models), nil
}

// Other adapter methods implementing domain.BuildingRepository...

// Helper functions for mapping between domain entities and database models
func mapToDomainBuilding(model BuildingModel) domain.Building {
    return domain.Building{
        ID:        int(model.ID),
        Name:      model.Name,
        Code:      model.Code,
        CreatedAt: model.CreatedAt,
        UpdatedAt: model.UpdatedAt,
    }
}

func mapToDomainBuildings(models []BuildingModel) []domain.Building {
    buildings := make([]domain.Building, len(models))
    for i, model := range models {
        buildings[i] = mapToDomainBuilding(model)
    }
    return buildings
}
```

## Application Entry Points

### HTTP Server (`cmd/server/`)

```go
func main() {
    // Initialize dependencies using Wire
    app, cleanup, err := InitializeApplication()
    if err != nil {
        log.Fatalf("failed to initialize application: %v", err)
    }
    defer cleanup()

    // Start server
    if err := app.Run(); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}
```

### Lambda Handler (`cmd/lambda/`)

```go
func main() {
    // Initialize dependencies using Wire
    handler, cleanup, err := InitializeHandler()
    if err != nil {
        log.Fatalf("failed to initialize lambda handler: %v", err)
    }
    defer cleanup()

    // Start Lambda handler
    lambda.Start(handler.Handle)
}
```

### CLI Application (`cmd/cli/`)

```go
func main() {
    // Initialize root command
    cmd := root.NewRootCommand()

    // Execute
    if err := cmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```
