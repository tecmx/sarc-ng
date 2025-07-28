package client

import "fmt"

// ReservationsService provides methods for reservation operations
type ReservationsService struct {
	client *Client
}

// Reservations returns the reservations service
func (c *Client) Reservations() *ReservationsService {
	return &ReservationsService{client: c}
}

// List retrieves all reservations with pagination
func (s *ReservationsService) List(page, pageSize int) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/reservations?page=%d&pageSize=%d", page, pageSize)
	resp, err := s.client.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Get retrieves a specific reservation by ID
func (s *ReservationsService) Get(id uint) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/reservations/%d", id)
	resp, err := s.client.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Create creates a new reservation
func (s *ReservationsService) Create(req interface{}) ([]byte, error) {
	resp, err := s.client.doRequest("POST", "/api/v1/reservations", req)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Update updates an existing reservation
func (s *ReservationsService) Update(id uint, req interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/reservations/%d", id)
	resp, err := s.client.doRequest("PUT", endpoint, req)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Delete removes a reservation by ID
func (s *ReservationsService) Delete(id uint) error {
	endpoint := fmt.Sprintf("/api/v1/reservations/%d", id)
	resp, err := s.client.doRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	_, err = s.client.handleRawResponse(resp)
	return err
}
