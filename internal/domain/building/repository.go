package building

import (
	"context"
	"sarc-ng/internal/domain/common"
)

// Repository defines the data access operations for buildings
// All methods are explicitly named with the Building entity
type Repository interface {
	ReadBuildingList(ctx context.Context, params common.PaginationParams) ([]Building, int64, error)
	ReadBuilding(ctx context.Context, id uint) (*Building, error)
	FindBuildingByCode(ctx context.Context, code string) (*Building, error)
	CreateBuilding(ctx context.Context, building *Building) error
	UpdateBuilding(ctx context.Context, building *Building) error
	DeleteBuilding(ctx context.Context, id uint) error
}
