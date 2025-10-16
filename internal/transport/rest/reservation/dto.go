package reservation

import (
	"time"
)

// CreateReservationDTO represents the data needed to create a reservation
type CreateReservationDTO struct {
	ResourceID  uint      `json:"resourceId" validate:"required"`
	UserID      uint      `json:"userId" validate:"required"`
	StartTime   time.Time `json:"startTime" validate:"required"`
	EndTime     time.Time `json:"endTime" validate:"required"`
	Purpose     string    `json:"purpose" validate:"required"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
}

// UpdateReservationDTO represents the data needed to update a reservation
type UpdateReservationDTO struct {
	ResourceID  uint      `json:"resourceId" validate:"required"`
	UserID      uint      `json:"userId" validate:"required"`
	StartTime   time.Time `json:"startTime" validate:"required"`
	EndTime     time.Time `json:"endTime" validate:"required"`
	Purpose     string    `json:"purpose" validate:"required"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
}

// ReservationDTO represents reservation data for application operations
type ReservationDTO struct {
	ID          uint      `json:"id"`
	ResourceID  uint      `json:"resourceId"`
	UserID      uint      `json:"userId"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Purpose     string    `json:"purpose"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
