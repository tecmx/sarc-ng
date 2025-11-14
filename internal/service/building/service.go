package building

import (
	"context"
	"fmt"
	"sarc-ng/internal/domain/building"
	"sarc-ng/internal/domain/common"
	"strings"
)

// Service implements building.Usecase interface
type Service struct {
	repo building.Repository
}

// Compile-time verification that Service implements building.Usecase
var _ building.Usecase = (*Service)(nil)

// NewService creates a new building service
func NewService(repo building.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetAllBuildings retrieves all buildings with pagination
func (s *Service) GetAllBuildings(ctx context.Context, params common.PaginationParams) ([]building.Building, common.PaginationResult, error) {
	// Validate all pagination parameters
	if err := building.ValidateParams(params); err != nil {
		return nil, common.PaginationResult{}, err
	}

	buildings, total, err := s.repo.ReadBuildingList(ctx, params)
	if err != nil {
		return nil, common.PaginationResult{}, err
	}

	result := common.NewPaginationResult(params.Page, params.PageSize, total)
	return buildings, result, nil
}

// GetBuilding retrieves a building by ID with validation
func (s *Service) GetBuilding(ctx context.Context, id uint) (*building.Building, error) {
	if id == 0 {
		return nil, fmt.Errorf("%w: building ID cannot be zero", common.ErrInvalidInput)
	}
	return s.repo.ReadBuilding(ctx, id)
}

// CreateBuilding creates a new building with validation
func (s *Service) CreateBuilding(ctx context.Context, b *building.Building) error {
	if strings.TrimSpace(b.Name) == "" {
		return fmt.Errorf("%w: building name cannot be empty", common.ErrInvalidInput)
	}

	if strings.TrimSpace(b.Code) == "" {
		return fmt.Errorf("%w: building code cannot be empty", common.ErrInvalidInput)
	}

	existing, err := s.repo.FindBuildingByCode(ctx, b.Code)
	if err != nil {
		return fmt.Errorf("failed to check for duplicate code: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("%w: building with code '%s' already exists", common.ErrConflict, b.Code)
	}

	return s.repo.CreateBuilding(ctx, b)
}

// UpdateBuilding updates an existing building with validation
func (s *Service) UpdateBuilding(ctx context.Context, b *building.Building) error {
	if b.ID == 0 {
		return fmt.Errorf("%w: building ID cannot be zero for update", common.ErrInvalidInput)
	}

	if strings.TrimSpace(b.Name) == "" {
		return fmt.Errorf("%w: building name cannot be empty", common.ErrInvalidInput)
	}

	if strings.TrimSpace(b.Code) == "" {
		return fmt.Errorf("%w: building code cannot be empty", common.ErrInvalidInput)
	}

	existing, err := s.repo.FindBuildingByCode(ctx, b.Code)
	if err != nil {
		return fmt.Errorf("failed to check for duplicate code: %w", err)
	}
	if existing != nil && existing.ID != b.ID {
		return fmt.Errorf("%w: building with code '%s' already exists", common.ErrConflict, b.Code)
	}

	return s.repo.UpdateBuilding(ctx, b)
}

// DeleteBuilding removes a building by ID
func (s *Service) DeleteBuilding(ctx context.Context, id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: building ID cannot be zero", common.ErrInvalidInput)
	}

	_, err := s.repo.ReadBuilding(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.DeleteBuilding(ctx, id)
}
