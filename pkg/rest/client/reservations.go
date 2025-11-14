package client

// ReservationsService provides methods for reservation operations
type ReservationsService struct {
	*baseService
}

// Reservations returns the reservations service
func (c *Client) Reservations() *ReservationsService {
	return &ReservationsService{
		baseService: newBaseService(c, "reservations"),
	}
}

// List retrieves all reservations with pagination
func (s *ReservationsService) List(page, pageSize int) ([]byte, error) {
	return s.list(page, pageSize)
}

// Get retrieves a specific reservation by ID
func (s *ReservationsService) Get(id uint) ([]byte, error) {
	return s.get(id)
}

// Create creates a new reservation
func (s *ReservationsService) Create(req interface{}) ([]byte, error) {
	return s.create(req)
}

// Update updates an existing reservation
func (s *ReservationsService) Update(id uint, req interface{}) ([]byte, error) {
	return s.update(id, req)
}

// Delete removes a reservation by ID
func (s *ReservationsService) Delete(id uint) error {
	return s.delete(id)
}
