package types

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	// DefaultPage is the default page number
	DefaultPage = 1
	// DefaultPageSize is the default number of items per page
	DefaultPageSize = 20
	// MaxPageSize is the maximum allowed page size
	MaxPageSize = 100
)

// PaginationParams represents pagination parameters from request
type PaginationParams struct {
	Page     int
	PageSize int
	Sort     string
	Order    string // asc or desc
}

// PaginationMeta represents pagination metadata in response
type PaginationMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse[T any] struct {
	Data []T            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// ExtractPaginationParams extracts pagination parameters from Gin context
func ExtractPaginationParams(c *gin.Context) PaginationParams {
	page := parseIntParam(c.Query("page"), DefaultPage)
	pageSize := parseIntParam(c.Query("pageSize"), DefaultPageSize)

	// Enforce max page size
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	// Ensure minimum values
	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}

	sort := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "asc")

	// Validate order
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Sort:     sort,
		Order:    order,
	}
}

// NewPaginatedResponse creates a new paginated response
func NewPaginatedResponse[T any](data []T, params PaginationParams, totalItems int) *PaginatedResponse[T] {
	totalPages := int(math.Ceil(float64(totalItems) / float64(params.PageSize)))

	return &PaginatedResponse[T]{
		Data: data,
		Meta: PaginationMeta{
			Page:       params.Page,
			PageSize:   params.PageSize,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}
}

// Offset calculates the database offset for the current page
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit returns the page size
func (p PaginationParams) Limit() int {
	return p.PageSize
}

// parseIntParam parses an integer parameter with a default value
func parseIntParam(param string, defaultValue int) int {
	if param == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(param)
	if err != nil {
		return defaultValue
	}

	return value
}
