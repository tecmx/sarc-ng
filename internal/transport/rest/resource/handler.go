package resource

import (
	"net/http"
	"sarc-ng/internal/domain/resource"
	"sarc-ng/internal/transport/common"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for resource operations
type Handler struct {
	*common.BaseHandler[resource.Resource, CreateResourceDTO, UpdateResourceDTO, ResourceDTO]
	service resource.Usecase
	mapper  *Mapper
}

// NewHandler creates a new resource handler
func NewHandler(service resource.Usecase) *Handler {
	mapper := NewMapper()
	baseHandler := common.NewBaseHandler[resource.Resource, CreateResourceDTO, UpdateResourceDTO, ResourceDTO](
		"resource")
	return &Handler{
		BaseHandler: baseHandler,
		service:     service,
		mapper:      mapper,
	}
}

// GetAll retrieves all resources with pagination
// @Summary Get all resources (paginated)
// @Description Retrieve a paginated list of all resources with support for search, filtering, and sorting
// @Tags resources
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param page_size query int false "Page size (default: 20, max: 100)" minimum(1) maximum(100)
// @Param search query string false "Search across name and description fields"
// @Param filter[name] query string false "Filter by exact name match"
// @Param filter[type] query string false "Filter by exact type match"
// @Param filter[location] query string false "Filter by exact location match"
// @Param filter[is_available] query boolean false "Filter by availability status"
// @Param sort query string false "Sort field (id, name, type, location, is_available, created_at, updated_at)"
// @Param order query string false "Sort order (asc or desc)" default(asc)
// @Success 200 {object} common.PaginatedResponse[ResourceDTO] "Paginated list of resources"
// @Failure 400 {object} common.ErrorResponse "Invalid filter or sort field"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /resources [get]
func (h *Handler) GetAll(c *gin.Context) {
	params := common.GetPaginationParams(c)
	entities, pagination, err := h.service.GetAllResources(c.Request.Context(), params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve "+h.GetEntityName()+"s", err.Error())
		return
	}

	dtos := make([]ResourceDTO, len(entities))
	for i, entity := range entities {
		dtos[i] = *h.mapper.FromDomain(&entity)
	}

	response := common.PaginatedResponse[ResourceDTO]{
		Data:       dtos,
		Pagination: pagination,
	}
	c.JSON(http.StatusOK, response)
}

// GetByID retrieves a resource by ID
// @Summary Get resource by ID
// @Description Retrieve a specific resource by its unique identifier
// @Tags resources
// @Accept json
// @Produce json
// @Param id path int true "Resource ID" minimum(1)
// @Success 200 {object} ResourceDTO "Resource details"
// @Failure 400 {object} common.ErrorResponse "Invalid resource ID"
// @Failure 404 {object} common.ErrorResponse "Resource not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /resources/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := common.ParseIDFromPath(c, h.GetEntityName())
	if err != nil {
		return
	}

	entity, err := h.service.GetResource(c.Request.Context(), id)
	if err != nil {
		common.HandleError(c, err, "Failed to retrieve "+h.GetEntityName())
		return
	}
	if entity == nil {
		common.RespondWithError(c, http.StatusNotFound, h.GetEntityName()+" not found", "No "+h.GetEntityName()+" exists with the specified ID")
		return
	}

	dto := h.mapper.FromDomain(entity)
	c.JSON(http.StatusOK, dto)
}

// Create creates a new resource
// @Summary Create a new resource
// @Description Create a new resource with name, type, and availability information
// @Tags resources
// @Accept json
// @Produce json
// @Param resource body CreateResourceDTO true "Resource creation data"
// @Success 201 {object} ResourceDTO "Created resource"
// @Failure 400 {object} common.ErrorResponse "Invalid input data"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /resources [post]
func (h *Handler) Create(c *gin.Context) {
	createDTO, err := h.BindCreateJSON(c)
	if err != nil {
		return
	}

	entity := h.mapper.ToDomain(createDTO)
	if err := h.service.CreateResource(c.Request.Context(), entity); err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to create "+h.GetEntityName(), err.Error())
		return
	}

	createdDTO := h.mapper.FromDomain(entity)
	c.JSON(http.StatusCreated, createdDTO)
}

// Update updates an existing resource
// @Summary Update an existing resource
// @Description Update an existing resource's name, type, and availability by ID
// @Tags resources
// @Accept json
// @Produce json
// @Param id path int true "Resource ID" minimum(1)
// @Param resource body UpdateResourceDTO true "Resource update data"
// @Success 200 {object} ResourceDTO "Updated resource"
// @Failure 400 {object} common.ErrorResponse "Invalid input data"
// @Failure 404 {object} common.ErrorResponse "Resource not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /resources/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, updateDTO, err := h.ParseIDAndBindJSON(c)
	if err != nil {
		return
	}

	entity := h.mapper.ToDomainWithID(updateDTO, id)
	if err := h.service.UpdateResource(c.Request.Context(), entity); err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to update "+h.GetEntityName(), err.Error())
		return
	}

	updatedDTO := h.mapper.FromDomain(entity)
	c.JSON(http.StatusOK, updatedDTO)
}

// Delete removes a resource
// @Summary Delete a resource
// @Description Delete a resource by its ID
// @Tags resources
// @Accept json
// @Produce json
// @Param id path int true "Resource ID" minimum(1)
// @Success 200 {object} common.SuccessResponse "Resource deleted successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid resource ID"
// @Failure 404 {object} common.ErrorResponse "Resource not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /resources/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := common.ParseIDFromPath(c, h.GetEntityName())
	if err != nil {
		return
	}

	if err := h.service.DeleteResource(c.Request.Context(), id); err != nil {
		common.HandleError(c, err, "Failed to delete "+h.GetEntityName())
		return
	}

	common.RespondWithSuccess(c, http.StatusOK, h.GetEntityName()+" deleted successfully")
}
