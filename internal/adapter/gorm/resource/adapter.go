package resource

import (
	"fmt"
	"sarc-ng/internal/adapter/gorm/common"
	domainCommon "sarc-ng/internal/domain/common"
	"sarc-ng/internal/domain/resource"

	"gorm.io/gorm"
)

// GormAdapter implements resource.Repository using GORM
type GormAdapter struct {
	db *gorm.DB
}

// Compile-time verification that GormAdapter implements resource.Repository
var _ resource.Repository = (*GormAdapter)(nil)

// NewGormAdapter creates a new resource GORM adapter
func NewGormAdapter(db *gorm.DB) *GormAdapter {
	return &GormAdapter{
		db: db,
	}
}

// ReadResourceList retrieves all resources
func (a *GormAdapter) ReadResourceList() ([]resource.Resource, error) {
	var models []GormModel
	if err := a.db.Find(&models).Error; err != nil {
		return nil, err
	}

	entities := make([]resource.Resource, len(models))
	for i, model := range models {
		entities[i] = modelToDomain(model)
	}
	return entities, nil
}

// ReadResource retrieves a resource by ID
func (a *GormAdapter) ReadResource(id uint) (*resource.Resource, error) {
	var model GormModel
	if err := a.db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("resource not found: %w", domainCommon.ErrNotFound)
		}
		return nil, err
	}

	entity := modelToDomain(model)
	return &entity, nil
}

// CreateResource adds a new resource
func (a *GormAdapter) CreateResource(r *resource.Resource) error {
	model := domainToModel(*r)
	if err := a.db.Create(&model).Error; err != nil {
		return err
	}

	// Update the entity with generated fields
	*r = modelToDomain(model)
	return nil
}

// UpdateResource modifies an existing resource
func (a *GormAdapter) UpdateResource(r *resource.Resource) error {
	model := domainToModel(*r)
	if err := a.db.Save(&model).Error; err != nil {
		return err
	}

	// Update the entity with modified fields
	*r = modelToDomain(model)
	return nil
}

// DeleteResource removes a resource
func (a *GormAdapter) DeleteResource(id uint) error {
	return a.db.Delete(&GormModel{}, id).Error
}

// domainToModel converts domain entity to GORM model
func domainToModel(entity resource.Resource) GormModel {
	return GormModel{
		ID:          entity.ID,
		Name:        entity.Name,
		Type:        entity.Type,
		Description: entity.Description,
		IsAvailable: entity.IsAvailable,
		Location:    entity.Location,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		DeletedAt:   common.ConvertTimeToGormDeletedAt(entity.DeletedAt),
	}
}

// modelToDomain converts GORM model to domain entity
func modelToDomain(model GormModel) resource.Resource {
	return resource.Resource{
		ID:          model.ID,
		Name:        model.Name,
		Type:        model.Type,
		Description: model.Description,
		IsAvailable: model.IsAvailable,
		Location:    model.Location,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		DeletedAt:   common.ConvertGormDeletedAtToTime(model.DeletedAt),
	}
}
