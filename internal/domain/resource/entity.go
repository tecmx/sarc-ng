package resource

import (
	"time"
)

// Resource represents a bookable resource in the system
type Resource struct {
	ID          uint
	Name        string
	Type        string
	Description string
	IsAvailable bool
	Location    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
