package class

import (
	"time"
)

// Class represents a classroom or space in the system
type Class struct {
	ID        uint
	Name      string
	Capacity  int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
