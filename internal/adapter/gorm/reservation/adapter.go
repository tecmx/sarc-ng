package reservation

import (
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

// ReadReservationList retrieves all reservations
func (a *GormAdapter) ReadReservationList() ([]reservation.Reservation, error) {
	var models []GormModel
	if err := a.db.Find(&models).Error; err != nil {
		return nil, err
	}

	entities := make([]reservation.Reservation, len(models))
	for i, model := range models {
		entities[i] = modelToDomain(model)
	}
	return entities, nil
}

// ReadReservation retrieves a reservation by ID
func (a *GormAdapter) ReadReservation(id uint) (*reservation.Reservation, error) {
	var model GormModel
	if err := a.db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("reservation not found: %w", domainCommon.ErrNotFound)
		}
		return nil, err
	}

	entity := modelToDomain(model)
	return &entity, nil
}

// CreateReservation adds a new reservation
func (a *GormAdapter) CreateReservation(r *reservation.Reservation) error {
	model := domainToModel(*r)
	if err := a.db.Create(&model).Error; err != nil {
		return err
	}

	// Update the entity with generated fields
	*r = modelToDomain(model)
	return nil
}

// UpdateReservation modifies an existing reservation
func (a *GormAdapter) UpdateReservation(r *reservation.Reservation) error {
	model := domainToModel(*r)
	if err := a.db.Save(&model).Error; err != nil {
		return err
	}

	// Update the entity with modified fields
	*r = modelToDomain(model)
	return nil
}

// DeleteReservation removes a reservation
func (a *GormAdapter) DeleteReservation(id uint) error {
	return a.db.Delete(&GormModel{}, id).Error
}

// domainToModel converts domain entity to GORM model
func domainToModel(entity reservation.Reservation) GormModel {
	return GormModel{
		ID:          entity.ID,
		ResourceID:  entity.ResourceID,
		UserID:      entity.UserID,
		StartTime:   entity.StartTime,
		EndTime:     entity.EndTime,
		Purpose:     entity.Purpose,
		Status:      entity.Status,
		Description: entity.Description,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		DeletedAt:   common.ConvertTimeToGormDeletedAt(entity.DeletedAt),
	}
}

// modelToDomain converts GORM model to domain entity
func modelToDomain(model GormModel) reservation.Reservation {
	return reservation.Reservation{
		ID:          model.ID,
		ResourceID:  model.ResourceID,
		UserID:      model.UserID,
		StartTime:   model.StartTime,
		EndTime:     model.EndTime,
		Purpose:     model.Purpose,
		Status:      model.Status,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		DeletedAt:   common.ConvertGormDeletedAtToTime(model.DeletedAt),
	}
}
