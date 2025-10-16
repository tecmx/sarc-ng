package reservation

import (
	"sarc-ng/internal/domain/reservation"
)

// Mapper handles conversions between domain entities and DTOs
type Mapper struct{}

// NewMapper creates a new reservation mapper
func NewMapper() *Mapper {
	return &Mapper{}
}

// FromDomain converts a domain entity to DTO
func (m *Mapper) FromDomain(entity *reservation.Reservation) *ReservationDTO {
	if entity == nil {
		return nil
	}
	return &ReservationDTO{
		ID:          entity.ID,
		ResourceID:  entity.ResourceID,
		UserID:      entity.UserID,
		StartTime:   entity.StartTime,
		EndTime:     entity.EndTime,
		Purpose:     entity.Purpose,
		Description: entity.Description,
		Status:      entity.Status,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

// ToDomain converts a create DTO to domain entity
func (m *Mapper) ToDomain(dto *CreateReservationDTO) *reservation.Reservation {
	if dto == nil {
		return nil
	}
	return &reservation.Reservation{
		ResourceID:  dto.ResourceID,
		UserID:      dto.UserID,
		StartTime:   dto.StartTime,
		EndTime:     dto.EndTime,
		Purpose:     dto.Purpose,
		Description: dto.Description,
		Status:      dto.Status,
	}
}

// ToDomainWithID converts an update DTO to domain entity with ID
func (m *Mapper) ToDomainWithID(dto *UpdateReservationDTO, id uint) *reservation.Reservation {
	if dto == nil {
		return nil
	}
	return &reservation.Reservation{
		ID:          id,
		ResourceID:  dto.ResourceID,
		UserID:      dto.UserID,
		StartTime:   dto.StartTime,
		EndTime:     dto.EndTime,
		Purpose:     dto.Purpose,
		Description: dto.Description,
		Status:      dto.Status,
	}
}
