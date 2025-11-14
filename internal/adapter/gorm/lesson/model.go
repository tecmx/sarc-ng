package lesson

import (
	"time"

	"gorm.io/gorm"
)

// GormModel represents the GORM database model for lessons
type GormModel struct {
	ID            uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ClassID       uint           `gorm:"not null;index" json:"classId"`
	ResourceID    uint           `gorm:"not null;index" json:"resourceId"`
	StartTime     time.Time      `gorm:"type:datetime;not null" json:"startTime"`
	EndTime       time.Time      `gorm:"type:datetime;not null" json:"endTime"`
	RepeatPattern string         `gorm:"type:varchar(50);not null" json:"repeatPattern"`
	RepeatUntil   *time.Time     `gorm:"type:datetime" json:"repeatUntil"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for the Lesson model
func (GormModel) TableName() string {
	return "lessons"
}
