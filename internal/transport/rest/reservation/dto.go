package reservation

import (
	"time"
)

// CreateReservationDTO represents the data needed to create a reservation
type CreateReservationDTO struct {
	ResourceID uint      `json:"resourceId" validate:"required"`
	StartTime  time.Time `json:"startTime" validate:"required"`
	EndTime    time.Time `json:"endTime" validate:"required"`
	Purpose    string    `json:"purpose" validate:"required"`
	Status     string    `json:"status" validate:"required,oneof=pending confirmed cancelled"`
}

// UpdateReservationDTO represents the data needed to update a reservation
type UpdateReservationDTO struct {
	ResourceID uint      `json:"resourceId" validate:"required"`
	StartTime  time.Time `json:"startTime" validate:"required"`
	EndTime    time.Time `json:"endTime" validate:"required"`
	Purpose    string    `json:"purpose" validate:"required"`
	Status     string    `json:"status" validate:"required,oneof=pending confirmed cancelled"`
}

// ReservationDTO represents reservation data for application operations
type ReservationDTO struct {
	ID         uint      `json:"id"`
	ResourceID uint      `json:"resourceId"`
	UserID     uint      `json:"userId"`
	UserName   string    `json:"userName"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	Purpose    string    `json:"purpose"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
