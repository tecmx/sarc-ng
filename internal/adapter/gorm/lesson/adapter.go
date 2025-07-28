package lesson

import (
	"fmt"
	"sarc-ng/internal/adapter/gorm/common"
	domainCommon "sarc-ng/internal/domain/common"
	"sarc-ng/internal/domain/lesson"

	"gorm.io/gorm"
)

// GormAdapter implements lesson.Repository using GORM
type GormAdapter struct {
	db *gorm.DB
}

// Compile-time verification that GormAdapter implements lesson.Repository
var _ lesson.Repository = (*GormAdapter)(nil)

// NewGormAdapter creates a new lesson GORM adapter
func NewGormAdapter(db *gorm.DB) *GormAdapter {
	return &GormAdapter{
		db: db,
	}
}

// ReadLessonList retrieves all lessons
func (a *GormAdapter) ReadLessonList() ([]lesson.Lesson, error) {
	var models []GormModel
	if err := a.db.Find(&models).Error; err != nil {
		return nil, err
	}

	entities := make([]lesson.Lesson, len(models))
	for i, model := range models {
		entities[i] = modelToDomain(model)
	}
	return entities, nil
}

// ReadLesson retrieves a lesson by ID
func (a *GormAdapter) ReadLesson(id uint) (*lesson.Lesson, error) {
	var model GormModel
	if err := a.db.First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("lesson not found: %w", domainCommon.ErrNotFound)
		}
		return nil, err
	}

	entity := modelToDomain(model)
	return &entity, nil
}

// CreateLesson adds a new lesson
func (a *GormAdapter) CreateLesson(l *lesson.Lesson) error {
	model := domainToModel(*l)
	if err := a.db.Create(&model).Error; err != nil {
		return err
	}

	// Update the entity with generated fields
	*l = modelToDomain(model)
	return nil
}

// UpdateLesson modifies an existing lesson
func (a *GormAdapter) UpdateLesson(l *lesson.Lesson) error {
	model := domainToModel(*l)
	if err := a.db.Save(&model).Error; err != nil {
		return err
	}

	// Update the entity with modified fields
	*l = modelToDomain(model)
	return nil
}

// DeleteLesson removes a lesson
func (a *GormAdapter) DeleteLesson(id uint) error {
	return a.db.Delete(&GormModel{}, id).Error
}

// domainToModel converts domain entity to GORM model
func domainToModel(entity lesson.Lesson) GormModel {
	return GormModel{
		ID:          entity.ID,
		Title:       entity.Title,
		Duration:    entity.Duration,
		Description: entity.Description,
		StartTime:   entity.StartTime,
		EndTime:     entity.EndTime,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		DeletedAt:   common.ConvertTimeToGormDeletedAt(entity.DeletedAt),
	}
}

// modelToDomain converts GORM model to domain entity
func modelToDomain(model GormModel) lesson.Lesson {
	return lesson.Lesson{
		ID:          model.ID,
		Title:       model.Title,
		Duration:    model.Duration,
		Description: model.Description,
		StartTime:   model.StartTime,
		EndTime:     model.EndTime,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		DeletedAt:   common.ConvertGormDeletedAtToTime(model.DeletedAt),
	}
}
