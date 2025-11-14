package lesson

import (
	"context"
	"fmt"
	"sarc-ng/internal/domain/common"
	"sarc-ng/internal/domain/lesson"
)

// Service implements lesson.Usecase interface
type Service struct {
	repo lesson.Repository
}

// Compile-time verification that Service implements lesson.Usecase
var _ lesson.Usecase = (*Service)(nil)

// NewService creates a new lesson service
func NewService(repo lesson.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetAllLessons retrieves all lessons with pagination
func (s *Service) GetAllLessons(ctx context.Context, params common.PaginationParams) ([]lesson.Lesson, common.PaginationResult, error) {
	// Validate all pagination parameters
	if err := lesson.ValidateParams(params); err != nil {
		return nil, common.PaginationResult{}, err
	}

	lessons, total, err := s.repo.ReadLessonList(ctx, params)
	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	result := common.NewPaginationResult(params.Page, params.PageSize, total)
	return lessons, result, nil
}

// GetLesson retrieves a lesson by ID with validation
func (s *Service) GetLesson(ctx context.Context, id uint) (*lesson.Lesson, error) {
	if id == 0 {
		return nil, fmt.Errorf("%w: lesson ID cannot be zero", common.ErrInvalidInput)
	}
	return s.repo.ReadLesson(ctx, id)
}

// CreateLesson creates a new lesson with validation
func (s *Service) CreateLesson(ctx context.Context, l *lesson.Lesson) error {
	// Validate class ID
	if l.ClassID == 0 {
		return fmt.Errorf("%w: class ID cannot be zero", common.ErrInvalidInput)
	}

	// Validate resource ID
	if l.ResourceID == 0 {
		return fmt.Errorf("%w: resource ID cannot be zero", common.ErrInvalidInput)
	}

	// Validate time range
	if l.EndTime.Before(l.StartTime) || l.EndTime.Equal(l.StartTime) {
		return fmt.Errorf("%w: end time must be after start time", common.ErrInvalidInput)
	}

	// Validate repeat pattern
	validPatterns := map[string]bool{"none": true, "daily": true, "weekly": true, "biweekly": true, "monthly": true}
	if !validPatterns[l.RepeatPattern] {
		return fmt.Errorf("%w: repeat pattern must be one of: none, daily, weekly, biweekly, monthly", common.ErrInvalidInput)
	}

	return s.repo.CreateLesson(ctx, l)
}

// UpdateLesson updates an existing lesson with validation
func (s *Service) UpdateLesson(ctx context.Context, l *lesson.Lesson) error {
	if l.ID == 0 {
		return fmt.Errorf("%w: lesson ID cannot be zero for update", common.ErrInvalidInput)
	}

	// Validate class ID
	if l.ClassID == 0 {
		return fmt.Errorf("%w: class ID cannot be zero", common.ErrInvalidInput)
	}

	// Validate resource ID
	if l.ResourceID == 0 {
		return fmt.Errorf("%w: resource ID cannot be zero", common.ErrInvalidInput)
	}

	// Validate time range
	if l.EndTime.Before(l.StartTime) || l.EndTime.Equal(l.StartTime) {
		return fmt.Errorf("%w: end time must be after start time", common.ErrInvalidInput)
	}

	// Validate repeat pattern
	validPatterns := map[string]bool{"none": true, "daily": true, "weekly": true, "biweekly": true, "monthly": true}
	if !validPatterns[l.RepeatPattern] {
		return fmt.Errorf("%w: repeat pattern must be one of: none, daily, weekly, biweekly, monthly", common.ErrInvalidInput)
	}

	return s.repo.UpdateLesson(ctx, l)
}

// DeleteLesson removes a lesson by ID
func (s *Service) DeleteLesson(ctx context.Context, id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: lesson ID cannot be zero", common.ErrInvalidInput)
	}

	_, err := s.repo.ReadLesson(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteLesson(ctx, id)
}
