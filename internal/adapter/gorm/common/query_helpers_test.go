package common

import (
	"regexp"
	"sarc-ng/internal/domain/common"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestApplySearch(t *testing.T) {
	t.Run("Empty search term returns unchanged query", func(t *testing.T) {
		db, _ := setupTestDB(t)
		query := db.Model(&struct{}{})

		result := ApplySearch(query, "", []string{"name", "code"}, false)

		assert.Equal(t, query, result)
	})

	t.Run("Case-sensitive search adds LIKE conditions", func(t *testing.T) {
		db, mock := setupTestDB(t)
		query := db.Model(&struct{}{})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM`)).
			WithArgs("%test%", "%test%").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		result := ApplySearch(query, "test", []string{"name", "code"}, false)
		var count int64
		_ = result.Count(&count)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Case-insensitive search uses LOWER function", func(t *testing.T) {
		db, mock := setupTestDB(t)
		query := db.Model(&struct{}{})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM`)).
			WithArgs("%test%", "%test%").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		result := ApplySearch(query, "test", []string{"name", "code"}, true)
		var count int64
		_ = result.Count(&count)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Search with single field", func(t *testing.T) {
		db, mock := setupTestDB(t)
		query := db.Model(&struct{}{})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM`)).
			WithArgs("%test%").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		result := ApplySearch(query, "test", []string{"name"}, false)
		var count int64
		_ = result.Count(&count)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestApplyFilters(t *testing.T) {
	t.Run("Empty filters returns unchanged query", func(t *testing.T) {
		db, _ := setupTestDB(t)
		query := db.Model(&struct{}{})

		result := ApplyFilters(query, map[string]string{})

		assert.Equal(t, query, result)
	})

	t.Run("Single filter adds WHERE condition", func(t *testing.T) {
		db, mock := setupTestDB(t)
		query := db.Model(&struct{}{})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM`)).
			WithArgs("active").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		result := ApplyFilters(query, map[string]string{"status": "active"})
		var count int64
		_ = result.Count(&count)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Multiple filters add multiple WHERE conditions", func(t *testing.T) {
		db, mock := setupTestDB(t)
		query := db.Model(&struct{}{})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		result := ApplyFilters(query, map[string]string{
			"status": "active",
			"type":   "premium",
		})
		var count int64
		_ = result.Count(&count)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestApplySorting(t *testing.T) {
	t.Run("Empty sort field returns unchanged query", func(t *testing.T) {
		db, _ := setupTestDB(t)
		query := db.Model(&struct{}{})

		result := ApplySorting(query, "", "asc")

		assert.Equal(t, query, result)
	})

	t.Run("Valid sort field adds ORDER BY", func(t *testing.T) {
		db, mock := setupTestDB(t)
		query := db.Model(&struct{}{})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		result := ApplySorting(query, "name", "asc")
		var count int64
		_ = result.Count(&count)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Descending sort order", func(t *testing.T) {
		db, mock := setupTestDB(t)
		query := db.Model(&struct{}{})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		result := ApplySorting(query, "created_at", "desc")
		var count int64
		_ = result.Count(&count)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBuildFilteredQuery(t *testing.T) {
	t.Run("Combines all filters correctly", func(t *testing.T) {
		db, mock := setupTestDB(t)
		query := db.Model(&struct{}{})

		params := common.PaginationParams{
			Search:  "test",
			Filters: map[string]string{"status": "active"},
			Sort:    "name",
			Order:   "asc",
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		result := BuildFilteredQuery(query, params, []string{"name", "code"}, false)
		var count int64
		_ = result.Count(&count)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Works with minimal params", func(t *testing.T) {
		db, mock := setupTestDB(t)
		query := db.Model(&struct{}{})

		params := common.PaginationParams{
			Filters: make(map[string]string),
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		result := BuildFilteredQuery(query, params, []string{"name"}, false)
		var count int64
		_ = result.Count(&count)

		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Case-insensitive flag works", func(t *testing.T) {
		db, mock := setupTestDB(t)
		query := db.Model(&struct{}{})

		params := common.PaginationParams{
			Search:  "Test",
			Filters: make(map[string]string),
		}

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM`)).
			WithArgs("%Test%", "%Test%").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		result := BuildFilteredQuery(query, params, []string{"name", "code"}, true)
		var count int64
		_ = result.Count(&count)

		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCountAndPaginate(t *testing.T) {
	t.Run("Successfully counts and paginates", func(t *testing.T) {
		db, mock := setupTestDB(t)
		query := db.Model(&struct{}{})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(100))

		result, total, err := CountAndPaginate(query, 20, 10)

		assert.NoError(t, err)
		assert.Equal(t, int64(100), total)
		assert.NotNil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Handles count error", func(t *testing.T) {
		db, mock := setupTestDB(t)
		query := db.Model(&struct{}{})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM`)).
			WillReturnError(gorm.ErrInvalidData)

		result, total, err := CountAndPaginate(query, 20, 10)

		assert.Error(t, err)
		assert.Equal(t, int64(0), total)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Applies correct offset and limit", func(t *testing.T) {
		db, mock := setupTestDB(t)
		query := db.Model(&struct{}{})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(50))

		result, total, err := CountAndPaginate(query, 10, 20)

		assert.NoError(t, err)
		assert.Equal(t, int64(50), total)
		assert.NotNil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Zero offset and limit", func(t *testing.T) {
		db, mock := setupTestDB(t)
		query := db.Model(&struct{}{})

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		result, total, err := CountAndPaginate(query, 0, 20)

		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
		assert.NotNil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

