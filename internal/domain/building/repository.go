package building

// Repository defines the data access operations for buildings
// All methods are explicitly named with the Building entity
type Repository interface {
	ReadBuildingList() ([]Building, error)
	ReadBuilding(id uint) (*Building, error)
	FindBuildingByCode(code string) (*Building, error)
	CreateBuilding(building *Building) error
	UpdateBuilding(building *Building) error
	DeleteBuilding(id uint) error
}
