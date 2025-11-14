package client

// ClassesService provides methods for class operations
type ClassesService struct {
	*baseService
}

// Classes returns the classes service
func (c *Client) Classes() *ClassesService {
	return &ClassesService{
		baseService: newBaseService(c, "classes"),
	}
}

// List retrieves all classes with pagination
func (s *ClassesService) List(page, pageSize int) ([]byte, error) {
	return s.list(page, pageSize)
}

// Get retrieves a specific class by ID
func (s *ClassesService) Get(id uint) ([]byte, error) {
	return s.get(id)
}

// Create creates a new class
func (s *ClassesService) Create(req interface{}) ([]byte, error) {
	return s.create(req)
}

// Update updates an existing class
func (s *ClassesService) Update(id uint, req interface{}) ([]byte, error) {
	return s.update(id, req)
}

// Delete removes a class by ID
func (s *ClassesService) Delete(id uint) error {
	return s.delete(id)
}
