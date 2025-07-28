package building

// Usecase defines the business logic operations for building management
type Usecase interface {
	GetAllBuildings() ([]Building, error)
	GetBuilding(id uint) (*Building, error)
	CreateBuilding(building *Building) error
	UpdateBuilding(building *Building) error
	DeleteBuilding(id uint) error
}
