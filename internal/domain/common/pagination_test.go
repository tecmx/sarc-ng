package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPaginationParams(t *testing.T) {
	t.Run("Valid parameters", func(t *testing.T) {
		params := NewPaginationParams(2, 50, "", nil, "", "asc")
		assert.Equal(t, 2, params.Page)
		assert.Equal(t, 50, params.PageSize)
		assert.Equal(t, 50, params.Offset())
		assert.Equal(t, 50, params.Limit())
		assert.Equal(t, "", params.Search)
		assert.NotNil(t, params.Filters)
		assert.Equal(t, "", params.Sort)
		assert.Equal(t, "asc", params.Order)
	})

	t.Run("Default page when zero", func(t *testing.T) {
		params := NewPaginationParams(0, 20, "", nil, "", "asc")
		assert.Equal(t, 1, params.Page)
		assert.Equal(t, 20, params.PageSize)
		assert.Equal(t, 0, params.Offset())
	})

	t.Run("Default page size when zero", func(t *testing.T) {
		params := NewPaginationParams(1, 0, "", nil, "", "asc")
		assert.Equal(t, 1, params.Page)
		assert.Equal(t, 20, params.PageSize)
	})

	t.Run("Max page size enforced", func(t *testing.T) {
		params := NewPaginationParams(1, 150, "", nil, "", "asc")
		assert.Equal(t, 1, params.Page)
		assert.Equal(t, 100, params.PageSize)
	})

	t.Run("Negative values default", func(t *testing.T) {
		params := NewPaginationParams(-5, -10, "", nil, "", "asc")
		assert.Equal(t, 1, params.Page)
		assert.Equal(t, 20, params.PageSize)
	})

	t.Run("With search and filters", func(t *testing.T) {
		filters := map[string]string{"status": "active"}
		params := NewPaginationParams(1, 20, "test", filters, "name", "desc")
		assert.Equal(t, "test", params.Search)
		assert.Equal(t, "active", params.Filters["status"])
		assert.Equal(t, "name", params.Sort)
		assert.Equal(t, "desc", params.Order)
	})

	t.Run("Invalid order defaults to asc", func(t *testing.T) {
		params := NewPaginationParams(1, 20, "", nil, "", "invalid")
		assert.Equal(t, "asc", params.Order)
	})
}

func TestPaginationParams_Offset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected int
	}{
		{"First page", 1, 20, 0},
		{"Second page", 2, 20, 20},
		{"Third page", 3, 50, 100},
		{"Page 10", 10, 10, 90},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := PaginationParams{Page: tt.page, PageSize: tt.pageSize}
			assert.Equal(t, tt.expected, params.Offset())
		})
	}
}

func TestNewPaginationResult(t *testing.T) {
	t.Run("Exact division", func(t *testing.T) {
		result := NewPaginationResult(2, 20, 100)
		assert.Equal(t, 2, result.Page)
		assert.Equal(t, 20, result.PageSize)
		assert.Equal(t, int64(100), result.TotalItems)
		assert.Equal(t, 5, result.TotalPages)
	})

	t.Run("With remainder", func(t *testing.T) {
		result := NewPaginationResult(1, 20, 95)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 20, result.PageSize)
		assert.Equal(t, int64(95), result.TotalItems)
		assert.Equal(t, 5, result.TotalPages) // Ceiling of 95/20
	})

	t.Run("Single page", func(t *testing.T) {
		result := NewPaginationResult(1, 20, 15)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 20, result.PageSize)
		assert.Equal(t, int64(15), result.TotalItems)
		assert.Equal(t, 1, result.TotalPages)
	})

	t.Run("Empty result", func(t *testing.T) {
		result := NewPaginationResult(1, 20, 0)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 20, result.PageSize)
		assert.Equal(t, int64(0), result.TotalItems)
		assert.Equal(t, 0, result.TotalPages)
	})
}
