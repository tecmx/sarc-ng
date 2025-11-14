package resource

import (
	"context"
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

// GetAllResources retrieves all resources with pagination
func (s *Service) GetAllResources(ctx context.Context, params common.PaginationParams) ([]resource.Resource, common.PaginationResult, error) {
	// Validate all pagination parameters
	if err := resource.ValidateParams(params); err != nil {
		return nil, common.PaginationResult{}, err
	}

	resources, total, err := s.repo.ReadResourceList(ctx, params)
	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	result := common.NewPaginationResult(params.Page, params.PageSize, total)
	return resources, result, nil
}

// GetResource retrieves a resource by ID with validation
func (s *Service) GetResource(ctx context.Context, id uint) (*resource.Resource, error) {
	if id == 0 {
		return nil, fmt.Errorf("%w: resource ID cannot be zero", common.ErrInvalidInput)
	}
	return s.repo.ReadResource(ctx, id)
}

// CreateResource creates a new resource with validation
func (s *Service) CreateResource(ctx context.Context, r *resource.Resource) error {
	// Validate name
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("%w: resource name cannot be empty", common.ErrInvalidInput)
	}

	// Validate type
	validTypes := map[string]bool{"classroom": true, "equipment": true, "lab": true}
	if strings.TrimSpace(r.Type) == "" {
		return fmt.Errorf("%w: resource type cannot be empty", common.ErrInvalidInput)
	}
	if !validTypes[r.Type] {
		return fmt.Errorf("%w: resource type must be one of: classroom, equipment, lab", common.ErrInvalidInput)
	}

	return s.repo.CreateResource(ctx, r)
}

// UpdateResource updates an existing resource with validation
func (s *Service) UpdateResource(ctx context.Context, r *resource.Resource) error {
	if r.ID == 0 {
		return fmt.Errorf("%w: resource ID cannot be zero for update", common.ErrInvalidInput)
	}

	// Validate name
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("%w: resource name cannot be empty", common.ErrInvalidInput)
	}

	// Validate type
	validTypes := map[string]bool{"classroom": true, "equipment": true, "lab": true}
	if strings.TrimSpace(r.Type) == "" {
		return fmt.Errorf("%w: resource type cannot be empty", common.ErrInvalidInput)
	}
	if !validTypes[r.Type] {
		return fmt.Errorf("%w: resource type must be one of: classroom, equipment, lab", common.ErrInvalidInput)
	}

	return s.repo.UpdateResource(ctx, r)
}

// DeleteResource removes a resource by ID
func (s *Service) DeleteResource(ctx context.Context, id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: resource ID cannot be zero", common.ErrInvalidInput)
	}

	_, err := s.repo.ReadResource(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteResource(ctx, id)
}
