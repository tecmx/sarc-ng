package resource

import (
	"context"
	"sarc-ng/internal/domain/common"
)

// Repository defines the data access operations for resources
// All methods are explicitly named with the Resource entity
type Repository interface {
	ReadResourceList(ctx context.Context, params common.PaginationParams) ([]Resource, int64, error)
	ReadResource(ctx context.Context, id uint) (*Resource, error)
	CreateResource(ctx context.Context, resource *Resource) error
	UpdateResource(ctx context.Context, resource *Resource) error
	DeleteResource(ctx context.Context, id uint) error
}
