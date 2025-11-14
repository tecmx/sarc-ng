package lesson

import (
	"context"
	"sarc-ng/internal/domain/common"
)

// Usecase defines the business logic operations for lesson management
type Usecase interface {
	GetAllLessons(ctx context.Context, params common.PaginationParams) ([]Lesson, common.PaginationResult, error)
	GetLesson(ctx context.Context, id uint) (*Lesson, error)
	CreateLesson(ctx context.Context, lesson *Lesson) error
	UpdateLesson(ctx context.Context, lesson *Lesson) error
	DeleteLesson(ctx context.Context, id uint) error
}
