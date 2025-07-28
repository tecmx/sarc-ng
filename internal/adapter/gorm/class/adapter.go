package class

import (
	"fmt"
	"sarc-ng/internal/adapter/gorm/common"
	"sarc-ng/internal/domain/class"
	domainCommon "sarc-ng/internal/domain/common"

	"gorm.io/gorm"
)

// GormAdapter implements class.Repository using GORM
type GormAdapter struct {
	db *gorm.DB
}

// Compile-time verification that GormAdapter implements class.Repository
var _ class.Repository = (*GormAdapter)(nil)

// NewGormAdapter creates a new class GORM adapter
func NewGormAdapter(db *gorm.DB) *GormAdapter {
	return &GormAdapter{
		db: db,
	}
}

// ReadClassList retrieves all classes
func (a *GormAdapter) ReadClassList() ([]class.Class, error) {
	var models []GormModel
	if err := a.db.Find(&models).Error; err != nil {
		return nil, err
	}

	entities := make([]class.Class, len(models))
	for i, model := range models {
		entities[i] = modelToDomain(model)
	}
	return entities, nil
}

// ReadClass retrieves a class by ID
func (a *GormAdapter) ReadClass(id uint) (*class.Class, error) {
	var model GormModel
	if err := a.db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("class not found: %w", domainCommon.ErrNotFound)
		}
		return nil, err
	}

	entity := modelToDomain(model)
	return &entity, nil
}

// CreateClass adds a new class
func (a *GormAdapter) CreateClass(c *class.Class) error {
	model := domainToModel(*c)
	if err := a.db.Create(&model).Error; err != nil {
		return err
	}

	// Update the entity with generated fields
	*c = modelToDomain(model)
	return nil
}

// UpdateClass modifies an existing class
func (a *GormAdapter) UpdateClass(c *class.Class) error {
	model := domainToModel(*c)
	if err := a.db.Save(&model).Error; err != nil {
		return err
	}

	// Update the entity with modified fields
	*c = modelToDomain(model)
	return nil
}

// DeleteClass removes a class
func (a *GormAdapter) DeleteClass(id uint) error {
	return a.db.Delete(&GormModel{}, id).Error
}

// domainToModel converts domain entity to GORM model
func domainToModel(entity class.Class) GormModel {
	return GormModel{
		ID:        entity.ID,
		Name:      entity.Name,
		Capacity:  entity.Capacity,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: common.ConvertTimeToGormDeletedAt(entity.DeletedAt),
	}
}

// modelToDomain converts GORM model to domain entity
func modelToDomain(model GormModel) class.Class {
	return class.Class{
		ID:        model.ID,
		Name:      model.Name,
		Capacity:  model.Capacity,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: common.ConvertGormDeletedAtToTime(model.DeletedAt),
	}
}
