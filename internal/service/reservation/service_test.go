package reservation

import (
	"context"
	"fmt"
	"testing"
	"time"

	"sarc-ng/internal/domain/common"
	"sarc-ng/internal/domain/reservation"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of reservation.Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ReadReservationList(ctx context.Context, params common.PaginationParams) ([]reservation.Reservation, int64, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]reservation.Reservation), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) ReadReservation(ctx context.Context, id uint) (*reservation.Reservation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reservation.Reservation), args.Error(1)
}

func (m *MockRepository) CreateReservation(ctx context.Context, r *reservation.Reservation) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockRepository) UpdateReservation(ctx context.Context, r *reservation.Reservation) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

func (m *MockRepository) DeleteReservation(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Test helper functions

// createTestTimeRange creates a valid future time range for testing
func createTestTimeRange() (start, end time.Time) {
	start = time.Now().Add(time.Hour)
	end = start.Add(2 * time.Hour)
	return start, end
}

// createValidReservation creates a valid reservation for testing
func createValidReservation(resourceID, userID uint, start, end time.Time, status string) *reservation.Reservation {
	return &reservation.Reservation{
		ResourceID: resourceID,
		UserID:     userID,
		StartTime:  start,
		EndTime:    end,
		Purpose:    "Team meeting",
		Status:     status,
	}
}

// setupMockService creates a mock repository and service for testing
func setupMockService() (*MockRepository, *Service) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)
	return mockRepo, service
}

func TestCreateReservation(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid reservation is created", func(t *testing.T) {
		mockRepo, service := setupMockService()
		start, end := createTestTimeRange()
		newReservation := createValidReservation(1, 100, start, end, "pending")

		mockRepo.On("ReadReservationList", ctx, mock.AnythingOfType("common.PaginationParams")).
			Return([]reservation.Reservation{}, int64(0), nil)
		mockRepo.On("CreateReservation", ctx, newReservation).Return(nil)

		err := service.CreateReservation(ctx, newReservation)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Empty resource ID returns error", func(t *testing.T) {
		_, service := setupMockService()
		start, end := createTestTimeRange()
		invalidReservation := createValidReservation(0, 100, start, end, "pending")

		err := service.CreateReservation(ctx, invalidReservation)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
		assert.Contains(t, err.Error(), "resource ID cannot be zero")
	})

	t.Run("Empty user ID returns error", func(t *testing.T) {
		_, service := setupMockService()
		start, end := createTestTimeRange()
		invalidReservation := createValidReservation(1, 0, start, end, "pending")

		err := service.CreateReservation(ctx, invalidReservation)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
		assert.Contains(t, err.Error(), "user ID cannot be zero")
	})

	t.Run("Empty purpose returns error", func(t *testing.T) {
		_, service := setupMockService()
		start, end := createTestTimeRange()
		invalidReservation := createValidReservation(1, 100, start, end, "pending")
		invalidReservation.Purpose = "   "

		err := service.CreateReservation(ctx, invalidReservation)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
		assert.Contains(t, err.Error(), "purpose cannot be empty")
	})

	t.Run("Start time after end time returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		start := time.Now().Add(2 * time.Hour)
		end := time.Now().Add(time.Hour) // End before start

		invalidReservation := &reservation.Reservation{
			ResourceID: 1,
			UserID:     100,
			StartTime:  start,
			EndTime:    end,
			Purpose:    "Meeting",
			Status:     "pending",
		}

		err := service.CreateReservation(ctx, invalidReservation)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
		assert.Contains(t, err.Error(), "start time must be before end time")
	})

	t.Run("Start time in the past returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		start := time.Now().Add(-time.Hour) // In the past
		end := start.Add(2 * time.Hour)

		invalidReservation := &reservation.Reservation{
			ResourceID: 1,
			UserID:     100,
			StartTime:  start,
			EndTime:    end,
			Purpose:    "Meeting",
			Status:     "pending",
		}

		err := service.CreateReservation(ctx, invalidReservation)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
		assert.Contains(t, err.Error(), "start time cannot be in the past")
	})

	t.Run("Conflicting reservation returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		start := time.Now().Add(time.Hour)
		end := start.Add(2 * time.Hour)

		existing := reservation.Reservation{
			ID:         1,
			ResourceID: 1,
			UserID:     50,
			StartTime:  start.Add(30 * time.Minute), // Overlapping
			EndTime:    end.Add(30 * time.Minute),
			Status:     "confirmed",
		}

		newReservation := &reservation.Reservation{
			ResourceID: 1,
			UserID:     100,
			StartTime:  start,
			EndTime:    end,
			Purpose:    "Team meeting",
			Status:     "pending",
		}

		mockRepo.On("ReadReservationList", ctx, mock.AnythingOfType("common.PaginationParams")).Return([]reservation.Reservation{existing}, int64(1), nil)

		err := service.CreateReservation(ctx, newReservation)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrConflict)
		assert.Contains(t, err.Error(), "not available")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Non-conflicting reservation succeeds (different resource)", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		start := time.Now().Add(time.Hour)
		end := start.Add(2 * time.Hour)

		existing := reservation.Reservation{
			ID:         1,
			ResourceID: 2, // Different resource
			UserID:     50,
			StartTime:  start,
			EndTime:    end,
			Status:     "confirmed",
		}

		newReservation := &reservation.Reservation{
			ResourceID: 1,
			UserID:     100,
			StartTime:  start,
			EndTime:    end,
			Purpose:    "Team meeting",
			Status:     "pending",
		}

		mockRepo.On("ReadReservationList", ctx, mock.AnythingOfType("common.PaginationParams")).Return([]reservation.Reservation{existing}, int64(1), nil)
		mockRepo.On("CreateReservation", ctx, newReservation).Return(nil)

		err := service.CreateReservation(ctx, newReservation)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Non-conflicting reservation succeeds (cancelled reservation)", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		start := time.Now().Add(time.Hour)
		end := start.Add(2 * time.Hour)

		existing := reservation.Reservation{
			ID:         1,
			ResourceID: 1,
			UserID:     50,
			StartTime:  start,
			EndTime:    end,
			Status:     "cancelled", // Cancelled, so no conflict
		}

		newReservation := &reservation.Reservation{
			ResourceID: 1,
			UserID:     100,
			StartTime:  start,
			EndTime:    end,
			Purpose:    "Team meeting",
			Status:     "pending",
		}

		mockRepo.On("ReadReservationList", ctx, mock.AnythingOfType("common.PaginationParams")).Return([]reservation.Reservation{existing}, int64(1), nil)
		mockRepo.On("CreateReservation", ctx, newReservation).Return(nil)

		err := service.CreateReservation(ctx, newReservation)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateReservation(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid update succeeds", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		start := time.Now().Add(time.Hour)
		end := start.Add(2 * time.Hour)

		existing := &reservation.Reservation{
			ID:         1,
			ResourceID: 1,
			UserID:     100,
			StartTime:  start,
			EndTime:    end,
			Purpose:    "Meeting",
			Status:     "pending",
		}

		updated := &reservation.Reservation{
			ID:         1,
			ResourceID: 1,
			UserID:     100,
			StartTime:  start,
			EndTime:    end,
			Purpose:    "Updated Meeting",
			Status:     "confirmed",
		}

		mockRepo.On("ReadReservation", ctx, uint(1)).Return(existing, nil)
		mockRepo.On("UpdateReservation", ctx, updated).Return(nil)

		err := service.UpdateReservation(ctx, updated)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Zero ID returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		start := time.Now().Add(time.Hour)
		end := start.Add(2 * time.Hour)

		invalid := &reservation.Reservation{
			ID:         0,
			ResourceID: 1,
			UserID:     100,
			StartTime:  start,
			EndTime:    end,
			Purpose:    "Meeting",
			Status:     "pending",
		}

		err := service.UpdateReservation(ctx, invalid)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
	})

	t.Run("Time change checks for conflicts", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		originalStart := time.Now().Add(time.Hour)
		originalEnd := originalStart.Add(2 * time.Hour)

		newStart := time.Now().Add(3 * time.Hour)
		newEnd := newStart.Add(2 * time.Hour)

		existing := &reservation.Reservation{
			ID:         1,
			ResourceID: 1,
			UserID:     100,
			StartTime:  originalStart,
			EndTime:    originalEnd,
			Purpose:    "Meeting",
			Status:     "pending",
		}

		conflicting := reservation.Reservation{
			ID:         2,
			ResourceID: 1,
			UserID:     50,
			StartTime:  newStart.Add(30 * time.Minute),
			EndTime:    newEnd.Add(30 * time.Minute),
			Status:     "confirmed",
		}

		updated := &reservation.Reservation{
			ID:         1,
			ResourceID: 1,
			UserID:     100,
			StartTime:  newStart, // New time
			EndTime:    newEnd,
			Purpose:    "Updated Meeting",
			Status:     "pending",
		}

		mockRepo.On("ReadReservation", ctx, uint(1)).Return(existing, nil)
		mockRepo.On("ReadReservationList", ctx, mock.AnythingOfType("common.PaginationParams")).Return([]reservation.Reservation{conflicting, *existing}, int64(2), nil)

		err := service.UpdateReservation(ctx, updated)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrConflict)
		mockRepo.AssertExpectations(t)
	})
}

func TestCancelReservation(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid cancel succeeds", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		start := time.Now().Add(time.Hour)
		end := start.Add(2 * time.Hour)

		existing := &reservation.Reservation{
			ID:         1,
			ResourceID: 1,
			UserID:     100,
			StartTime:  start,
			EndTime:    end,
			Purpose:    "Meeting",
			Status:     "confirmed",
		}

		mockRepo.On("ReadReservation", ctx, uint(1)).Return(existing, nil)
		mockRepo.On("UpdateReservation", ctx, mock.MatchedBy(func(r *reservation.Reservation) bool {
			return r.ID == 1 && r.Status == "cancelled"
		})).Return(nil)

		err := service.CancelReservation(ctx, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Already cancelled returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		start := time.Now().Add(time.Hour)
		end := start.Add(2 * time.Hour)

		existing := &reservation.Reservation{
			ID:         1,
			ResourceID: 1,
			UserID:     100,
			StartTime:  start,
			EndTime:    end,
			Purpose:    "Meeting",
			Status:     "cancelled",
		}

		mockRepo.On("ReadReservation", ctx, uint(1)).Return(existing, nil)

		err := service.CancelReservation(ctx, 1)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrConflict)
		assert.Contains(t, err.Error(), "already cancelled")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Not found returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		mockRepo.On("ReadReservation", ctx, uint(999)).Return(nil, fmt.Errorf("not found: %w", common.ErrNotFound))

		err := service.CancelReservation(ctx, 999)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestCheckReservationAvailability(t *testing.T) {
	ctx := context.Background()

	t.Run("Available slot returns true", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		start := time.Now().Add(time.Hour)
		end := start.Add(2 * time.Hour)

		mockRepo.On("ReadReservationList", ctx, mock.AnythingOfType("common.PaginationParams")).Return([]reservation.Reservation{}, int64(0), nil)

		available, err := service.CheckReservationAvailability(ctx, 1, start, end)

		assert.NoError(t, err)
		assert.True(t, available)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Conflicting reservation returns false", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		start := time.Now().Add(time.Hour)
		end := start.Add(2 * time.Hour)

		existing := reservation.Reservation{
			ID:         1,
			ResourceID: 1,
			UserID:     50,
			StartTime:  start.Add(30 * time.Minute),
			EndTime:    end.Add(30 * time.Minute),
			Status:     "confirmed",
		}

		mockRepo.On("ReadReservationList", ctx, mock.AnythingOfType("common.PaginationParams")).Return([]reservation.Reservation{existing}, int64(1), nil)

		available, err := service.CheckReservationAvailability(ctx, 1, start, end)

		assert.NoError(t, err)
		assert.False(t, available)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Adjacent reservations don't conflict", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		start := time.Now().Add(time.Hour)
		end := start.Add(2 * time.Hour)

		existing := reservation.Reservation{
			ID:         1,
			ResourceID: 1,
			UserID:     50,
			StartTime:  end, // Starts exactly when we end
			EndTime:    end.Add(2 * time.Hour),
			Status:     "confirmed",
		}

		mockRepo.On("ReadReservationList", ctx, mock.AnythingOfType("common.PaginationParams")).Return([]reservation.Reservation{existing}, int64(1), nil)

		available, err := service.CheckReservationAvailability(ctx, 1, start, end)

		assert.NoError(t, err)
		assert.True(t, available, "Adjacent reservations should not conflict")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Past time returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		start := time.Now().Add(-time.Hour)
		end := time.Now()

		available, err := service.CheckReservationAvailability(ctx, 1, start, end)

		assert.Error(t, err)
		assert.False(t, available)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
		assert.Contains(t, err.Error(), "start time cannot be in the past")
	})
}

func TestDeleteReservation(t *testing.T) {
	ctx := context.Background()

	t.Run("Valid delete succeeds", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		start := time.Now().Add(time.Hour)
		end := start.Add(2 * time.Hour)

		existing := &reservation.Reservation{
			ID:         1,
			ResourceID: 1,
			UserID:     100,
			StartTime:  start,
			EndTime:    end,
			Purpose:    "Meeting",
			Status:     "pending",
		}

		mockRepo.On("ReadReservation", ctx, uint(1)).Return(existing, nil)
		mockRepo.On("DeleteReservation", ctx, uint(1)).Return(nil)

		err := service.DeleteReservation(ctx, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Zero ID returns error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		service := NewService(mockRepo)

		err := service.DeleteReservation(ctx, 0)

		assert.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
	})
}

func TestTimeRangesOverlap(t *testing.T) {
	now := time.Now()

	t.Run("Overlapping ranges", func(t *testing.T) {
		start1 := now
		end1 := now.Add(2 * time.Hour)
		start2 := now.Add(1 * time.Hour)
		end2 := now.Add(3 * time.Hour)

		assert.True(t, timeRangesOverlap(start1, end1, start2, end2))
	})

	t.Run("Non-overlapping ranges (before)", func(t *testing.T) {
		start1 := now
		end1 := now.Add(1 * time.Hour)
		start2 := now.Add(2 * time.Hour)
		end2 := now.Add(3 * time.Hour)

		assert.False(t, timeRangesOverlap(start1, end1, start2, end2))
	})

	t.Run("Non-overlapping ranges (after)", func(t *testing.T) {
		start1 := now.Add(2 * time.Hour)
		end1 := now.Add(3 * time.Hour)
		start2 := now
		end2 := now.Add(1 * time.Hour)

		assert.False(t, timeRangesOverlap(start1, end1, start2, end2))
	})

	t.Run("Adjacent ranges (no overlap)", func(t *testing.T) {
		start1 := now
		end1 := now.Add(1 * time.Hour)
		start2 := end1
		end2 := now.Add(2 * time.Hour)

		assert.False(t, timeRangesOverlap(start1, end1, start2, end2))
	})
}

