package resource

import (
	"time"
)

// Resource represents a bookable resource in the system
type Resource struct {
	ID         uint
	Name       string
	Type       string
	Capacity   *int
	Location   string
	BuildingID *uint
	Details    map[string]interface{}
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}
