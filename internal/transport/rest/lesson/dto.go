package lesson

import (
	"time"
)

// CreateLessonDTO represents the data needed to create a lesson
type CreateLessonDTO struct {
	Title     string    `json:"title" validate:"required"`
	Duration  int       `json:"duration" validate:"min=1"`
	StartTime time.Time `json:"startTime,omitempty"`
}

// UpdateLessonDTO represents the data needed to update a lesson
type UpdateLessonDTO struct {
	Title     string    `json:"title" validate:"required"`
	Duration  int       `json:"duration" validate:"min=1"`
	StartTime time.Time `json:"startTime,omitempty"`
}

// LessonDTO represents lesson data for application operations
type LessonDTO struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Duration  int       `json:"duration"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
