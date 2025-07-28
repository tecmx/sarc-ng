package common

import "errors"

// Domain errors that can be checked by type across layers
var (
	// ErrNotFound indicates that a requested entity was not found
	ErrNotFound = errors.New("entity not found")

	// ErrInvalidInput indicates that provided input is invalid
	ErrInvalidInput = errors.New("invalid input")

	// ErrConflict indicates that the operation conflicts with existing data
	ErrConflict = errors.New("resource conflict")

	// ErrUnauthorized indicates that the operation is not authorized
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden indicates that the operation is forbidden
	ErrForbidden = errors.New("forbidden")
)

// IsNotFoundError checks if an error is a "not found" error
func IsNotFoundError(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsInvalidInputError checks if an error is an "invalid input" error
func IsInvalidInputError(err error) bool {
	return errors.Is(err, ErrInvalidInput)
}

// IsConflictError checks if an error is a "conflict" error
func IsConflictError(err error) bool {
	return errors.Is(err, ErrConflict)
}

// IsUnauthorizedError checks if an error is an "unauthorized" error
func IsUnauthorizedError(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsForbiddenError checks if an error is a "forbidden" error
func IsForbiddenError(err error) bool {
	return errors.Is(err, ErrForbidden)
}
