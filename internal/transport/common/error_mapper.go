package common

import (
	"net/http"
	domainCommon "sarc-ng/internal/domain/common"

	"github.com/gin-gonic/gin"
)

// ErrorMapping defines how domain errors map to HTTP responses
type ErrorMapping struct {
	StatusCode int
	Message    string
}

// ErrorChecker defines a function that checks if an error matches a domain error type
type ErrorChecker func(error) bool

// errorMappings defines the mapping from domain error checkers to HTTP responses
var errorMappings = []struct {
	Checker ErrorChecker
	Mapping ErrorMapping
}{
	{domainCommon.IsNotFoundError, ErrorMapping{http.StatusNotFound, "Resource not found"}},
	{domainCommon.IsInvalidInputError, ErrorMapping{http.StatusBadRequest, "Invalid input provided"}},
	{domainCommon.IsConflictError, ErrorMapping{http.StatusConflict, "Resource conflict"}},
	{domainCommon.IsUnauthorizedError, ErrorMapping{http.StatusUnauthorized, "Unauthorized access"}},
	{domainCommon.IsForbiddenError, ErrorMapping{http.StatusForbidden, "Access forbidden"}},
}

// HandleError automatically maps domain errors to appropriate HTTP responses
func HandleError(c *gin.Context, err error, defaultMessage string) {
	if err == nil {
		return
	}

	// Check for mapped domain errors using checker functions
	for _, mapping := range errorMappings {
		if mapping.Checker(err) {
			RespondWithError(c, mapping.Mapping.StatusCode, mapping.Mapping.Message, err.Error())
			return
		}
	}

	// Default to internal server error for unmapped errors
	RespondWithError(c, http.StatusInternalServerError, defaultMessage, err.Error())
}
