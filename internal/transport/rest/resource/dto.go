package resource

import (
	"time"
)

// CreateResourceDTO represents the data needed to create a resource
type CreateResourceDTO struct {
	Name       string                 `json:"name" validate:"required"`
	Type       string                 `json:"type" validate:"required,oneof=classroom equipment lab"`
	Capacity   *int                   `json:"capacity"`
	Location   string                 `json:"location"`
	BuildingID *uint                  `json:"buildingId"`
	Details    map[string]interface{} `json:"details"`
}

// UpdateResourceDTO represents the data needed to update a resource
type UpdateResourceDTO struct {
	Name       string                 `json:"name" validate:"required"`
	Type       string                 `json:"type" validate:"required,oneof=classroom equipment lab"`
	Capacity   *int                   `json:"capacity"`
	Location   string                 `json:"location"`
	BuildingID *uint                  `json:"buildingId"`
	Details    map[string]interface{} `json:"details"`
}

// ResourceDTO represents resource data for application operations
type ResourceDTO struct {
	ID         uint                   `json:"id"`
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Capacity   *int                   `json:"capacity,omitempty"`
	Location   string                 `json:"location"`
	BuildingID *uint                  `json:"buildingId,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	CreatedAt  time.Time              `json:"createdAt"`
	UpdatedAt  time.Time              `json:"updatedAt"`
}
