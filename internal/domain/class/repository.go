package class

import (
	"context"
	"sarc-ng/internal/domain/common"
)

// Repository defines the data access operations for classes
// All methods are explicitly named with the Class entity
type Repository interface {
	ReadClassList(ctx context.Context, params common.PaginationParams) ([]Class, int64, error)
	ReadClass(ctx context.Context, id uint) (*Class, error)
	CreateClass(ctx context.Context, class *Class) error
	UpdateClass(ctx context.Context, class *Class) error
	DeleteClass(ctx context.Context, id uint) error
}
