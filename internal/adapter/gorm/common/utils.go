package common

import (
	"time"

	"gorm.io/gorm"
)

// ConvertTimeToGormDeletedAt converts a *time.Time to gorm.DeletedAt
// This eliminates the duplicate helper function across all adapters
func ConvertTimeToGormDeletedAt(t *time.Time) gorm.DeletedAt {
	if t == nil {
		return gorm.DeletedAt{}
	}
	return gorm.DeletedAt{Time: *t, Valid: true}
}

// ConvertGormDeletedAtToTime converts gorm.DeletedAt to *time.Time
// This eliminates the duplicate helper function across all adapters
func ConvertGormDeletedAtToTime(gd gorm.DeletedAt) *time.Time {
	if !gd.Valid {
		return nil
	}
	return &gd.Time
}
