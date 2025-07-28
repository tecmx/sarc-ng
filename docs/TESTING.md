# Testing

Testing strategies, procedures, and best practices for SARC-NG.

## Testing Overview

SARC-NG implements a comprehensive testing strategy with multiple layers to ensure reliability and maintainability. The testing approach follows Go best practices and includes unit tests, integration tests, and API validation.

### Testing Philosophy

- **Fast Feedback**: Unit tests provide immediate feedback during development
- **Confidence**: Integration tests validate complete workflows
- **Reliability**: Tests run consistently across different environments
- **Maintainability**: Tests are easy to understand and modify

## Testing Levels

### 1. Unit Tests

**Purpose**: Test individual functions and methods in isolation

**Characteristics:**
- Fast execution (< 100ms per test)
- No external dependencies
- High code coverage
- Test business logic and edge cases

**Location**: Alongside source code (`*_test.go` files)

**Example Structure:**
```go
func TestResourceService_Create(t *testing.T) {
    tests := []struct {
        name    string
        input   *resource.Resource
        want    error
        wantErr bool
    }{
        {
            name: "valid resource",
            input: &resource.Resource{
                Name: "Test Resource",
                Type: "equipment",
            },
            want:    nil,
            wantErr: false,
        },
        {
            name: "invalid resource - empty name",
            input: &resource.Resource{
                Name: "",
                Type: "equipment",
            },
            want:    ErrInvalidResourceName,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### 2. Integration Tests

**Purpose**: Test complete request/response flows with real dependencies

**Characteristics:**
- Use real database via Docker
- Test HTTP endpoints end-to-end
- Validate API contracts
- Test error scenarios

**Location**: `test/integration/` directory

**Build Tag**: `//go:build integration`

**Current Implementation:**
```go
//go:build integration
// +build integration

package integration

import (
    "encoding/json"
    "net/http"
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestAPIEndpoints(t *testing.T) {
    endpoints := []struct {
        name string
        path string
    }{
        {"Buildings", "/api/v1/buildings"},
        {"Classes", "/api/v1/classes"},
        {"Lessons", "/api/v1/lessons"},
        {"Reservations", "/api/v1/reservations"},
        {"Resources", "/api/v1/resources"},
    }

    for _, endpoint := range endpoints {
        t.Run(endpoint.name, func(t *testing.T) {
            resp, err := http.Get(baseURL + endpoint.path)
            require.NoError(t, err)
            defer resp.Body.Close()

            assert.Equal(t, http.StatusOK, resp.StatusCode)
            
            var result interface{}
            err = json.NewDecoder(resp.Body).Decode(&result)
            require.NoError(t, err)
        })
    }
}
```

### 3. HTTP Testing

**Purpose**: Test HTTP handlers and middleware in isolation

**Characteristics:**
- Use `httptest` package
- Mock dependencies
- Test middleware behavior
- Validate response formats

**Example:**
```go
func TestBuildingHandler_GetAll(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    mockService := &mockBuildingService{}
    handler := building.NewHandler(mockService)
    
    router := gin.New()
    router.GET("/buildings", handler.GetAll)
    
    req := httptest.NewRequest(http.MethodGet, "/buildings", nil)
    w := httptest.NewRecorder()
    
    router.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

## Testing Commands

### Local Testing

**Unit Tests:**
```bash
# Run all unit tests
make test

# Run tests with coverage
make coverage

# Run tests with race detection
go test -race ./...

# Run specific package tests
go test ./internal/service/building/...

# Verbose output
go test -v ./...
```

**Integration Tests:**
```bash
# Start Docker environment and run integration tests
docker compose up -d && \
docker compose run --rm app go test -tags=integration ./test/integration/... && \
docker compose down -v --remove-orphans

# Run integration tests with verbose output
docker compose up -d && \
docker compose run --rm app go test -v -tags=integration ./test/integration/... && \
docker compose down -v --remove-orphans
```

### Test Coverage

**Generate Coverage Report:**
```bash
# Generate HTML coverage report
make coverage

# View coverage in terminal
go test -cover ./...

# Detailed coverage per function
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Open HTML report
open build/coverage.html
```

**Coverage Targets:**
- **Domain Layer**: 90%+ coverage
- **Service Layer**: 85%+ coverage
- **Handler Layer**: 80%+ coverage
- **Overall Project**: 75%+ coverage

## Testing Frameworks and Libraries

### Core Testing Libraries

**Standard Library:**
- `testing` - Core testing framework
- `httptest` - HTTP testing utilities
- `net/http/httptest` - Test HTTP servers

**Third-Party Libraries:**
- `github.com/stretchr/testify` - Assertions and test utilities
  - `assert` - Basic assertions
  - `require` - Assertions that stop test execution
  - `mock` - Mock generation
  - `suite` - Test suites

### Assertion Examples

**Basic Assertions:**
```go
// Equality
assert.Equal(t, expected, actual)
assert.NotEqual(t, unexpected, actual)

// Nil checks
assert.Nil(t, err)
assert.NotNil(t, result)

// Boolean assertions
assert.True(t, condition)
assert.False(t, condition)

// String assertions
assert.Contains(t, "hello world", "world")
assert.Empty(t, "")
assert.NotEmpty(t, "content")

// Numeric assertions
assert.Greater(t, 5, 3)
assert.Less(t, 3, 5)
assert.Zero(t, 0)
```

**Requirements (Stop on Failure):**
```go
// Use require for critical assertions
require.NoError(t, err)
require.NotNil(t, result)
require.Equal(t, expected, actual)
```

## Test Data Management

### Test Database

**Docker-based Testing:**
```yaml
# docker-compose.yml test configuration
services:
  db:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: testpassword
      MYSQL_DATABASE: sarcng_test
    ports:
      - "3306:3306"
```

**Database Setup:**
```go
func setupTestDB(t *testing.T) *gorm.DB {
    dsn := "root:testpassword@tcp(localhost:3306)/sarcng_test?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    require.NoError(t, err)
    
    // Migrate tables
    err = db.AutoMigrate(&building.Building{}, &resource.Resource{})
    require.NoError(t, err)
    
    return db
}

func cleanupTestDB(t *testing.T, db *gorm.DB) {
    // Clean up test data
    db.Exec("DELETE FROM buildings")
    db.Exec("DELETE FROM resources")
}
```

### Test Fixtures

**Sample Data Creation:**
```go
func createTestBuilding() *building.Building {
    return &building.Building{
        Name: "Test Building",
        Code: "TEST-001",
    }
}

func createTestResource() *resource.Resource {
    return &resource.Resource{
        Name:        "Test Resource",
        Type:        "equipment",
        IsAvailable: true,
    }
}
```

**Factory Pattern:**
```go
type BuildingFactory struct{}

func (f *BuildingFactory) Create(opts ...func(*building.Building)) *building.Building {
    b := &building.Building{
        Name: "Default Building",
        Code: "DEFAULT",
    }
    
    for _, opt := range opts {
        opt(b)
    }
    
    return b
}

func WithName(name string) func(*building.Building) {
    return func(b *building.Building) {
        b.Name = name
    }
}

// Usage
building := factory.Create(WithName("Custom Building"))
```

## Mocking and Test Doubles

### Interface Mocking

**Manual Mocks:**
```go
type mockBuildingRepository struct {
    buildings []building.Building
    err       error
}

func (m *mockBuildingRepository) FindAll() ([]building.Building, error) {
    if m.err != nil {
        return nil, m.err
    }
    return m.buildings, nil
}

func (m *mockBuildingRepository) Create(b *building.Building) error {
    if m.err != nil {
        return m.err
    }
    b.ID = uint(len(m.buildings) + 1)
    m.buildings = append(m.buildings, *b)
    return nil
}
```

**Service Layer Testing with Mocks:**
```go
func TestBuildingService_Create(t *testing.T) {
    mockRepo := &mockBuildingRepository{}
    service := buildingService.NewService(mockRepo)
    
    building := &building.Building{
        Name: "Test Building",
        Code: "TEST",
    }
    
    err := service.Create(building)
    
    assert.NoError(t, err)
    assert.NotZero(t, building.ID)
}
```

### Testify Mock

**Generated Mocks:**
```go
//go:generate mockery --name Repository --dir ./internal/domain/building --output ./mocks
type Repository interface {
    FindAll() ([]Building, error)
    Create(building *Building) error
}

// Usage in tests
func TestWithGeneratedMock(t *testing.T) {
    mockRepo := mocks.NewRepository(t)
    mockRepo.On("FindAll").Return([]building.Building{}, nil)
    
    service := buildingService.NewService(mockRepo)
    buildings, err := service.GetAll()
    
    assert.NoError(t, err)
    assert.Empty(t, buildings)
    mockRepo.AssertExpectations(t)
}
```

## Testing Patterns

### Table-Driven Tests

**Structure:**
```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    bool
        wantErr error
    }{
        {
            name:    "valid input",
            input:   "valid@example.com",
            want:    true,
            wantErr: nil,
        },
        {
            name:    "invalid input",
            input:   "invalid-email",
            want:    false,
            wantErr: ErrInvalidEmail,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ValidateEmail(tt.input)
            
            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}
```

### Test Suites

**Testify Suites:**
```go
type ServiceTestSuite struct {
    suite.Suite
    service building.Usecase
    repo    *mockRepository
}

func (s *ServiceTestSuite) SetupTest() {
    s.repo = &mockRepository{}
    s.service = buildingService.NewService(s.repo)
}

func (s *ServiceTestSuite) TestCreate() {
    building := &building.Building{Name: "Test", Code: "TEST"}
    
    err := s.service.Create(building)
    
    s.NoError(err)
    s.NotZero(building.ID)
}

func TestServiceSuite(t *testing.T) {
    suite.Run(t, new(ServiceTestSuite))
}
```

### Helper Functions

**Common Test Helpers:**
```go
func assertHTTPError(t *testing.T, w *httptest.ResponseRecorder, expectedCode int) {
    assert.Equal(t, expectedCode, w.Code)
    
    var errorResp map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &errorResp)
    assert.NoError(t, err)
    assert.Contains(t, errorResp, "error")
}

func createTestContext() (*gin.Context, *httptest.ResponseRecorder) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    return c, w
}
```

## Testing Best Practices

### Test Organization

**File Naming:**
- `*_test.go` for unit tests in same package
- `*_integration_test.go` for integration tests
- Separate package for external tests (`package_test`)

**Test Naming:**
- Function: `TestFunction_Scenario`
- Method: `TestStruct_Method_Scenario`
- Subtests: Use descriptive names

**Example:**
```go
func TestBuildingService_Create_ValidInput(t *testing.T) { }
func TestBuildingService_Create_InvalidInput(t *testing.T) { }
func TestBuildingHandler_GetAll_EmptyDatabase(t *testing.T) { }
```

### Test Structure

**AAA Pattern (Arrange, Act, Assert):**
```go
func TestBuildingService_Create(t *testing.T) {
    // Arrange
    mockRepo := &mockRepository{}
    service := buildingService.NewService(mockRepo)
    building := &building.Building{Name: "Test", Code: "TEST"}
    
    // Act
    err := service.Create(building)
    
    // Assert
    assert.NoError(t, err)
    assert.NotZero(t, building.ID)
}
```

### Error Testing

**Test Error Conditions:**
```go
func TestService_HandleErrors(t *testing.T) {
    tests := []struct {
        name        string
        setupMock   func(*mockRepository)
        expectError bool
        errorType   error
    }{
        {
            name: "database error",
            setupMock: func(m *mockRepository) {
                m.err = errors.New("database error")
            },
            expectError: true,
            errorType:   ErrDatabaseError,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := &mockRepository{}
            tt.setupMock(mockRepo)
            
            service := buildingService.NewService(mockRepo)
            _, err := service.GetAll()
            
            if tt.expectError {
                assert.Error(t, err)
                assert.ErrorIs(t, err, tt.errorType)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## Continuous Integration Testing

### GitHub Actions Integration

**Test Workflow:**
```yaml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: testpassword
          MYSQL_DATABASE: sarcng_test
        ports:
          - 3306:3306
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.24
    
    - name: Run tests
      run: |
        go test -race -coverprofile=coverage.out ./...
        go tool cover -html=coverage.out -o coverage.html
    
    - name: Upload coverage
      uses: actions/upload-artifact@v3
      with:
        name: coverage-report
        path: coverage.html
```

### Docker Testing

**Test in Docker Environment:**
```bash
# Build test image
docker build --target development -t sarc-ng:test .

# Run tests in container
docker run --rm \
  -v $(pwd):/app \
  -w /app \
  sarc-ng:test \
  go test -race ./...

# Integration tests with compose
docker compose -f docker-compose.test.yml up --abort-on-container-exit
```

## Performance Testing

### Benchmark Tests

**Basic Benchmarks:**
```go
func BenchmarkBuildingService_Create(b *testing.B) {
    service := setupService()
    building := createTestBuilding()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.Create(building)
    }
}

func BenchmarkBuildingService_GetAll(b *testing.B) {
    service := setupServiceWithData()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.GetAll()
    }
}
```

**Run Benchmarks:**
```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkBuildingService_Create ./internal/service/building

# Memory allocation profiling
go test -bench=. -benchmem ./...
```

### Load Testing

**API Load Testing:**
```go
func TestAPI_LoadTest(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping load test in short mode")
    }
    
    client := &http.Client{Timeout: 5 * time.Second}
    var wg sync.WaitGroup
    errors := make(chan error, 100)
    
    // Simulate 100 concurrent requests
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            resp, err := client.Get("http://localhost:8080/api/v1/buildings")
            if err != nil {
                errors <- err
                return
            }
            defer resp.Body.Close()
            
            if resp.StatusCode != http.StatusOK {
                errors <- fmt.Errorf("unexpected status: %d", resp.StatusCode)
            }
        }()
    }
    
    wg.Wait()
    close(errors)
    
    errorCount := 0
    for err := range errors {
        t.Logf("Request error: %v", err)
        errorCount++
    }
    
    // Allow up to 5% error rate
    assert.Less(t, errorCount, 5)
}
```

## Test Maintenance

### Test Reliability

**Avoid Flaky Tests:**
- Don't depend on timing
- Clean up test data
- Use deterministic test data
- Avoid shared global state

**Stable Test Data:**
```go
func TestWithCleanup(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    // Test implementation
}
```

### Test Documentation

**Self-Documenting Tests:**
```go
func TestReservationService_Create_RejectsOverlappingReservations(t *testing.T) {
    // Given: An existing reservation from 10:00 to 12:00
    existing := createReservation("10:00", "12:00")
    service.Create(existing)
    
    // When: Attempting to create overlapping reservation from 11:00 to 13:00
    overlapping := createReservation("11:00", "13:00")
    err := service.Create(overlapping)
    
    // Then: The reservation should be rejected with conflict error
    assert.Error(t, err)
    assert.ErrorIs(t, err, ErrReservationConflict)
}
```

### Debugging Tests

**Debug Test Failures:**
```bash
# Run specific test with verbose output
go test -v -run TestSpecificTest ./package

# Run tests with detailed output
go test -v -json ./... | jq

# Debug with delve
dlv test ./package -- -test.run TestSpecificTest
```

## Testing Tools

### Development Tools

**Code Coverage:**
```bash
# Install coverage tools
go install golang.org/x/tools/cmd/cover@latest

# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View coverage by function
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

**Static Analysis:**
```bash
# Vet for common errors
go vet ./...

# Lint with golangci-lint
golangci-lint run

# Race detection
go test -race ./...
```

**Test Utilities:**
```bash
# Find tests without assertions
grep -r "func Test" --include="*_test.go" . | \
  xargs -I {} sh -c 'file="{}"; if ! grep -q "assert\|require" "$file"; then echo "$file"; fi'

# Count test functions
find . -name "*_test.go" -exec grep -l "func Test" {} \; | \
  xargs grep "func Test" | wc -l
```

This comprehensive testing strategy ensures SARC-NG maintains high quality and reliability throughout development and deployment. 