package lesson

import (
	"context"
	"sarc-ng/internal/domain/common"
)

// Repository defines the data access operations for lessons
// All methods are explicitly named with the Lesson entity
type Repository interface {
	ReadLessonList(ctx context.Context, params common.PaginationParams) ([]Lesson, int64, error)
	ReadLesson(ctx context.Context, id uint) (*Lesson, error)
	CreateLesson(ctx context.Context, lesson *Lesson) error
	UpdateLesson(ctx context.Context, lesson *Lesson) error
	DeleteLesson(ctx context.Context, id uint) error
}
