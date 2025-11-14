package reservation

import (
	"context"
	"sarc-ng/internal/domain/common"
)

// Repository defines the data access operations for reservations
// All methods are explicitly named with the Reservation entity
type Repository interface {
	ReadReservationList(ctx context.Context, params common.PaginationParams) ([]Reservation, int64, error)
	ReadReservation(ctx context.Context, id uint) (*Reservation, error)
	CreateReservation(ctx context.Context, reservation *Reservation) error
	UpdateReservation(ctx context.Context, reservation *Reservation) error
	DeleteReservation(ctx context.Context, id uint) error
}
