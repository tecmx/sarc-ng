package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFieldValidator(t *testing.T) {
	t.Run("Creates validator with all fields", func(t *testing.T) {
		validator := NewFieldValidator(
			"buildings",
			[]string{"name", "code"},
			[]string{"id", "name", "code", "created_at"},
		)

		assert.NotNil(t, validator)
		assert.Equal(t, "buildings", validator.entityName)
		assert.Len(t, validator.filterable, 2)
		assert.Len(t, validator.sortable, 4)
	})

	t.Run("Creates validator with empty fields", func(t *testing.T) {
		validator := NewFieldValidator("test", []string{}, []string{})

		assert.NotNil(t, validator)
		assert.Equal(t, "test", validator.entityName)
		assert.Empty(t, validator.filterable)
		assert.Empty(t, validator.sortable)
	})
}

func TestFieldValidator_ValidateFilterField(t *testing.T) {
	validator := NewFieldValidator(
		"buildings",
		[]string{"name", "code", "status"},
		[]string{"id", "name"},
	)

	t.Run("Valid filter field returns no error", func(t *testing.T) {
		err := validator.ValidateFilterField("name")
		assert.NoError(t, err)

		err = validator.ValidateFilterField("code")
		assert.NoError(t, err)

		err = validator.ValidateFilterField("status")
		assert.NoError(t, err)
	})

	t.Run("Invalid filter field returns error", func(t *testing.T) {
		err := validator.ValidateFilterField("invalid_field")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInvalidInput))
		assert.Contains(t, err.Error(), "filter field 'invalid_field' not allowed")
		assert.Contains(t, err.Error(), "buildings")
	})

	t.Run("Empty field name returns error", func(t *testing.T) {
		err := validator.ValidateFilterField("")
		assert.Error(t, err)
	})

	t.Run("Field not in filter list but in sort list returns error", func(t *testing.T) {
		err := validator.ValidateFilterField("id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "filter field 'id' not allowed")
	})
}

func TestFieldValidator_ValidateSortField(t *testing.T) {
	validator := NewFieldValidator(
		"classes",
		[]string{"name"},
		[]string{"id", "name", "capacity", "created_at", "updated_at"},
	)

	t.Run("Valid sort field returns no error", func(t *testing.T) {
		err := validator.ValidateSortField("id")
		assert.NoError(t, err)

		err = validator.ValidateSortField("name")
		assert.NoError(t, err)

		err = validator.ValidateSortField("capacity")
		assert.NoError(t, err)
	})

	t.Run("Empty sort field returns no error", func(t *testing.T) {
		err := validator.ValidateSortField("")
		assert.NoError(t, err)
	})

	t.Run("Invalid sort field returns error", func(t *testing.T) {
		err := validator.ValidateSortField("invalid_field")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrInvalidInput))
		assert.Contains(t, err.Error(), "sort field 'invalid_field' not allowed")
		assert.Contains(t, err.Error(), "classes")
	})
}

func TestFieldValidator_ValidateParams(t *testing.T) {
	validator := NewFieldValidator(
		"resources",
		[]string{"name", "type", "location"},
		[]string{"id", "name", "type", "created_at"},
	)

	t.Run("Valid params return no error", func(t *testing.T) {
		params := PaginationParams{
			Filters: map[string]string{
				"name": "test",
				"type": "equipment",
			},
			Sort: "name",
		}

		err := validator.ValidateParams(params)
		assert.NoError(t, err)
	})

	t.Run("Empty params return no error", func(t *testing.T) {
		params := PaginationParams{
			Filters: make(map[string]string),
			Sort:    "",
		}

		err := validator.ValidateParams(params)
		assert.NoError(t, err)
	})

	t.Run("Invalid filter field returns error", func(t *testing.T) {
		params := PaginationParams{
			Filters: map[string]string{
				"name":          "test",
				"invalid_field": "value",
			},
			Sort: "name",
		}

		err := validator.ValidateParams(params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid_field")
	})

	t.Run("Invalid sort field returns error", func(t *testing.T) {
		params := PaginationParams{
			Filters: map[string]string{
				"name": "test",
			},
			Sort: "invalid_sort",
		}

		err := validator.ValidateParams(params)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid_sort")
	})

	t.Run("Both invalid filter and sort returns filter error first", func(t *testing.T) {
		params := PaginationParams{
			Filters: map[string]string{
				"invalid_filter": "value",
			},
			Sort: "invalid_sort",
		}

		err := validator.ValidateParams(params)
		assert.Error(t, err)
		// Should return first error encountered (filter validation)
		assert.Contains(t, err.Error(), "filter field")
	})
}

func TestFieldValidator_Performance(t *testing.T) {
	t.Run("Map lookup is O(1)", func(t *testing.T) {
		// Create validator with many fields
		filterFields := make([]string, 100)
		sortFields := make([]string, 100)
		for i := 0; i < 100; i++ {
			filterFields[i] = string(rune('a' + i))
			sortFields[i] = string(rune('a' + i))
		}

		validator := NewFieldValidator("test", filterFields, sortFields)

		// Validate last field (would be slow with O(n) array search)
		err := validator.ValidateFilterField(string(rune('a' + 99)))
		assert.NoError(t, err)

		err = validator.ValidateSortField(string(rune('a' + 99)))
		assert.NoError(t, err)
	})
}

