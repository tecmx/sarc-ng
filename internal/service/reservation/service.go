package reservation

import (
	"context"
	"fmt"
	"sarc-ng/internal/domain/common"
	"sarc-ng/internal/domain/reservation"
	"strings"
	"time"
)

const (
	statusPending   = "pending"
	statusConfirmed = "confirmed"
	statusCancelled = "cancelled"
	statusRejected  = "rejected"
)

// Service implements reservation.Usecase interface
type Service struct {
	repo reservation.Repository
}

// Compile-time verification that Service implements reservation.Usecase
var _ reservation.Usecase = (*Service)(nil)

// NewService creates a new reservation service
func NewService(repo reservation.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetAllReservations retrieves all reservations with pagination
func (s *Service) GetAllReservations(ctx context.Context, params common.PaginationParams) ([]reservation.Reservation, common.PaginationResult, error) {
	// Validate all pagination parameters
	if err := reservation.ValidateParams(params); err != nil {
		return nil, common.PaginationResult{}, err
	}

	reservations, total, err := s.repo.ReadReservationList(ctx, params)
	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	result := common.NewPaginationResult(params.Page, params.PageSize, total)
	return reservations, result, nil
}

// GetReservation retrieves a reservation by ID with validation
func (s *Service) GetReservation(ctx context.Context, id uint) (*reservation.Reservation, error) {
	if id == 0 {
		return nil, fmt.Errorf("%w: reservation ID cannot be zero", common.ErrInvalidInput)
	}
	return s.repo.ReadReservation(ctx, id)
}

// CreateReservation creates a new reservation with validation
func (s *Service) CreateReservation(ctx context.Context, r *reservation.Reservation) error {
	// Validate resource ID
	if r.ResourceID == 0 {
		return fmt.Errorf("%w: resource ID cannot be zero", common.ErrInvalidInput)
	}

	// Validate user ID
	if r.UserID == 0 {
		return fmt.Errorf("%w: user ID cannot be zero", common.ErrInvalidInput)
	}

	// Validate purpose
	if strings.TrimSpace(r.Purpose) == "" {
		return fmt.Errorf("%w: reservation purpose cannot be empty", common.ErrInvalidInput)
	}

	// Validate time range
	if r.StartTime.IsZero() {
		return fmt.Errorf("%w: start time is required", common.ErrInvalidInput)
	}

	if r.EndTime.IsZero() {
		return fmt.Errorf("%w: end time is required", common.ErrInvalidInput)
	}

	if r.StartTime.After(r.EndTime) || r.StartTime.Equal(r.EndTime) {
		return fmt.Errorf("%w: start time must be before end time", common.ErrInvalidInput)
	}

	if r.StartTime.Before(time.Now()) {
		return fmt.Errorf("%w: start time cannot be in the past", common.ErrInvalidInput)
	}

	// Validate status
	validStatuses := map[string]bool{statusPending: true, statusConfirmed: true, statusCancelled: true}
	if !validStatuses[r.Status] {
		return fmt.Errorf("%w: status must be one of: %s, %s, %s", common.ErrInvalidInput, statusPending, statusConfirmed, statusCancelled)
	}

	// Check for conflicts
	available, err := s.CheckReservationAvailability(ctx, r.ResourceID, r.StartTime, r.EndTime)
	if err != nil {
		return err
	}
	if !available {
		return fmt.Errorf("%w: resource is not available for the requested time", common.ErrConflict)
	}

	return s.repo.CreateReservation(ctx, r)
}

// UpdateReservation updates an existing reservation with validation
func (s *Service) UpdateReservation(ctx context.Context, r *reservation.Reservation) error {
	if r.ID == 0 {
		return fmt.Errorf("%w: reservation ID cannot be zero for update", common.ErrInvalidInput)
	}

	// Validate resource ID
	if r.ResourceID == 0 {
		return fmt.Errorf("%w: resource ID cannot be zero", common.ErrInvalidInput)
	}

	// Validate purpose
	if strings.TrimSpace(r.Purpose) == "" {
		return fmt.Errorf("%w: reservation purpose cannot be empty", common.ErrInvalidInput)
	}

	// Validate status
	validStatuses := map[string]bool{statusPending: true, statusConfirmed: true, statusCancelled: true}
	if !validStatuses[r.Status] {
		return fmt.Errorf("%w: status must be one of: %s, %s, %s", common.ErrInvalidInput, statusPending, statusConfirmed, statusCancelled)
	}

	// Validate time range
	if r.StartTime.After(r.EndTime) || r.StartTime.Equal(r.EndTime) {
		return fmt.Errorf("%w: start time must be before end time", common.ErrInvalidInput)
	}

	// Get existing reservation to check if time/resource changed
	existing, err := s.repo.ReadReservation(ctx, r.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("%w: reservation not found", common.ErrNotFound)
	}

	// Check for conflicts if time or resource changed
	if existing.ResourceID != r.ResourceID ||
		!existing.StartTime.Equal(r.StartTime) ||
		!existing.EndTime.Equal(r.EndTime) {
		available, err := s.checkReservationAvailabilityExcluding(ctx, r.ResourceID, r.StartTime, r.EndTime, r.ID)
		if err != nil {
			return err
		}
		if !available {
			return fmt.Errorf("%w: resource is not available for the requested time", common.ErrConflict)
		}
	}

	return s.repo.UpdateReservation(ctx, r)
}

// DeleteReservation removes a reservation by ID
func (s *Service) DeleteReservation(ctx context.Context, id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: reservation ID cannot be zero", common.ErrInvalidInput)
	}

	_, err := s.repo.ReadReservation(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteReservation(ctx, id)
}

// CancelReservation cancels a reservation by setting its status
func (s *Service) CancelReservation(ctx context.Context, id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: reservation ID cannot be zero", common.ErrInvalidInput)
	}

	reservation, err := s.repo.ReadReservation(ctx, id)
	if err != nil {
		return err
	}
	if reservation == nil {
		return fmt.Errorf("%w: reservation not found", common.ErrNotFound)
	}

	if reservation.Status == statusCancelled {
		return fmt.Errorf("%w: reservation is already cancelled", common.ErrConflict)
	}

	reservation.Status = statusCancelled
	return s.repo.UpdateReservation(ctx, reservation)
}

// CheckReservationAvailability checks if a resource is available for the given time period
func (s *Service) CheckReservationAvailability(ctx context.Context, resourceID uint, start, end time.Time) (bool, error) {
	if resourceID == 0 {
		return false, fmt.Errorf("%w: resource ID cannot be zero", common.ErrInvalidInput)
	}

	if start.After(end) {
		return false, fmt.Errorf("%w: start time cannot be after end time", common.ErrInvalidInput)
	}

	if start.Before(time.Now()) {
		return false, fmt.Errorf("%w: start time cannot be in the past", common.ErrInvalidInput)
	}

	// Get all reservations and check for conflicts
	// Use a large page size to get all relevant reservations
	allReservations, _, err := s.repo.ReadReservationList(ctx, common.PaginationParams{Page: 1, PageSize: 10000})
	if err != nil {
		return false, fmt.Errorf("failed to check availability: %w", err)
	}

	// Check for time conflicts with active reservations
	for _, existing := range allReservations {
		if existing.ResourceID == resourceID &&
			existing.Status != statusCancelled &&
			existing.Status != statusRejected {
			// Check if time ranges overlap
			if timeRangesOverlap(start, end, existing.StartTime, existing.EndTime) {
				return false, nil
			}
		}
	}

	return true, nil
}

// checkReservationAvailabilityExcluding checks availability excluding a specific reservation
func (s *Service) checkReservationAvailabilityExcluding(ctx context.Context, resourceID uint, start, end time.Time, excludeID uint) (bool, error) {
	// Get all reservations and check for conflicts
	// Use a large page size to get all relevant reservations
	allReservations, _, err := s.repo.ReadReservationList(ctx, common.PaginationParams{Page: 1, PageSize: 10000})
	if err != nil {
		return false, fmt.Errorf("failed to check availability: %w", err)
	}

	// Check for time conflicts with active reservations
	for _, existing := range allReservations {
		if existing.ID == excludeID {
			continue
		}

		if existing.ResourceID == resourceID &&
			existing.Status != statusCancelled &&
			existing.Status != statusRejected {
			// Check if time ranges overlap
			if timeRangesOverlap(start, end, existing.StartTime, existing.EndTime) {
				return false, nil
			}
		}
	}

	return true, nil
}

// timeRangesOverlap checks if two time ranges overlap
func timeRangesOverlap(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && start2.Before(end1)
}
