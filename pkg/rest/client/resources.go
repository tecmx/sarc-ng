package client

// ResourcesService provides methods for resource operations
type ResourcesService struct {
	*baseService
}

// Resources returns the resources service
func (c *Client) Resources() *ResourcesService {
	return &ResourcesService{
		baseService: newBaseService(c, "resources"),
	}
}

// List retrieves all resources with pagination
func (s *ResourcesService) List(page, pageSize int) ([]byte, error) {
	return s.list(page, pageSize)
}

// Get retrieves a specific resource by ID
func (s *ResourcesService) Get(id uint) ([]byte, error) {
	return s.get(id)
}

// Create creates a new resource
func (s *ResourcesService) Create(req interface{}) ([]byte, error) {
	return s.create(req)
}

// Update updates an existing resource
func (s *ResourcesService) Update(id uint, req interface{}) ([]byte, error) {
	return s.update(id, req)
}

// Delete removes a resource by ID
func (s *ResourcesService) Delete(id uint) error {
	return s.delete(id)
}
