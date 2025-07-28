package lesson

import (
	"sarc-ng/internal/domain/lesson"
)

// Mapper handles conversions between domain entities and DTOs
type Mapper struct{}

// NewMapper creates a new lesson mapper
func NewMapper() *Mapper {
	return &Mapper{}
}

// FromDomain converts a domain entity to DTO
func (m *Mapper) FromDomain(entity *lesson.Lesson) *LessonDTO {
	if entity == nil {
		return nil
	}
	return &LessonDTO{
		ID:        entity.ID,
		Title:     entity.Title,
		Duration:  entity.Duration,
		StartTime: entity.StartTime,
		EndTime:   entity.EndTime,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

// ToDomain converts a create DTO to domain entity
func (m *Mapper) ToDomain(dto *CreateLessonDTO) *lesson.Lesson {
	if dto == nil {
		return nil
	}
	return &lesson.Lesson{
		Title:     dto.Title,
		Duration:  dto.Duration,
		StartTime: dto.StartTime,
	}
}

// ToDomainWithID converts an update DTO to domain entity with ID
func (m *Mapper) ToDomainWithID(dto *UpdateLessonDTO, id uint) *lesson.Lesson {
	if dto == nil {
		return nil
	}
	return &lesson.Lesson{
		ID:        id,
		Title:     dto.Title,
		Duration:  dto.Duration,
		StartTime: dto.StartTime,
	}
}
