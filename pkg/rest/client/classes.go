package client

import "fmt"

// ClassesService provides methods for class operations
type ClassesService struct {
	client *Client
}

// Classes returns the classes service
func (c *Client) Classes() *ClassesService {
	return &ClassesService{client: c}
}

// List retrieves all classes with pagination
func (s *ClassesService) List(page, pageSize int) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/classes?page=%d&pageSize=%d", page, pageSize)
	resp, err := s.client.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Get retrieves a specific class by ID
func (s *ClassesService) Get(id uint) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/classes/%d", id)
	resp, err := s.client.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Create creates a new class
func (s *ClassesService) Create(req interface{}) ([]byte, error) {
	resp, err := s.client.doRequest("POST", "/api/v1/classes", req)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Update updates an existing class
func (s *ClassesService) Update(id uint, req interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/classes/%d", id)
	resp, err := s.client.doRequest("PUT", endpoint, req)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Delete removes a class by ID
func (s *ClassesService) Delete(id uint) error {
	endpoint := fmt.Sprintf("/api/v1/classes/%d", id)
	resp, err := s.client.doRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	_, err = s.client.handleRawResponse(resp)
	return err
}
