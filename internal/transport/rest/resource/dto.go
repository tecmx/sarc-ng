package resource

import (
	"time"
)

// CreateResourceDTO represents the data needed to create a resource
type CreateResourceDTO struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required"`
	Description string `json:"description"`
	Location    string `json:"location"`
	IsAvailable bool   `json:"isAvailable"`
}

// UpdateResourceDTO represents the data needed to update a resource
type UpdateResourceDTO struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required"`
	Description string `json:"description"`
	Location    string `json:"location"`
	IsAvailable bool   `json:"isAvailable"`
}

// ResourceDTO represents resource data for application operations
type ResourceDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	IsAvailable bool      `json:"isAvailable"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
