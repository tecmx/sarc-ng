package lesson

import (
	"errors"
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

// GetAllLessons retrieves all lessons
func (s *Service) GetAllLessons() ([]lesson.Lesson, error) {
	return s.repo.ReadLessonList()
}

// GetLesson retrieves a lesson by ID with validation
func (s *Service) GetLesson(id uint) (*lesson.Lesson, error) {
	if id == 0 {
		return nil, errors.New("lesson ID cannot be zero")
	}
	return s.repo.ReadLesson(id)
}

// CreateLesson creates a new lesson
func (s *Service) CreateLesson(l *lesson.Lesson) error {
	return s.repo.CreateLesson(l)
}

// UpdateLesson updates an existing lesson
func (s *Service) UpdateLesson(l *lesson.Lesson) error {
	if l.ID == 0 {
		return errors.New("lesson ID cannot be zero for update")
	}

	return s.repo.UpdateLesson(l)
}

// DeleteLesson removes a lesson by ID
func (s *Service) DeleteLesson(id uint) error {
	if id == 0 {
		return errors.New("lesson ID cannot be zero")
	}

	// Check if lesson exists before deletion
	_, err := s.repo.ReadLesson(id)
	if err != nil {
		return err // Repository already returns "lesson not found" for missing records
	}

	return s.repo.DeleteLesson(id)
}
