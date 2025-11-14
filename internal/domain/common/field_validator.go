package common

import "fmt"

// FieldValidator validates field names against whitelists
// Uses map for O(1) lookup performance vs O(n) array iteration
type FieldValidator struct {
	filterable map[string]bool
	sortable   map[string]bool
	entityName string
}

// NewFieldValidator creates a validator with field whitelists
func NewFieldValidator(entity string, filterFields, sortFields []string) *FieldValidator {
	v := &FieldValidator{
		filterable: make(map[string]bool),
		sortable:   make(map[string]bool),
		entityName: entity,
	}
	for _, f := range filterFields {
		v.filterable[f] = true
	}
	for _, f := range sortFields {
		v.sortable[f] = true
	}
	return v
}

// ValidateFilterField checks if field is filterable
func (v *FieldValidator) ValidateFilterField(field string) error {
	if !v.filterable[field] {
		return fmt.Errorf("%w: filter field '%s' not allowed for %s",
			ErrInvalidInput, field, v.entityName)
	}
	return nil
}

// ValidateSortField checks if field is sortable
func (v *FieldValidator) ValidateSortField(field string) error {
	if field == "" {
		return nil
	}
	if !v.sortable[field] {
		return fmt.Errorf("%w: sort field '%s' not allowed for %s",
			ErrInvalidInput, field, v.entityName)
	}
	return nil
}

// ValidateParams validates both filters and sort in one call
// Simplifies service layer validation from 8 lines to 3 lines
func (v *FieldValidator) ValidateParams(params PaginationParams) error {
	for field := range params.Filters {
		if err := v.ValidateFilterField(field); err != nil {
			return err
		}
	}
	return v.ValidateSortField(params.Sort)
}

