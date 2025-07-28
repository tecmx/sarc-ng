package building

import (
	"time"
)

// Building represents a physical building in the system
type Building struct {
	ID        uint
	Name      string
	Code      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
