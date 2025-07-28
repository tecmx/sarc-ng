package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// validator instance for validating structs
var validate = validator.New()

// BaseHandler provides common utilities for REST handlers
type BaseHandler[Entity any, CreateDTO any, UpdateDTO any, ResponseDTO any] struct {
	entityName string
}

// NewBaseHandler creates a new base handler
func NewBaseHandler[Entity any, CreateDTO any, UpdateDTO any, ResponseDTO any](
	entityName string,
) *BaseHandler[Entity, CreateDTO, UpdateDTO, ResponseDTO] {
	return &BaseHandler[Entity, CreateDTO, UpdateDTO, ResponseDTO]{
		entityName: entityName,
	}
}

// GetEntityName returns the entity name for error messages
func (h *BaseHandler[Entity, CreateDTO, UpdateDTO, ResponseDTO]) GetEntityName() string {
	return h.entityName
}

// ParseIDAndBindJSON is a helper method that parses ID from path and binds JSON for update operations
func (h *BaseHandler[Entity, CreateDTO, UpdateDTO, ResponseDTO]) ParseIDAndBindJSON(c *gin.Context) (uint, *UpdateDTO, error) {
	id, err := ParseIDFromPath(c, h.entityName)
	if err != nil {
		return 0, nil, err
	}

	var updateDTO UpdateDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid JSON format", err.Error())
		return 0, nil, err
	}

	// Validate the DTO using struct tags
	if err := validate.Struct(&updateDTO); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Validation failed", err.Error())
		return 0, nil, err
	}

	return id, &updateDTO, nil
}

// BindCreateJSON is a helper method that binds JSON for create operations
func (h *BaseHandler[Entity, CreateDTO, UpdateDTO, ResponseDTO]) BindCreateJSON(c *gin.Context) (*CreateDTO, error) {
	var createDTO CreateDTO
	if err := c.ShouldBindJSON(&createDTO); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Invalid JSON format", err.Error())
		return nil, err
	}

	// Validate the DTO using struct tags
	if err := validate.Struct(&createDTO); err != nil {
		RespondWithError(c, http.StatusBadRequest, "Validation failed", err.Error())
		return nil, err
	}

	return &createDTO, nil
}
