package client

import "fmt"

// ResourcesService provides methods for resource operations
type ResourcesService struct {
	client *Client
}

// Resources returns the resources service
func (c *Client) Resources() *ResourcesService {
	return &ResourcesService{client: c}
}

// List retrieves all resources with pagination
func (s *ResourcesService) List(page, pageSize int) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/resources?page=%d&pageSize=%d", page, pageSize)
	resp, err := s.client.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Get retrieves a specific resource by ID
func (s *ResourcesService) Get(id uint) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/resources/%d", id)
	resp, err := s.client.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Create creates a new resource
func (s *ResourcesService) Create(req interface{}) ([]byte, error) {
	resp, err := s.client.doRequest("POST", "/api/v1/resources", req)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Update updates an existing resource
func (s *ResourcesService) Update(id uint, req interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/resources/%d", id)
	resp, err := s.client.doRequest("PUT", endpoint, req)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Delete removes a resource by ID
func (s *ResourcesService) Delete(id uint) error {
	endpoint := fmt.Sprintf("/api/v1/resources/%d", id)
	resp, err := s.client.doRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	_, err = s.client.handleRawResponse(resp)
	return err
}
