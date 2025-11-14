package class

import (
	"time"
)

// CreateClassDTO represents the data needed to create a class
type CreateClassDTO struct {
	Name           string  `json:"name" validate:"required"`
	Code           string  `json:"code" validate:"required"`
	InstructorName string  `json:"instructorName" validate:"required"`
	Semester       *string `json:"semester"`
	Capacity       int     `json:"capacity" validate:"min=1"`
}

// UpdateClassDTO represents the data needed to update a class
type UpdateClassDTO struct {
	Name           string  `json:"name" validate:"required"`
	Code           string  `json:"code" validate:"required"`
	InstructorName string  `json:"instructorName" validate:"required"`
	Semester       *string `json:"semester"`
	Capacity       int     `json:"capacity" validate:"min=1"`
}

// ClassDTO represents class data for application operations
type ClassDTO struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	Code           string    `json:"code"`
	InstructorName string    `json:"instructorName"`
	Semester       *string   `json:"semester,omitempty"`
	Capacity       int       `json:"capacity"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
