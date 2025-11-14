package class

import (
	"context"
	"sarc-ng/internal/domain/common"
)

// Usecase defines the business logic operations for class management
type Usecase interface {
	GetAllClasses(ctx context.Context, params common.PaginationParams) ([]Class, common.PaginationResult, error)
	GetClass(ctx context.Context, id uint) (*Class, error)
	CreateClass(ctx context.Context, class *Class) error
	UpdateClass(ctx context.Context, class *Class) error
	DeleteClass(ctx context.Context, id uint) error
}
