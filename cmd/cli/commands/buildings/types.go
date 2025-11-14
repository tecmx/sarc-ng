package buildings

import "time"

// BuildingRequest represents a building creation/update request
type BuildingRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// Building represents a building response
type Building struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
