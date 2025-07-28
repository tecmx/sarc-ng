package class

import (
	"time"

	"gorm.io/gorm"
)

// GormModel represents the GORM database model for classes
type GormModel struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Capacity  int            `gorm:"not null;default:0" json:"capacity"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for the Class model
func (GormModel) TableName() string {
	return "classes"
}
