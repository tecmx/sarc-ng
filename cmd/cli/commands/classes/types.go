package classes

import "time"

// ClassRequest represents a class creation/update request
type ClassRequest struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

// Class represents a class response
type Class struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Capacity  int       `json:"capacity"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
