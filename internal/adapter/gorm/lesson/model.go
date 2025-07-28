package lesson

import (
	"time"

	"gorm.io/gorm"
)

// GormModel represents the GORM database model for lessons
type GormModel struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Duration    int            `gorm:"not null;default:60" json:"duration"` // Duration in minutes
	Description string         `gorm:"type:text" json:"description"`
	StartTime   time.Time      `gorm:"type:datetime" json:"startTime"`
	EndTime     time.Time      `gorm:"type:datetime" json:"endTime"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for the Lesson model
func (GormModel) TableName() string {
	return "lessons"
}
