package lesson

import (
	"time"
)

// Lesson represents a scheduled class session in the system
type Lesson struct {
	ID            uint
	ClassID       uint
	ResourceID    uint
	StartTime     time.Time
	EndTime       time.Time
	RepeatPattern string
	RepeatUntil   *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}
