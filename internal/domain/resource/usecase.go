package resource

import (
	"context"
	"sarc-ng/internal/domain/common"
)

// Usecase defines the business logic operations for resource management
type Usecase interface {
	GetAllResources(ctx context.Context, params common.PaginationParams) ([]Resource, common.PaginationResult, error)
	GetResource(ctx context.Context, id uint) (*Resource, error)
	CreateResource(ctx context.Context, resource *Resource) error
	UpdateResource(ctx context.Context, resource *Resource) error
	DeleteResource(ctx context.Context, id uint) error
}
