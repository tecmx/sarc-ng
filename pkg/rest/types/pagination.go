package types

// PaginationMeta represents pagination metadata for REST responses
type PaginationMeta struct {
	Total       int  `json:"total" example:"100"`
	Page        int  `json:"page" example:"1"`
	PageSize    int  `json:"pageSize" example:"10"`
	TotalPages  int  `json:"totalPages" example:"10"`
	HasNext     bool `json:"hasNext" example:"true"`
	HasPrevious bool `json:"hasPrevious" example:"false"`
}

// PaginatedResponse represents a generic paginated response envelope
type PaginatedResponse[T any] struct {
	Data       []T            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// SearchRequest represents common search parameters
type SearchRequest struct {
	Query    string `json:"query,omitempty" form:"query"`
	Page     int    `json:"page,omitempty" form:"page" minimum:"1"`
	PageSize int    `json:"pageSize,omitempty" form:"pageSize" minimum:"1" maximum:"100"`
}

// CalculatePaginationMeta calculates pagination metadata
func CalculatePaginationMeta(total, page, pageSize int) PaginationMeta {
	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationMeta{
		Total:       total,
		Page:        page,
		PageSize:    pageSize,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}
}
