package resource

// Repository defines the data access operations for resources
// All methods are explicitly named with the Resource entity
type Repository interface {
	ReadResourceList() ([]Resource, error)
	ReadResource(id uint) (*Resource, error)
	CreateResource(resource *Resource) error
	UpdateResource(resource *Resource) error
	DeleteResource(id uint) error
}
