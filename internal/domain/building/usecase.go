package building

import (
	"context"
	"sarc-ng/internal/domain/common"
)

// Usecase defines the business logic operations for building management
type Usecase interface {
	GetAllBuildings(ctx context.Context, params common.PaginationParams) ([]Building, common.PaginationResult, error)
	GetBuilding(ctx context.Context, id uint) (*Building, error)
	CreateBuilding(ctx context.Context, building *Building) error
	UpdateBuilding(ctx context.Context, building *Building) error
	DeleteBuilding(ctx context.Context, id uint) error
}
