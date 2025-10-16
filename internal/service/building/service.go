package building

import (
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

// GetAllBuildings retrieves all buildings
func (s *Service) GetAllBuildings() ([]building.Building, error) {
	return s.repo.ReadBuildingList()
}

// GetBuilding retrieves a building by ID with validation
func (s *Service) GetBuilding(id uint) (*building.Building, error) {
	if id == 0 {
		return nil, fmt.Errorf("%w: building ID cannot be zero", common.ErrInvalidInput)
	}
	return s.repo.ReadBuilding(id)
}

// CreateBuilding creates a new building with validation
func (s *Service) CreateBuilding(b *building.Building) error {
	if strings.TrimSpace(b.Name) == "" {
		return fmt.Errorf("%w: building name cannot be empty", common.ErrInvalidInput)
	}

	if strings.TrimSpace(b.Code) == "" {
		return fmt.Errorf("%w: building code cannot be empty", common.ErrInvalidInput)
	}

	existing, err := s.repo.FindBuildingByCode(b.Code)
	if err != nil {
		return fmt.Errorf("failed to check for duplicate code: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("%w: building with code '%s' already exists", common.ErrConflict, b.Code)
	}

	return s.repo.CreateBuilding(b)
}

// UpdateBuilding updates an existing building with validation
func (s *Service) UpdateBuilding(b *building.Building) error {
	if b.ID == 0 {
		return fmt.Errorf("%w: building ID cannot be zero for update", common.ErrInvalidInput)
	}

	if strings.TrimSpace(b.Name) == "" {
		return fmt.Errorf("%w: building name cannot be empty", common.ErrInvalidInput)
	}

	if strings.TrimSpace(b.Code) == "" {
		return fmt.Errorf("%w: building code cannot be empty", common.ErrInvalidInput)
	}

	existing, err := s.repo.FindBuildingByCode(b.Code)
	if err != nil {
		return fmt.Errorf("failed to check for duplicate code: %w", err)
	}
	if existing != nil && existing.ID != b.ID {
		return fmt.Errorf("%w: building with code '%s' already exists", common.ErrConflict, b.Code)
	}

	return s.repo.UpdateBuilding(b)
}

// DeleteBuilding removes a building by ID
func (s *Service) DeleteBuilding(id uint) error {
	if id == 0 {
		return fmt.Errorf("%w: building ID cannot be zero", common.ErrInvalidInput)
	}

	_, err := s.repo.ReadBuilding(id)
	if err != nil {
		return err
	}

	return s.repo.DeleteBuilding(id)
}
