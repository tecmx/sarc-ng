package resource

import (
	"errors"
	"sarc-ng/internal/domain/resource"
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
		return nil, errors.New("resource ID cannot be zero")
	}
	return s.repo.ReadResource(id)
}

// CreateResource creates a new resource
func (s *Service) CreateResource(r *resource.Resource) error {
	return s.repo.CreateResource(r)
}

// UpdateResource updates an existing resource
func (s *Service) UpdateResource(r *resource.Resource) error {
	if r.ID == 0 {
		return errors.New("resource ID cannot be zero for update")
	}

	return s.repo.UpdateResource(r)
}

// DeleteResource removes a resource by ID
func (s *Service) DeleteResource(id uint) error {
	if id == 0 {
		return errors.New("resource ID cannot be zero")
	}

	// Check if resource exists before deletion
	_, err := s.repo.ReadResource(id)
	if err != nil {
		return err // Repository already returns "resource not found" for missing records
	}

	return s.repo.DeleteResource(id)
}

// SetResourceAvailability sets the availability status of a resource
func (s *Service) SetResourceAvailability(id uint, available bool) error {
	if id == 0 {
		return errors.New("resource ID cannot be zero")
	}

	resource, err := s.repo.ReadResource(id)
	if err != nil {
		return err
	}
	if resource == nil {
		return errors.New("resource not found")
	}

	resource.IsAvailable = available
	return s.repo.UpdateResource(resource)
}
