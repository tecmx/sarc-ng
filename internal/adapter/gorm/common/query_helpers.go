package common

import (
	domainCommon "sarc-ng/internal/domain/common"

	"gorm.io/gorm"
)

// ApplySearch adds LIKE search conditions for specified fields
// Supports optional case-insensitive search for better UX
// Note: Uses GORM parameterized queries to prevent SQL injection
func ApplySearch(query *gorm.DB, searchTerm string, fields []string, caseInsensitive bool) *gorm.DB {
	if searchTerm == "" {
		return query
	}

	searchPattern := "%" + searchTerm + "%"
	searchQuery := query.Session(&gorm.Session{})

	for i, field := range fields {
		condition := field + " LIKE ?"
		if caseInsensitive {
			condition = "LOWER(" + field + ") LIKE LOWER(?)"
		}

		if i == 0 {
			searchQuery = searchQuery.Where(condition, searchPattern)
		} else {
			searchQuery = searchQuery.Or(condition, searchPattern)
		}
	}

	return query.Where(searchQuery)
}

// ApplyFilters adds WHERE conditions for each filter
// Field names come from validated whitelists, preventing SQL injection
func ApplyFilters(query *gorm.DB, filters map[string]string) *gorm.DB {
	for field, value := range filters {
		query = query.Where(field+" = ?", value)
	}
	return query
}

// ApplySorting adds ORDER BY clause
// Field names validated before calling, safe from injection
func ApplySorting(query *gorm.DB, sortField, sortOrder string) *gorm.DB {
	if sortField != "" {
		query = query.Order(sortField + " " + sortOrder)
	}
	return query
}

// BuildFilteredQuery is a convenience wrapper that applies all filters at once
// Use this for standard filtering, or call individual functions for custom logic
func BuildFilteredQuery(
	query *gorm.DB,
	params domainCommon.PaginationParams,
	searchFields []string,
	caseInsensitive bool,
) *gorm.DB {
	query = ApplySearch(query, params.Search, searchFields, caseInsensitive)
	query = ApplyFilters(query, params.Filters)
	query = ApplySorting(query, params.Sort, params.Order)
	return query
}

// CountAndPaginate executes count and applies pagination
// Returns modified query and total count for building pagination metadata
func CountAndPaginate(query *gorm.DB, offset, limit int) (*gorm.DB, int64, error) {
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return query.Offset(offset).Limit(limit), total, nil
}

