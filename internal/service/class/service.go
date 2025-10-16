package class

import (
	"fmt"
	"sarc-ng/internal/domain/class"
	"sarc-ng/internal/domain/common"
	"strings"
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
		return nil, fmt.Errorf("%w: class ID cannot be zero", common.ErrInvalidInput)
	}
	return s.repo.ReadClass(id)
}

// CreateClass creates a new class with validation
func (s *Service) CreateClass(c *class.Class) error {
	// Validate name
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("%w: class name cannot be empty", common.ErrInvalidInput)
	}

	// Validate capacity
	if c.Capacity <= 0 {
		return fmt.Errorf("%w: class capacity must be greater than zero", common.ErrInvalidInput)
	}

	return s.repo.CreateClass(c)
}

// UpdateClass updates an existing class with validation
func (s *Service) UpdateClass(c *class.Class) error {
	if c.ID == 0 {
		return fmt.Errorf("%w: class ID cannot be zero for update", common.ErrInvalidInput)
	}

	// Validate name
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("%w: class name cannot be empty", common.ErrInvalidInput)
	}

	// Validate capacity
	if c.Capacity <= 0 {
		return fmt.Errorf("%w: class capacity must be greater than zero", common.ErrInvalidInput)
	}

	return s.repo.UpdateClass(c)
}

// DeleteClass removes a class by ID
func (s *Service) DeleteClass(id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: class ID cannot be zero", common.ErrInvalidInput)
	}

	_, err := s.repo.ReadClass(id)
	if err != nil {
		return err
	}

	return s.repo.DeleteClass(id)
}
