package building

import (
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

// ReadBuildingList retrieves all buildings
func (m *MockRepository) ReadBuildingList() ([]building.Building, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]building.Building), args.Error(1)
}

// ReadBuilding retrieves a building by ID
func (m *MockRepository) ReadBuilding(id uint) (*building.Building, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*building.Building), args.Error(1)
}

// FindBuildingByCode retrieves a building by code
func (m *MockRepository) FindBuildingByCode(code string) (*building.Building, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*building.Building), args.Error(1)
}

// CreateBuilding creates a new building
func (m *MockRepository) CreateBuilding(b *building.Building) error {
	args := m.Called(b)
	return args.Error(0)
}

// UpdateBuilding updates an existing building
func (m *MockRepository) UpdateBuilding(b *building.Building) error {
	args := m.Called(b)
	return args.Error(0)
}

// DeleteBuilding removes a building
func (m *MockRepository) DeleteBuilding(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetBuilding(t *testing.T) {
	t.Run("Valid ID returns building", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		expectedBuilding := &building.Building{
			ID:   1,
			Name: "Main Building",
			Code: "MB01",
		}

		mockRepo.On("ReadBuilding", uint(1)).Return(expectedBuilding, nil)

		result, err := service.GetBuilding(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedBuilding, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Zero ID returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		result, err := service.GetBuilding(0)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
	})

	t.Run("Not found returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		mockRepo.On("ReadBuilding", uint(999)).Return(nil, fmt.Errorf("not found: %w", common.ErrNotFound))

		result, err := service.GetBuilding(999)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestCreateBuilding(t *testing.T) {
	t.Run("Valid building is created", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		newBuilding := &building.Building{
			Name: "New Building",
			Code: "NB01",
		}

		mockRepo.On("FindBuildingByCode", "NB01").Return(nil, nil)
		mockRepo.On("CreateBuilding", newBuilding).Return(nil)

		err := service.CreateBuilding(newBuilding)

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

		err := service.CreateBuilding(invalidBuilding)

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

		err := service.CreateBuilding(invalidBuilding)

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

		mockRepo.On("FindBuildingByCode", "EB01").Return(existingBuilding, nil)

		err := service.CreateBuilding(newBuilding)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrConflict)
		assert.Contains(t, err.Error(), "already exists")
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateBuilding(t *testing.T) {
	t.Run("Valid update succeeds", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		updateBuilding := &building.Building{
			ID:   1,
			Name: "Updated Building",
			Code: "UB01",
		}

		mockRepo.On("FindBuildingByCode", "UB01").Return(nil, nil)
		mockRepo.On("UpdateBuilding", updateBuilding).Return(nil)

		err := service.UpdateBuilding(updateBuilding)

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

		err := service.UpdateBuilding(invalidBuilding)

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

		mockRepo.On("FindBuildingByCode", "OB01").Return(existingBuilding, nil)

		err := service.UpdateBuilding(updateBuilding)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrConflict)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteBuilding(t *testing.T) {
	t.Run("Valid delete succeeds", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		existingBuilding := &building.Building{
			ID:   1,
			Name: "Building to Delete",
			Code: "BD01",
		}

		mockRepo.On("ReadBuilding", uint(1)).Return(existingBuilding, nil)
		mockRepo.On("DeleteBuilding", uint(1)).Return(nil)

		err := service.DeleteBuilding(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Zero ID returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		err := service.DeleteBuilding(0)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
	})

	t.Run("Not found returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		mockRepo.On("ReadBuilding", uint(999)).Return(nil, fmt.Errorf("not found: %w", common.ErrNotFound))

		err := service.DeleteBuilding(999)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

