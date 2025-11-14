package reservation

import (
	"time"
)

// Reservation represents a booking in the system
type Reservation struct {
	ID         uint
	ResourceID uint
	UserID     uint
	UserName   string
	StartTime  time.Time
	EndTime    time.Time
	Purpose    string
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
