package client

// LessonsService provides methods for lesson operations
type LessonsService struct {
	*baseService
}

// Lessons returns the lessons service
func (c *Client) Lessons() *LessonsService {
	return &LessonsService{
		baseService: newBaseService(c, "lessons"),
	}
}

// List retrieves all lessons with pagination
func (s *LessonsService) List(page, pageSize int) ([]byte, error) {
	return s.list(page, pageSize)
}

// Get retrieves a specific lesson by ID
func (s *LessonsService) Get(id uint) ([]byte, error) {
	return s.get(id)
}

// Create creates a new lesson
func (s *LessonsService) Create(req interface{}) ([]byte, error) {
	return s.create(req)
}

// Update updates an existing lesson
func (s *LessonsService) Update(id uint, req interface{}) ([]byte, error) {
	return s.update(id, req)
}

// Delete removes a lesson by ID
func (s *LessonsService) Delete(id uint) error {
	return s.delete(id)
}
