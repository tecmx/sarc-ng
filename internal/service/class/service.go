package class

import (
	"errors"
	"sarc-ng/internal/domain/class"
)

// Service implements class.Usecase interface
type Service struct {
	repo class.Repository
}

// Compile-time verification that Service implements class.Usecase
var _ class.Usecase = (*Service)(nil)

// NewService creates a new class service
func NewService(repo class.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetAllClasses retrieves all classes
func (s *Service) GetAllClasses() ([]class.Class, error) {
	return s.repo.ReadClassList()
}

// GetClass retrieves a class by ID with validation
func (s *Service) GetClass(id uint) (*class.Class, error) {
	if id == 0 {
		return nil, errors.New("class ID cannot be zero")
	}
	return s.repo.ReadClass(id)
}

// CreateClass creates a new class
func (s *Service) CreateClass(c *class.Class) error {
	return s.repo.CreateClass(c)
}

// UpdateClass updates an existing class
func (s *Service) UpdateClass(c *class.Class) error {
	if c.ID == 0 {
		return errors.New("class ID cannot be zero for update")
	}

	return s.repo.UpdateClass(c)
}

// DeleteClass removes a class by ID
func (s *Service) DeleteClass(id uint) error {
	if id == 0 {
		return errors.New("class ID cannot be zero")
	}

	// Check if class exists before deletion
	_, err := s.repo.ReadClass(id)
	if err != nil {
		return err // Repository already returns "class not found" for missing records
	}

	return s.repo.DeleteClass(id)
}
