package resource

import (
	"sarc-ng/internal/domain/resource"
)

// Mapper handles conversions between domain entities and DTOs
type Mapper struct{}

// NewMapper creates a new resource mapper
func NewMapper() *Mapper {
	return &Mapper{}
}

// FromDomain converts a domain entity to DTO
func (m *Mapper) FromDomain(entity *resource.Resource) *ResourceDTO {
	if entity == nil {
		return nil
	}
	return &ResourceDTO{
		ID:          entity.ID,
		Name:        entity.Name,
		Type:        entity.Type,
		Description: entity.Description,
		Location:    entity.Location,
		IsAvailable: entity.IsAvailable,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

// ToDomain converts a create DTO to domain entity
func (m *Mapper) ToDomain(dto *CreateResourceDTO) *resource.Resource {
	if dto == nil {
		return nil
	}
	return &resource.Resource{
		Name:        dto.Name,
		Type:        dto.Type,
		Description: dto.Description,
		Location:    dto.Location,
		IsAvailable: dto.IsAvailable,
	}
}

// ToDomainWithID converts an update DTO to domain entity with ID
func (m *Mapper) ToDomainWithID(dto *UpdateResourceDTO, id uint) *resource.Resource {
	if dto == nil {
		return nil
	}
	return &resource.Resource{
		ID:          id,
		Name:        dto.Name,
		Type:        dto.Type,
		Description: dto.Description,
		Location:    dto.Location,
		IsAvailable: dto.IsAvailable,
	}
}
