package lesson

import (
	"context"
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

// ReadLessonList retrieves lessons with pagination
func (a *GormAdapter) ReadLessonList(ctx context.Context, params domainCommon.PaginationParams) ([]lesson.Lesson, int64, error) {
	var models []GormModel

	query := a.db.WithContext(ctx).Model(&GormModel{})
	query = common.BuildFilteredQuery(query, params, lesson.SearchableFields, false)

	query, total, err := common.CountAndPaginate(query, params.Offset(), params.Limit())
	if err != nil {
		return nil, 0, err
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]lesson.Lesson, len(models))
	for i, model := range models {
		entities[i] = modelToDomain(model)
	}
	return entities, total, nil
}

// ReadLesson retrieves a lesson by ID
func (a *GormAdapter) ReadLesson(ctx context.Context, id uint) (*lesson.Lesson, error) {
	var model GormModel
	if err := a.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("lesson not found: %w", domainCommon.ErrNotFound)
		}
		return nil, err
	}

	entity := modelToDomain(model)
	return &entity, nil
}

// CreateLesson adds a new lesson
func (a *GormAdapter) CreateLesson(ctx context.Context, l *lesson.Lesson) error {
	model := domainToModel(*l)
	if err := a.db.WithContext(ctx).Create(&model).Error; err != nil {
		return err
	}

	// Update the entity with generated fields
	*l = modelToDomain(model)
	return nil
}

// UpdateLesson modifies an existing lesson
func (a *GormAdapter) UpdateLesson(ctx context.Context, l *lesson.Lesson) error {
	model := domainToModel(*l)
	if err := a.db.WithContext(ctx).Save(&model).Error; err != nil {
		return err
	}

	// Update the entity with modified fields
	*l = modelToDomain(model)
	return nil
}

// DeleteLesson removes a lesson
func (a *GormAdapter) DeleteLesson(ctx context.Context, id uint) error {
	return a.db.WithContext(ctx).Delete(&GormModel{}, id).Error
}

// Migrate ensures the lessons table matches the adapter schema
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&GormModel{})
}

// domainToModel converts domain entity to GORM model
func domainToModel(entity lesson.Lesson) GormModel {
	return GormModel{
		ID:            entity.ID,
		ClassID:       entity.ClassID,
		ResourceID:    entity.ResourceID,
		StartTime:     entity.StartTime,
		EndTime:       entity.EndTime,
		RepeatPattern: entity.RepeatPattern,
		RepeatUntil:   entity.RepeatUntil,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		DeletedAt:     common.ConvertTimeToGormDeletedAt(entity.DeletedAt),
	}
}

// modelToDomain converts GORM model to domain entity
func modelToDomain(model GormModel) lesson.Lesson {
	return lesson.Lesson{
		ID:            model.ID,
		ClassID:       model.ClassID,
		ResourceID:    model.ResourceID,
		StartTime:     model.StartTime,
		EndTime:       model.EndTime,
		RepeatPattern: model.RepeatPattern,
		RepeatUntil:   model.RepeatUntil,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
		DeletedAt:     common.ConvertGormDeletedAtToTime(model.DeletedAt),
	}
}
