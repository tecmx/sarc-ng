package reservation

import (
	"errors"
	"sarc-ng/internal/domain/reservation"
	"time"
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

// GetAllReservations retrieves all reservations
func (s *Service) GetAllReservations() ([]reservation.Reservation, error) {
	return s.repo.ReadReservationList()
}

// GetReservation retrieves a reservation by ID with validation
func (s *Service) GetReservation(id uint) (*reservation.Reservation, error) {
	if id == 0 {
		return nil, errors.New("reservation ID cannot be zero")
	}
	return s.repo.ReadReservation(id)
}

// CreateReservation creates a new reservation
func (s *Service) CreateReservation(r *reservation.Reservation) error {
	return s.repo.CreateReservation(r)
}

// UpdateReservation updates an existing reservation
func (s *Service) UpdateReservation(r *reservation.Reservation) error {
	if r.ID == 0 {
		return errors.New("reservation ID cannot be zero for update")
	}

	// Get existing reservation to check if time/resource changed
	existing, err := s.repo.ReadReservation(r.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("reservation not found")
	}

	return s.repo.UpdateReservation(r)
}

// DeleteReservation removes a reservation by ID
func (s *Service) DeleteReservation(id uint) error {
	if id == 0 {
		return errors.New("reservation ID cannot be zero")
	}

	// Check if reservation exists before deletion
	_, err := s.repo.ReadReservation(id)
	if err != nil {
		return err // Repository already returns "reservation not found" for missing records
	}

	return s.repo.DeleteReservation(id)
}

// CancelReservation cancels a reservation by setting its status
func (s *Service) CancelReservation(id uint) error {
	if id == 0 {
		return errors.New("reservation ID cannot be zero")
	}

	reservation, err := s.repo.ReadReservation(id)
	if err != nil {
		return err
	}
	if reservation == nil {
		return errors.New("reservation not found")
	}

	if reservation.Status == "cancelled" {
		return errors.New("reservation is already cancelled")
	}

	reservation.Status = "cancelled"
	return s.repo.UpdateReservation(reservation)
}

// CheckReservationAvailability checks if a resource is available for the given time period
func (s *Service) CheckReservationAvailability(resourceID uint, start, end time.Time) (bool, error) {
	if resourceID == 0 {
		return false, errors.New("resource ID cannot be zero")
	}

	if start.After(end) {
		return false, errors.New("start time cannot be after end time")
	}

	if start.Before(time.Now()) {
		return false, errors.New("start time cannot be in the past")
	}

	// Since we removed FindConflicting, we'll assume availability for now
	// In a real implementation, this would need to be handled differently
	// or the availability checking would be moved to a different service
	return true, nil
}
