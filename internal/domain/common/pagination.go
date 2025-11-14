package common

// PaginationParams contains parameters for paginated requests
type PaginationParams struct {
	Page     int               // Current page number (1-based)
	PageSize int               // Number of items per page
	Search   string            // General text search across searchable fields
	Filters  map[string]string // Specific field filters (field -> value)
	Sort     string            // Sort field name
	Order    string            // Sort order: "asc" or "desc"
}

// PaginationResult contains pagination metadata
type PaginationResult struct {
	Page       int   `json:"page"`        // Current page number
	PageSize   int   `json:"page_size"`   // Items per page
	TotalItems int64 `json:"total_items"` // Total number of items
	TotalPages int   `json:"total_pages"` // Total number of pages
}

// NewPaginationParams creates pagination parameters with defaults
func NewPaginationParams(page, pageSize int, search string, filters map[string]string, sort, order string) PaginationParams {
	// Default values
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20 // Default page size
	}
	// Max page size to prevent abuse
	if pageSize > 100 {
		pageSize = 100
	}

	// Default sort order
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// Initialize filters map if nil
	if filters == nil {
		filters = make(map[string]string)
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
		Filters:  filters,
		Sort:     sort,
		Order:    order,
	}
}

// Offset calculates the database offset for the current page
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit returns the page size (for clarity in usage)
func (p PaginationParams) Limit() int {
	return p.PageSize
}

// NewPaginationResult creates pagination result metadata
func NewPaginationResult(page, pageSize int, totalItems int64) PaginationResult {
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		totalPages++
	}

	return PaginationResult{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
