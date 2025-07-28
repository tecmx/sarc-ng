package class

import (
	"sarc-ng/internal/domain/class"
)

// Mapper handles conversions between domain entities and DTOs
type Mapper struct{}

// NewMapper creates a new class mapper
func NewMapper() *Mapper {
	return &Mapper{}
}

// FromDomain converts a domain entity to DTO
func (m *Mapper) FromDomain(entity *class.Class) *ClassDTO {
	if entity == nil {
		return nil
	}
	return &ClassDTO{
		ID:        entity.ID,
		Name:      entity.Name,
		Capacity:  entity.Capacity,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// ToDomain converts a create DTO to domain entity
func (m *Mapper) ToDomain(dto *CreateClassDTO) *class.Class {
	if dto == nil {
		return nil
	}
	return &class.Class{
		Name:     dto.Name,
		Capacity: dto.Capacity,
	}
}

// ToDomainWithID converts an update DTO to domain entity with ID
func (m *Mapper) ToDomainWithID(dto *UpdateClassDTO, id uint) *class.Class {
	if dto == nil {
		return nil
	}
	return &class.Class{
		ID:       id,
		Name:     dto.Name,
		Capacity: dto.Capacity,
	}
}
