package reservations

import "time"

// ReservationRequest represents a reservation creation/update request
type ReservationRequest struct {
	ResourceID uint      `json:"resourceId"`
	UserID     uint      `json:"userId"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	Status     string    `json:"status,omitempty"`
}

// Reservation represents a reservation response
type Reservation struct {
	ID         uint      `json:"id"`
	ResourceID uint      `json:"resourceId"`
	UserID     uint      `json:"userId"`
	StartTime  time.Time `json:"startTime"`
	EndTime    time.Time `json:"endTime"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
