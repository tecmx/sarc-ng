package lesson

import (
	"fmt"
	"sarc-ng/internal/domain/common"
	"sarc-ng/internal/domain/lesson"
	"strings"
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

// GetAllLessons retrieves all lessons
func (s *Service) GetAllLessons() ([]lesson.Lesson, error) {
	return s.repo.ReadLessonList()
}

// GetLesson retrieves a lesson by ID with validation
func (s *Service) GetLesson(id uint) (*lesson.Lesson, error) {
	if id == 0 {
		return nil, fmt.Errorf("%w: lesson ID cannot be zero", common.ErrInvalidInput)
	}
	return s.repo.ReadLesson(id)
}

// CreateLesson creates a new lesson with validation
func (s *Service) CreateLesson(l *lesson.Lesson) error {
	// Validate title
	if strings.TrimSpace(l.Title) == "" {
		return fmt.Errorf("%w: lesson title cannot be empty", common.ErrInvalidInput)
	}

	// Validate duration
	if l.Duration <= 0 {
		return fmt.Errorf("%w: lesson duration must be greater than zero", common.ErrInvalidInput)
	}

	return s.repo.CreateLesson(l)
}

// UpdateLesson updates an existing lesson with validation
func (s *Service) UpdateLesson(l *lesson.Lesson) error {
	if l.ID == 0 {
		return fmt.Errorf("%w: lesson ID cannot be zero for update", common.ErrInvalidInput)
	}

	// Validate title
	if strings.TrimSpace(l.Title) == "" {
		return fmt.Errorf("%w: lesson title cannot be empty", common.ErrInvalidInput)
	}

	// Validate duration
	if l.Duration <= 0 {
		return fmt.Errorf("%w: lesson duration must be greater than zero", common.ErrInvalidInput)
	}

	return s.repo.UpdateLesson(l)
}

// DeleteLesson removes a lesson by ID
func (s *Service) DeleteLesson(id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: lesson ID cannot be zero", common.ErrInvalidInput)
	}

	_, err := s.repo.ReadLesson(id)
	if err != nil {
		return err
	}

	return s.repo.DeleteLesson(id)
}
