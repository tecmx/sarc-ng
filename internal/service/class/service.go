package class

import (
	"context"
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

// GetAllClasses retrieves all classes with pagination
func (s *Service) GetAllClasses(ctx context.Context, params common.PaginationParams) ([]class.Class, common.PaginationResult, error) {
	// Validate all pagination parameters
	if err := class.ValidateParams(params); err != nil {
		return nil, common.PaginationResult{}, err
	}

	classes, total, err := s.repo.ReadClassList(ctx, params)
	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	result := common.NewPaginationResult(params.Page, params.PageSize, total)
	return classes, result, nil
}

// GetClass retrieves a class by ID with validation
func (s *Service) GetClass(ctx context.Context, id uint) (*class.Class, error) {
	if id == 0 {
		return nil, fmt.Errorf("%w: class ID cannot be zero", common.ErrInvalidInput)
	}
	return s.repo.ReadClass(ctx, id)
}

// CreateClass creates a new class with validation
func (s *Service) CreateClass(ctx context.Context, c *class.Class) error {
	// Validate name
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("%w: class name cannot be empty", common.ErrInvalidInput)
	}

	// Validate code
	if strings.TrimSpace(c.Code) == "" {
		return fmt.Errorf("%w: class code cannot be empty", common.ErrInvalidInput)
	}

	// Validate instructor name
	if strings.TrimSpace(c.InstructorName) == "" {
		return fmt.Errorf("%w: instructor name cannot be empty", common.ErrInvalidInput)
	}

	// Validate capacity
	if c.Capacity <= 0 {
		return fmt.Errorf("%w: class capacity must be greater than zero", common.ErrInvalidInput)
	}

	return s.repo.CreateClass(ctx, c)
}

// UpdateClass updates an existing class with validation
func (s *Service) UpdateClass(ctx context.Context, c *class.Class) error {
	if c.ID == 0 {
		return fmt.Errorf("%w: class ID cannot be zero for update", common.ErrInvalidInput)
	}

	// Validate name
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("%w: class name cannot be empty", common.ErrInvalidInput)
	}

	// Validate code
	if strings.TrimSpace(c.Code) == "" {
		return fmt.Errorf("%w: class code cannot be empty", common.ErrInvalidInput)
	}

	// Validate instructor name
	if strings.TrimSpace(c.InstructorName) == "" {
		return fmt.Errorf("%w: instructor name cannot be empty", common.ErrInvalidInput)
	}

	// Validate capacity
	if c.Capacity <= 0 {
		return fmt.Errorf("%w: class capacity must be greater than zero", common.ErrInvalidInput)
	}

	return s.repo.UpdateClass(ctx, c)
}

// DeleteClass removes a class by ID
func (s *Service) DeleteClass(ctx context.Context, id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: class ID cannot be zero", common.ErrInvalidInput)
	}

	_, err := s.repo.ReadClass(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteClass(ctx, id)
}
