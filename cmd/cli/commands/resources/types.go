package resources

import "time"

// ResourceRequest represents a resource creation/update request
type ResourceRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	IsAvailable bool   `json:"isAvailable"`
}

// Resource represents a resource response
type Resource struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	IsAvailable bool      `json:"isAvailable"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
