package lesson

import (
	"time"
)

// CreateLessonDTO represents the data needed to create a lesson
type CreateLessonDTO struct {
	ClassID       uint       `json:"classId" validate:"required"`
	ResourceID    uint       `json:"resourceId" validate:"required"`
	StartTime     time.Time  `json:"startTime" validate:"required"`
	EndTime       time.Time  `json:"endTime" validate:"required"`
	RepeatPattern string     `json:"repeatPattern" validate:"required,oneof=none daily weekly biweekly monthly"`
	RepeatUntil   *time.Time `json:"repeatUntil"`
}

// UpdateLessonDTO represents the data needed to update a lesson
type UpdateLessonDTO struct {
	ClassID       uint       `json:"classId" validate:"required"`
	ResourceID    uint       `json:"resourceId" validate:"required"`
	StartTime     time.Time  `json:"startTime" validate:"required"`
	EndTime       time.Time  `json:"endTime" validate:"required"`
	RepeatPattern string     `json:"repeatPattern" validate:"required,oneof=none daily weekly biweekly monthly"`
	RepeatUntil   *time.Time `json:"repeatUntil"`
}

// LessonDTO represents lesson data for application operations
type LessonDTO struct {
	ID            uint       `json:"id"`
	ClassID       uint       `json:"classId"`
	ResourceID    uint       `json:"resourceId"`
	StartTime     time.Time  `json:"startTime"`
	EndTime       time.Time  `json:"endTime"`
	RepeatPattern string     `json:"repeatPattern"`
	RepeatUntil   *time.Time `json:"repeatUntil,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}
