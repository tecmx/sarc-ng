package common

import (
	"sarc-ng/internal/domain/common"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// PaginatedResponse wraps the data with pagination metadata
type PaginatedResponse[T any] struct {
	Data       []T                        `json:"data"`
	Pagination common.PaginationResult `json:"pagination"`
}

// GetPaginationParams extracts pagination parameters from query string
func GetPaginationParams(c *gin.Context) common.PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Extract search parameter
	search := c.Query("search")

	// Extract filters from filter[field] syntax
	filters := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if strings.HasPrefix(key, "filter[") && strings.HasSuffix(key, "]") {
			// Extract field name from filter[field] syntax
			field := key[7 : len(key)-1]
			if len(values) > 0 && field != "" {
				filters[field] = values[0]
			}
		}
	}

	// Extract sort and order parameters
	sort := c.Query("sort")
	order := c.DefaultQuery("order", "asc")

	return common.NewPaginationParams(page, pageSize, search, filters, sort, order)
}

