package building

import (
	"context"
	"fmt"
	"testing"

	"sarc-ng/internal/domain/building"
	"sarc-ng/internal/domain/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of building.Repository
type MockRepository struct {
	mock.Mock
}

// ReadBuildingList retrieves buildings with pagination
func (m *MockRepository) ReadBuildingList(ctx context.Context, params common.PaginationParams) ([]building.Building, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]building.Building), args.Get(1).(int64), args.Error(2)
}

// ReadBuilding retrieves a building by ID
func (m *MockRepository) ReadBuilding(ctx context.Context, id uint) (*building.Building, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*building.Building), args.Error(1)
}

// FindBuildingByCode retrieves a building by code
func (m *MockRepository) FindBuildingByCode(ctx context.Context, code string) (*building.Building, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*building.Building), args.Error(1)
}

// CreateBuilding creates a new building
func (m *MockRepository) CreateBuilding(ctx context.Context, b *building.Building) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

// UpdateBuilding updates an existing building
func (m *MockRepository) UpdateBuilding(ctx context.Context, b *building.Building) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

// DeleteBuilding removes a building
func (m *MockRepository) DeleteBuilding(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGetBuilding(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid ID returns building", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		expectedBuilding := &building.Building{
			ID:   1,
			Name: "Main Building",
			Code: "MB01",
		}

		mockRepo.On("ReadBuilding", ctx, uint(1)).Return(expectedBuilding, nil)

		result, err := service.GetBuilding(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedBuilding, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Zero ID returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		result, err := service.GetBuilding(ctx, 0)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
	})

	t.Run("Not found returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		mockRepo.On("ReadBuilding", ctx, uint(999)).Return(nil, fmt.Errorf("not found: %w", common.ErrNotFound))

		result, err := service.GetBuilding(ctx, 999)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateBuilding(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid building is created", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		newBuilding := &building.Building{
			Name: "New Building",
			Code: "NB01",
		}

		mockRepo.On("FindBuildingByCode", ctx, "NB01").Return(nil, nil)
		mockRepo.On("CreateBuilding", ctx, newBuilding).Return(nil)

		err := service.CreateBuilding(ctx, newBuilding)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty name returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		invalidBuilding := &building.Building{
			Name: "  ",
			Code: "NB01",
		}

		err := service.CreateBuilding(ctx, invalidBuilding)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
		assert.Contains(t, err.Error(), "name cannot be empty")
	})

	t.Run("Empty code returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		invalidBuilding := &building.Building{
			Name: "New Building",
			Code: "  ",
		}

		err := service.CreateBuilding(ctx, invalidBuilding)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
		assert.Contains(t, err.Error(), "code cannot be empty")
	})

	t.Run("Duplicate code returns conflict error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		existingBuilding := &building.Building{
			ID:   1,
			Name: "Existing Building",
			Code: "EB01",
		}

		newBuilding := &building.Building{
			Name: "New Building",
			Code: "EB01",
		}

		mockRepo.On("FindBuildingByCode", ctx, "EB01").Return(existingBuilding, nil)

		err := service.CreateBuilding(ctx, newBuilding)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrConflict)
		assert.Contains(t, err.Error(), "already exists")
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateBuilding(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid update succeeds", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		updateBuilding := &building.Building{
			ID:   1,
			Name: "Updated Building",
			Code: "UB01",
		}

		mockRepo.On("FindBuildingByCode", ctx, "UB01").Return(nil, nil)
		mockRepo.On("UpdateBuilding", ctx, updateBuilding).Return(nil)

		err := service.UpdateBuilding(ctx, updateBuilding)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Zero ID returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		invalidBuilding := &building.Building{
			ID:   0,
			Name: "Building",
			Code: "B01",
		}

		err := service.UpdateBuilding(ctx, invalidBuilding)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
	})

	t.Run("Duplicate code for different building returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		existingBuilding := &building.Building{
			ID:   2,
			Name: "Other Building",
			Code: "OB01",
		}

		updateBuilding := &building.Building{
			ID:   1,
			Name: "Updated Building",
			Code: "OB01",
		}

		mockRepo.On("FindBuildingByCode", ctx, "OB01").Return(existingBuilding, nil)

		err := service.UpdateBuilding(ctx, updateBuilding)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrConflict)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteBuilding(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid delete succeeds", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		existingBuilding := &building.Building{
			ID:   1,
			Name: "Building to Delete",
			Code: "BD01",
		}

		mockRepo.On("ReadBuilding", ctx, uint(1)).Return(existingBuilding, nil)
		mockRepo.On("DeleteBuilding", ctx, uint(1)).Return(nil)

		err := service.DeleteBuilding(ctx, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Zero ID returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		err := service.DeleteBuilding(ctx, 0)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
	})

	t.Run("Not found returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		mockRepo.On("ReadBuilding", ctx, uint(999)).Return(nil, fmt.Errorf("not found: %w", common.ErrNotFound))

		err := service.DeleteBuilding(ctx, 999)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
