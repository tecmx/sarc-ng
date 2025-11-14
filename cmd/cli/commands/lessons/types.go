package lessons

import "time"

// LessonRequest represents a lesson creation/update request
type LessonRequest struct {
	Title     string    `json:"title"`
	Duration  int       `json:"duration"`
	StartTime time.Time `json:"startTime,omitempty"`
}

// Lesson represents a lesson response
type Lesson struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Duration  int       `json:"duration"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
