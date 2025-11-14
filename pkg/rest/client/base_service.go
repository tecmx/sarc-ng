package client

import "fmt"

// baseService provides common CRUD operations for REST resources
type baseService struct {
	client       *Client
	resourceName string
}

// newBaseService creates a new base service for a given resource
func newBaseService(client *Client, resourceName string) *baseService {
	return &baseService{
		client:       client,
		resourceName: resourceName,
	}
}

// list retrieves all resources with pagination
func (s *baseService) list(page, pageSize int) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/%s?page=%d&pageSize=%d", s.resourceName, page, pageSize)
	resp, err := s.client.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// get retrieves a specific resource by ID
func (s *baseService) get(id uint) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/%s/%d", s.resourceName, id)
	resp, err := s.client.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// create creates a new resource
func (s *baseService) create(req interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/%s", s.resourceName)
	resp, err := s.client.doRequest("POST", endpoint, req)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// update updates an existing resource
func (s *baseService) update(id uint, req interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/%s/%d", s.resourceName, id)
	resp, err := s.client.doRequest("PUT", endpoint, req)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// delete removes a resource by ID
func (s *baseService) delete(id uint) error {
	endpoint := fmt.Sprintf("/api/v1/%s/%d", s.resourceName, id)
	resp, err := s.client.doRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	_, err = s.client.handleRawResponse(resp)
	return err
}

