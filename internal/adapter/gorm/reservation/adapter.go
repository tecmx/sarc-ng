package reservation

import (
	"context"
	"fmt"
	"sarc-ng/internal/adapter/gorm/common"
	domainCommon "sarc-ng/internal/domain/common"
	"sarc-ng/internal/domain/reservation"

	"gorm.io/gorm"
)

// GormAdapter implements reservation.Repository using GORM
type GormAdapter struct {
	db *gorm.DB
}

// Compile-time verification that GormAdapter implements reservation.Repository
var _ reservation.Repository = (*GormAdapter)(nil)

// NewGormAdapter creates a new reservation GORM adapter
func NewGormAdapter(db *gorm.DB) *GormAdapter {
	return &GormAdapter{
		db: db,
	}
}

// ReadReservationList retrieves reservations with pagination
func (a *GormAdapter) ReadReservationList(ctx context.Context, params domainCommon.PaginationParams) ([]reservation.Reservation, int64, error) {
	var models []GormModel

	query := a.db.WithContext(ctx).Model(&GormModel{})
	query = common.BuildFilteredQuery(query, params, reservation.SearchableFields, false)

	query, total, err := common.CountAndPaginate(query, params.Offset(), params.Limit())
	if err != nil {
		return nil, 0, err
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]reservation.Reservation, len(models))
	for i, model := range models {
		entities[i] = modelToDomain(model)
	}
	return entities, total, nil
}

// ReadReservation retrieves a reservation by ID
func (a *GormAdapter) ReadReservation(ctx context.Context, id uint) (*reservation.Reservation, error) {
	var model GormModel
	if err := a.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("reservation not found: %w", domainCommon.ErrNotFound)
		}
		return nil, err
	}

	entity := modelToDomain(model)
	return &entity, nil
}

// CreateReservation adds a new reservation
func (a *GormAdapter) CreateReservation(ctx context.Context, r *reservation.Reservation) error {
	model := domainToModel(*r)
	if err := a.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	// Update the entity with generated fields
	*r = modelToDomain(model)
	return nil
}

// UpdateReservation modifies an existing reservation
func (a *GormAdapter) UpdateReservation(ctx context.Context, r *reservation.Reservation) error {
	model := domainToModel(*r)
	if err := a.db.WithContext(ctx).Save(&model).Error; err != nil {
		return err
	}

	// Update the entity with modified fields
	*r = modelToDomain(model)
	return nil
}

// DeleteReservation removes a reservation
func (a *GormAdapter) DeleteReservation(ctx context.Context, id uint) error {
	return a.db.WithContext(ctx).Delete(&GormModel{}, id).Error
}

// Migrate ensures the reservations table matches the adapter schema
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&GormModel{})
}

// domainToModel converts domain entity to GORM model
func domainToModel(entity reservation.Reservation) GormModel {
	return GormModel{
		ID:         entity.ID,
		ResourceID: entity.ResourceID,
		UserID:     entity.UserID,
		UserName:   entity.UserName,
		StartTime:  entity.StartTime,
		EndTime:    entity.EndTime,
		Purpose:    entity.Purpose,
		Status:     entity.Status,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
		DeletedAt:  common.ConvertTimeToGormDeletedAt(entity.DeletedAt),
	}
}

// modelToDomain converts GORM model to domain entity
func modelToDomain(model GormModel) reservation.Reservation {
	return reservation.Reservation{
		ID:         model.ID,
		ResourceID: model.ResourceID,
		UserID:     model.UserID,
		UserName:   model.UserName,
		StartTime:  model.StartTime,
		EndTime:    model.EndTime,
		Purpose:    model.Purpose,
		Status:     model.Status,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
		DeletedAt:  common.ConvertGormDeletedAtToTime(model.DeletedAt),
	}
}
