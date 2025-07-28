package building

import (
	"time"
)

// CreateBuildingDTO represents the data needed to create a building
type CreateBuildingDTO struct {
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required"`
}

// UpdateBuildingDTO represents the data needed to update a building
type UpdateBuildingDTO struct {
	ID   uint   `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	Code string `json:"code" validate:"required"`
}

// BuildingDTO represents building data for application operations
type BuildingDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
