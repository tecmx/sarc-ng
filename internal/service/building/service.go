package building

import (
	"errors"
	"sarc-ng/internal/domain/building"
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
		return nil, errors.New("building ID cannot be zero")
	}
	return s.repo.ReadBuilding(id)
}

// CreateBuilding creates a new building
func (s *Service) CreateBuilding(b *building.Building) error {
	return s.repo.CreateBuilding(b)
}

// UpdateBuilding updates an existing building
func (s *Service) UpdateBuilding(b *building.Building) error {
	if b.ID == 0 {
		return errors.New("building ID cannot be zero for update")
	}

	return s.repo.UpdateBuilding(b)
}

// DeleteBuilding removes a building by ID
func (s *Service) DeleteBuilding(id uint) error {
	if id == 0 {
		return errors.New("building ID cannot be zero")
	}

	// Check if building exists before deletion
	_, err := s.repo.ReadBuilding(id)
	if err != nil {
		return err // Repository already returns "building not found" for missing records
	}

	return s.repo.DeleteBuilding(id)
}
