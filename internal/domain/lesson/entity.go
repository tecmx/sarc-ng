package lesson

import (
	"time"
)

// Lesson represents a teaching session in the system
type Lesson struct {
	ID          uint
	Title       string
	Duration    int // Duration in minutes
	Description string
	StartTime   time.Time
	EndTime     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
