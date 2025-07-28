package resource

// Usecase defines the business logic operations for resource management
type Usecase interface {
	GetAllResources() ([]Resource, error)
	GetResource(id uint) (*Resource, error)
	CreateResource(resource *Resource) error
	UpdateResource(resource *Resource) error
	DeleteResource(id uint) error
	SetResourceAvailability(id uint, available bool) error
}
