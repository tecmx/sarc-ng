package lesson

import "sarc-ng/internal/domain/common"

// SearchableFields defines fields that can be searched using the search parameter
var SearchableFields = []string{"repeat_pattern"}

var fieldValidator = common.NewFieldValidator(
	"lessons",
	[]string{"class_id", "resource_id", "repeat_pattern"},
	[]string{"id", "class_id", "resource_id", "start_time", "end_time", "repeat_pattern", "created_at", "updated_at"},
)

// ValidateFilterField checks if a field name is valid for filtering
func ValidateFilterField(field string) error {
	return fieldValidator.ValidateFilterField(field)
}

// ValidateSortField checks if a field name is valid for sorting
func ValidateSortField(field string) error {
	return fieldValidator.ValidateSortField(field)
}

// ValidateParams validates all pagination parameters
func ValidateParams(params common.PaginationParams) error {
	return fieldValidator.ValidateParams(params)
}

