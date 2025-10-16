package resource

import (
	"fmt"
	"sarc-ng/internal/domain/common"
	"sarc-ng/internal/domain/resource"
	"strings"
)

// Service implements resource.Usecase interface
type Service struct {
	repo resource.Repository
}

// Compile-time verification that Service implements resource.Usecase
var _ resource.Usecase = (*Service)(nil)

// NewService creates a new resource service
func NewService(repo resource.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetAllResources retrieves all resources
func (s *Service) GetAllResources() ([]resource.Resource, error) {
	return s.repo.ReadResourceList()
}

// GetResource retrieves a resource by ID with validation
func (s *Service) GetResource(id uint) (*resource.Resource, error) {
	if id == 0 {
		return nil, fmt.Errorf("%w: resource ID cannot be zero", common.ErrInvalidInput)
	}
	return s.repo.ReadResource(id)
}

// CreateResource creates a new resource with validation
func (s *Service) CreateResource(r *resource.Resource) error {
	// Validate name
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("%w: resource name cannot be empty", common.ErrInvalidInput)
	}

	// Validate type
	if strings.TrimSpace(r.Type) == "" {
		return fmt.Errorf("%w: resource type cannot be empty", common.ErrInvalidInput)
	}

	return s.repo.CreateResource(r)
}

// UpdateResource updates an existing resource with validation
func (s *Service) UpdateResource(r *resource.Resource) error {
	if r.ID == 0 {
		return fmt.Errorf("%w: resource ID cannot be zero for update", common.ErrInvalidInput)
	}

	// Validate name
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("%w: resource name cannot be empty", common.ErrInvalidInput)
	}

	// Validate type
	if strings.TrimSpace(r.Type) == "" {
		return fmt.Errorf("%w: resource type cannot be empty", common.ErrInvalidInput)
	}

	return s.repo.UpdateResource(r)
}

// DeleteResource removes a resource by ID
func (s *Service) DeleteResource(id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: resource ID cannot be zero", common.ErrInvalidInput)
	}

	_, err := s.repo.ReadResource(id)
	if err != nil {
		return err
	}

	return s.repo.DeleteResource(id)
}

// SetResourceAvailability sets the availability status of a resource
func (s *Service) SetResourceAvailability(id uint, available bool) error {
	if id == 0 {
		return fmt.Errorf("%w: resource ID cannot be zero", common.ErrInvalidInput)
	}

	resource, err := s.repo.ReadResource(id)
	if err != nil {
		return err
	}
	if resource == nil {
		return fmt.Errorf("%w: resource not found", common.ErrNotFound)
	}

	resource.IsAvailable = available
	return s.repo.UpdateResource(resource)
}
