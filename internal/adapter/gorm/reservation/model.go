package reservation

import (
	"time"

	"gorm.io/gorm"
)

// GormModel represents the GORM database model for reservations
type GormModel struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	ResourceID  uint           `gorm:"not null;index" json:"resourceId"`
	UserID      uint           `gorm:"not null;index" json:"userId"`
	StartTime   time.Time      `gorm:"not null" json:"startTime"`
	EndTime     time.Time      `gorm:"not null" json:"endTime"`
	Purpose     string         `gorm:"type:varchar(255)" json:"purpose"`
	Status      string         `gorm:"type:varchar(50);default:'active'" json:"status"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for the Reservation model
func (GormModel) TableName() string {
	return "reservations"
}
