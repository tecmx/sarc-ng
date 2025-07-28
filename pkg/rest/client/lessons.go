package client

import "fmt"

// LessonsService provides methods for lesson operations
type LessonsService struct {
	client *Client
}

// Lessons returns the lessons service
func (c *Client) Lessons() *LessonsService {
	return &LessonsService{client: c}
}

// List retrieves all lessons with pagination
func (s *LessonsService) List(page, pageSize int) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/lessons?page=%d&pageSize=%d", page, pageSize)
	resp, err := s.client.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Get retrieves a specific lesson by ID
func (s *LessonsService) Get(id uint) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/lessons/%d", id)
	resp, err := s.client.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Create creates a new lesson
func (s *LessonsService) Create(req interface{}) ([]byte, error) {
	resp, err := s.client.doRequest("POST", "/api/v1/lessons", req)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Update updates an existing lesson
func (s *LessonsService) Update(id uint, req interface{}) ([]byte, error) {
	endpoint := fmt.Sprintf("/api/v1/lessons/%d", id)
	resp, err := s.client.doRequest("PUT", endpoint, req)
	if err != nil {
		return nil, err
	}

	return s.client.handleRawResponse(resp)
}

// Delete removes a lesson by ID
func (s *LessonsService) Delete(id uint) error {
	endpoint := fmt.Sprintf("/api/v1/lessons/%d", id)
	resp, err := s.client.doRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}

	_, err = s.client.handleRawResponse(resp)
	return err
}
