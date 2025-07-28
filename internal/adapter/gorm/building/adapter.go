package building

import (
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

// ReadBuildingList retrieves all buildings
func (a *GormAdapter) ReadBuildingList() ([]building.Building, error) {
	var models []GormModel
	if err := a.db.Find(&models).Error; err != nil {
		return nil, err
	}

	entities := make([]building.Building, len(models))
	for i, model := range models {
		entities[i] = modelToDomain(model)
	}
	return entities, nil
}

// ReadBuilding retrieves a building by ID
func (a *GormAdapter) ReadBuilding(id uint) (*building.Building, error) {
	var model GormModel
	if err := a.db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("building not found: %w", domainCommon.ErrNotFound)
		}
		return nil, err
	}

	entity := modelToDomain(model)
	return &entity, nil
}

// CreateBuilding adds a new building
func (a *GormAdapter) CreateBuilding(b *building.Building) error {
	model := domainToModel(*b)
	if err := a.db.Create(&model).Error; err != nil {
		return err
	}

	// Update the entity with generated fields
	*b = modelToDomain(model)
	return nil
}

// UpdateBuilding modifies an existing building
func (a *GormAdapter) UpdateBuilding(b *building.Building) error {
	model := domainToModel(*b)
	if err := a.db.Save(&model).Error; err != nil {
		return err
	}

	// Update the entity with modified fields
	*b = modelToDomain(model)
	return nil
}

// DeleteBuilding removes a building
func (a *GormAdapter) DeleteBuilding(id uint) error {
	return a.db.Delete(&GormModel{}, id).Error
}

// domainToModel converts domain entity to GORM model
func domainToModel(entity building.Building) GormModel {
	return GormModel{
		ID:        entity.ID,
		Name:      entity.Name,
		Code:      entity.Code,
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
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: common.ConvertGormDeletedAtToTime(model.DeletedAt),
	}
}
