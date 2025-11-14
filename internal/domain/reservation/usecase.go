package reservation

import (
	"context"
	"sarc-ng/internal/domain/common"
	"time"
)

// Usecase defines the business logic operations for reservation management
type Usecase interface {
	GetAllReservations(ctx context.Context, params common.PaginationParams) ([]Reservation, common.PaginationResult, error)
	GetReservation(ctx context.Context, id uint) (*Reservation, error)
	CreateReservation(ctx context.Context, reservation *Reservation) error
	UpdateReservation(ctx context.Context, reservation *Reservation) error
	DeleteReservation(ctx context.Context, id uint) error
	CancelReservation(ctx context.Context, id uint) error
	CheckReservationAvailability(ctx context.Context, resourceID uint, start, end time.Time) (bool, error)
}
