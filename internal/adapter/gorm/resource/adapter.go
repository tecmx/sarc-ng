package resource

import (
	"context"
	"encoding/json"
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

// ReadResourceList retrieves resources with pagination
func (a *GormAdapter) ReadResourceList(ctx context.Context, params domainCommon.PaginationParams) ([]resource.Resource, int64, error) {
	var models []GormModel

	query := a.db.WithContext(ctx).Model(&GormModel{})
	query = common.BuildFilteredQuery(query, params, resource.SearchableFields, false)

	query, total, err := common.CountAndPaginate(query, params.Offset(), params.Limit())
	if err != nil {
		return nil, 0, err
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]resource.Resource, len(models))
	for i, model := range models {
		entities[i] = modelToDomain(model)
	}
	return entities, total, nil
}

// ReadResource retrieves a resource by ID
func (a *GormAdapter) ReadResource(ctx context.Context, id uint) (*resource.Resource, error) {
	var model GormModel
	if err := a.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("resource not found: %w", domainCommon.ErrNotFound)
		}
		return nil, err
	}

	entity := modelToDomain(model)
	return &entity, nil
}

// CreateResource adds a new resource
func (a *GormAdapter) CreateResource(ctx context.Context, r *resource.Resource) error {
	model := domainToModel(*r)
	if err := a.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	// Update the entity with generated fields
	*r = modelToDomain(model)
	return nil
}

// UpdateResource modifies an existing resource
func (a *GormAdapter) UpdateResource(ctx context.Context, r *resource.Resource) error {
	model := domainToModel(*r)
	if err := a.db.WithContext(ctx).Save(&model).Error; err != nil {
		return err
	}

	// Update the entity with modified fields
	*r = modelToDomain(model)
	return nil
}

// DeleteResource removes a resource
func (a *GormAdapter) DeleteResource(ctx context.Context, id uint) error {
	return a.db.WithContext(ctx).Delete(&GormModel{}, id).Error
}

// Migrate ensures the resources table matches the adapter schema
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&GormModel{})
}

// domainToModel converts domain entity to GORM model
func domainToModel(entity resource.Resource) GormModel {
	var detailsJSON string
	if entity.Details != nil {
		if bytes, err := json.Marshal(entity.Details); err == nil {
			detailsJSON = string(bytes)
		}
	}

	return GormModel{
		ID:         entity.ID,
		Name:       entity.Name,
		Type:       entity.Type,
		Capacity:   entity.Capacity,
		Location:   entity.Location,
		BuildingID: entity.BuildingID,
		Details:    detailsJSON,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
		DeletedAt:  common.ConvertTimeToGormDeletedAt(entity.DeletedAt),
	}
}

// modelToDomain converts GORM model to domain entity
func modelToDomain(model GormModel) resource.Resource {
	var details map[string]interface{}
	if model.Details != "" {
		json.Unmarshal([]byte(model.Details), &details)
	}

	return resource.Resource{
		ID:         model.ID,
		Name:       model.Name,
		Type:       model.Type,
		Capacity:   model.Capacity,
		Location:   model.Location,
		BuildingID: model.BuildingID,
		Details:    details,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
		DeletedAt:  common.ConvertGormDeletedAtToTime(model.DeletedAt),
	}
}
