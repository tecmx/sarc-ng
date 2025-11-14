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
		ID:            entity.ID,
		ClassID:       entity.ClassID,
		ResourceID:    entity.ResourceID,
		StartTime:     entity.StartTime,
		EndTime:       entity.EndTime,
		RepeatPattern: entity.RepeatPattern,
		RepeatUntil:   entity.RepeatUntil,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
}

// ToDomain converts a create DTO to domain entity
func (m *Mapper) ToDomain(dto *CreateLessonDTO) *lesson.Lesson {
	if dto == nil {
		return nil
	}
	return &lesson.Lesson{
		ClassID:       dto.ClassID,
		ResourceID:    dto.ResourceID,
		StartTime:     dto.StartTime,
		EndTime:       dto.EndTime,
		RepeatPattern: dto.RepeatPattern,
		RepeatUntil:   dto.RepeatUntil,
	}
}

// ToDomainWithID converts an update DTO to domain entity with ID
func (m *Mapper) ToDomainWithID(dto *UpdateLessonDTO, id uint) *lesson.Lesson {
	if dto == nil {
		return nil
	}
	return &lesson.Lesson{
		ID:            id,
		ClassID:       dto.ClassID,
		ResourceID:    dto.ResourceID,
		StartTime:     dto.StartTime,
		EndTime:       dto.EndTime,
		RepeatPattern: dto.RepeatPattern,
		RepeatUntil:   dto.RepeatUntil,
	}
}
