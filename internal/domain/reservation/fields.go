package reservation

import "sarc-ng/internal/domain/common"

// SearchableFields defines fields that can be searched using the search parameter
var SearchableFields = []string{"purpose", "user_name"}

var fieldValidator = common.NewFieldValidator(
	"reservations",
	[]string{"resource_id", "user_id", "status", "purpose"},
	[]string{"id", "resource_id", "user_id", "user_name", "start_time", "end_time", "status", "created_at", "updated_at"},
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

