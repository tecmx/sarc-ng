package building

import (
	"context"
	"fmt"
	"sarc-ng/internal/adapter/gorm/common"
	"sarc-ng/internal/domain/building"
	domainCommon "sarc-ng/internal/domain/common"

	"gorm.io/gorm"
)

// GormAdapter implements building.Repository using GORM
type GormAdapter struct {
	db *gorm.DB
}

// Compile-time verification that GormAdapter implements building.Repository
var _ building.Repository = (*GormAdapter)(nil)

// NewGormAdapter creates a new building GORM adapter
func NewGormAdapter(db *gorm.DB) *GormAdapter {
	return &GormAdapter{
		db: db,
	}
}

// ReadBuildingList retrieves buildings with pagination
func (a *GormAdapter) ReadBuildingList(ctx context.Context, params domainCommon.PaginationParams) ([]building.Building, int64, error) {
	var models []GormModel

	query := a.db.WithContext(ctx).Model(&GormModel{})
	query = common.BuildFilteredQuery(query, params, building.SearchableFields, false)

	query, total, err := common.CountAndPaginate(query, params.Offset(), params.Limit())
	if err != nil {
		return nil, 0, err
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]building.Building, len(models))
	for i, model := range models {
		entities[i] = modelToDomain(model)
	}
	return entities, total, nil
}

// ReadBuilding retrieves a building by ID
func (a *GormAdapter) ReadBuilding(ctx context.Context, id uint) (*building.Building, error) {
	var model GormModel
	if err := a.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("building not found: %w", domainCommon.ErrNotFound)
		}
		return nil, err
	}

	entity := modelToDomain(model)
	return &entity, nil
}

// FindBuildingByCode retrieves a building by code
func (a *GormAdapter) FindBuildingByCode(ctx context.Context, code string) (*building.Building, error) {
	var model GormModel
	if err := a.db.WithContext(ctx).Where("code = ?", code).First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	entity := modelToDomain(model)
	return &entity, nil
}

// CreateBuilding adds a new building
func (a *GormAdapter) CreateBuilding(ctx context.Context, b *building.Building) error {
	model := domainToModel(*b)
	if err := a.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	// Update the entity with generated fields
	*b = modelToDomain(model)
	return nil
}

// UpdateBuilding modifies an existing building
func (a *GormAdapter) UpdateBuilding(ctx context.Context, b *building.Building) error {
	model := domainToModel(*b)
	if err := a.db.WithContext(ctx).Save(&model).Error; err != nil {
		return err
	}

	// Update the entity with modified fields
	*b = modelToDomain(model)
	return nil
}

// DeleteBuilding removes a building
func (a *GormAdapter) DeleteBuilding(ctx context.Context, id uint) error {
	return a.db.WithContext(ctx).Delete(&GormModel{}, id).Error
}

// Migrate ensures the buildings table matches the adapter schema
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&GormModel{})
}

// domainToModel converts domain entity to GORM model
func domainToModel(entity building.Building) GormModel {
	return GormModel{
		ID:        entity.ID,
		Name:      entity.Name,
		Code:      entity.Code,
		Address:   entity.Address,
		Floors:    entity.Floors,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: common.ConvertTimeToGormDeletedAt(entity.DeletedAt),
	}
}

// modelToDomain converts GORM model to domain entity
func modelToDomain(model GormModel) building.Building {
	return building.Building{
		ID:        model.ID,
		Name:      model.Name,
		Code:      model.Code,
		Address:   model.Address,
		Floors:    model.Floors,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: common.ConvertGormDeletedAtToTime(model.DeletedAt),
	}
}
