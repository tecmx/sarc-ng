package building

import (
	"sarc-ng/internal/domain/building"
)

// Mapper handles conversions between domain entities and DTOs
type Mapper struct{}

// NewMapper creates a new building mapper
func NewMapper() *Mapper {
	return &Mapper{}
}

// FromDomain converts a domain entity to DTO
func (m *Mapper) FromDomain(entity *building.Building) *BuildingDTO {
	if entity == nil {
		return nil
	}
	return &BuildingDTO{
		ID:        entity.ID,
		Name:      entity.Name,
		Code:      entity.Code,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// ToDomain converts a create DTO to domain entity
func (m *Mapper) ToDomain(dto *CreateBuildingDTO) *building.Building {
	if dto == nil {
		return nil
	}
	return &building.Building{
		Name: dto.Name,
		Code: dto.Code,
	}
}

// ToDomainWithID converts an update DTO to domain entity with ID
func (m *Mapper) ToDomainWithID(dto *UpdateBuildingDTO, id uint) *building.Building {
	if dto == nil {
		return nil
	}
	return &building.Building{
		ID:   id,
		Name: dto.Name,
		Code: dto.Code,
	}
}
