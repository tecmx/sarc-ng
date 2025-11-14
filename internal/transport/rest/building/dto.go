package building

import (
	"time"
)

// CreateBuildingDTO represents the data needed to create a building
type CreateBuildingDTO struct {
	Name    string  `json:"name" validate:"required"`
	Code    string  `json:"code" validate:"required"`
	Address *string `json:"address"`
	Floors  *int    `json:"floors"`
}

// UpdateBuildingDTO represents the data needed to update a building
type UpdateBuildingDTO struct {
	Name    string  `json:"name" validate:"required"`
	Code    string  `json:"code" validate:"required"`
	Address *string `json:"address"`
	Floors  *int    `json:"floors"`
}

// BuildingDTO represents building data for application operations
type BuildingDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Address   *string   `json:"address,omitempty"`
	Floors    *int      `json:"floors,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
