package class

import (
	"context"
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

// ReadClassList retrieves classes with pagination
func (a *GormAdapter) ReadClassList(ctx context.Context, params domainCommon.PaginationParams) ([]class.Class, int64, error) {
	var models []GormModel

	query := a.db.WithContext(ctx).Model(&GormModel{})
	query = common.BuildFilteredQuery(query, params, class.SearchableFields, false)

	query, total, err := common.CountAndPaginate(query, params.Offset(), params.Limit())
	if err != nil {
		return nil, 0, err
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]class.Class, len(models))
	for i, model := range models {
		entities[i] = modelToDomain(model)
	}
	return entities, total, nil
}

// ReadClass retrieves a class by ID
func (a *GormAdapter) ReadClass(ctx context.Context, id uint) (*class.Class, error) {
	var model GormModel
	if err := a.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("class not found: %w", domainCommon.ErrNotFound)
		}
		return nil, err
	}

	entity := modelToDomain(model)
	return &entity, nil
}

// CreateClass adds a new class
func (a *GormAdapter) CreateClass(ctx context.Context, c *class.Class) error {
	model := domainToModel(*c)
	if err := a.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	// Update the entity with generated fields
	*c = modelToDomain(model)
	return nil
}

// UpdateClass modifies an existing class
func (a *GormAdapter) UpdateClass(ctx context.Context, c *class.Class) error {
	model := domainToModel(*c)
	if err := a.db.WithContext(ctx).Save(&model).Error; err != nil {
		return err
	}

	// Update the entity with modified fields
	*c = modelToDomain(model)
	return nil
}

// DeleteClass removes a class
func (a *GormAdapter) DeleteClass(ctx context.Context, id uint) error {
	return a.db.WithContext(ctx).Delete(&GormModel{}, id).Error
}

// Migrate ensures the classes table matches the adapter schema
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&GormModel{})
}

// domainToModel converts domain entity to GORM model
func domainToModel(entity class.Class) GormModel {
	return GormModel{
		ID:             entity.ID,
		Name:           entity.Name,
		Code:           entity.Code,
		InstructorName: entity.InstructorName,
		Semester:       entity.Semester,
		Capacity:       entity.Capacity,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		DeletedAt:      common.ConvertTimeToGormDeletedAt(entity.DeletedAt),
	}
}

// modelToDomain converts GORM model to domain entity
func modelToDomain(model GormModel) class.Class {
	return class.Class{
		ID:             model.ID,
		Name:           model.Name,
		Code:           model.Code,
		InstructorName: model.InstructorName,
		Semester:       model.Semester,
		Capacity:       model.Capacity,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
		DeletedAt:      common.ConvertGormDeletedAtToTime(model.DeletedAt),
	}
}
