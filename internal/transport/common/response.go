package common

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid input"`
	Message string `json:"message,omitempty" example:"The provided data is invalid"`
	Code    int    `json:"code" example:"400"`
}

// SuccessResponse represents a standard success response
type SuccessResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
	Data    any    `json:"data,omitempty"`
}

// ParseIDFromPath extracts and validates ID from URL path parameter
func ParseIDFromPath(c *gin.Context, entityName string) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid " + entityName + " ID",
			Message: entityName + " ID must be a positive integer",
			Code:    http.StatusBadRequest,
		})
		return 0, err
	}
	return uint(id), nil
}

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, statusCode int, error string, message string) {
	c.JSON(statusCode, ErrorResponse{
		Error:   error,
		Message: message,
		Code:    statusCode,
	})
}

// RespondWithSuccess sends a standardized success response
func RespondWithSuccess(c *gin.Context, statusCode int, message string, data ...any) {
	response := SuccessResponse{Message: message}
	if len(data) > 0 {
		response.Data = data[0]
	}
	c.JSON(statusCode, response)
}
