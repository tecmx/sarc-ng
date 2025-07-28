package reservation

import (
	"time"
)

// CreateReservationDTO represents the data needed to create a reservation
type CreateReservationDTO struct {
	ResourceID uint      `json:"resourceId" validate:"required"`
	UserID     uint      `json:"userId" validate:"required"`
	StartTime  time.Time `json:"startTime" validate:"required"`
	EndTime    time.Time `json:"endTime" validate:"required"`
	Status     string    `json:"status"`
}

// UpdateReservationDTO represents the data needed to update a reservation
type UpdateReservationDTO struct {
	ID         uint      `json:"id" validate:"required"`
	ResourceID uint      `json:"resourceId" validate:"required"`
	UserID     uint      `json:"userId" validate:"required"`
	StartTime  time.Time `json:"startTime" validate:"required"`
	EndTime    time.Time `json:"endTime" validate:"required"`
	Status     string    `json:"status"`
}

// ReservationDTO represents reservation data for application operations
type ReservationDTO struct {
	ID         uint      `json:"id"`
	ResourceID uint      `json:"resourceId"`
	UserID     uint      `json:"userId"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
