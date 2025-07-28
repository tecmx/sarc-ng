package reservation

import "time"

// Usecase defines the business logic operations for reservation management
type Usecase interface {
	GetAllReservations() ([]Reservation, error)
	GetReservation(id uint) (*Reservation, error)
	CreateReservation(reservation *Reservation) error
	UpdateReservation(reservation *Reservation) error
	DeleteReservation(id uint) error
	CancelReservation(id uint) error
	CheckReservationAvailability(resourceID uint, start, end time.Time) (bool, error)
}
