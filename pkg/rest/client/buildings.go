package client

import (
	"fmt"
)

// BuildingsService provides methods for building operations
type BuildingsService struct {
	client *Client
}

// Buildings returns the buildings service
func (c *Client) Buildings() *BuildingsService {
	return &BuildingsService{client: c}
}

// List retrieves all buildings with pagination
func (s *BuildingsService) List(page, pageSize int) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/buildings?page=%d&pageSize=%d", page, pageSize)

	resp, err := s.client.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Get retrieves a specific building by ID
func (s *BuildingsService) Get(id uint) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/buildings/%d", id)

	resp, err := s.client.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Search searches buildings by name and/or code with pagination
func (s *BuildingsService) Search(name, code string, page, pageSize int) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/buildings/search?name=%s&code=%s&page=%d&pageSize=%d", name, code, page, pageSize)

	resp, err := s.client.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Create creates a new building
func (s *BuildingsService) Create(req interface{}) ([]byte, error) {
	resp, err := s.client.doRequest("POST", "/api/v1/buildings", req)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Update updates an existing building
func (s *BuildingsService) Update(id uint, req interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/buildings/%d", id)

	resp, err := s.client.doRequest("PUT", endpoint, req)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Delete removes a building by ID
func (s *BuildingsService) Delete(id uint) error {
	endpoint := fmt.Sprintf("/api/v1/buildings/%d", id)
	resp, err := s.client.doRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	_, err = s.client.handleRawResponse(resp)
	return err
}
