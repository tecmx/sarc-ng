package resource

import (
	"time"

	"gorm.io/gorm"
)

// GormModel represents the GORM database model for resources
type GormModel struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Type        string         `gorm:"type:varchar(100);not null" json:"type"`
	Description string         `gorm:"type:text" json:"description"`
	IsAvailable bool           `gorm:"default:true" json:"isAvailable"`
	Location    string         `gorm:"type:varchar(255)" json:"location"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for the Resource model
func (GormModel) TableName() string {
	return "resources"
}
