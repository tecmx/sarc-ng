package class

import (
	"net/http"
	"sarc-ng/internal/domain/class"
	"sarc-ng/internal/transport/common"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for class operations
type Handler struct {
	*common.BaseHandler[class.Class, CreateClassDTO, UpdateClassDTO, ClassDTO]
	service class.Usecase
	mapper  *Mapper
}

// NewHandler creates a new class handler
func NewHandler(service class.Usecase) *Handler {
	mapper := NewMapper()
	baseHandler := common.NewBaseHandler[class.Class, CreateClassDTO, UpdateClassDTO, ClassDTO](
		"class")
	return &Handler{
		BaseHandler: baseHandler,
		service:     service,
		mapper:      mapper,
	}
}

// GetAll retrieves all classes with pagination
// @Summary Get all classes (paginated)
// @Description Retrieve a paginated list of all classes with support for search, filtering, and sorting
// @Tags classes
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param page_size query int false "Page size (default: 20, max: 100)" minimum(1) maximum(100)
// @Param search query string false "Search across name field"
// @Param filter[name] query string false "Filter by exact name match"
// @Param filter[capacity] query int false "Filter by exact capacity"
// @Param sort query string false "Sort field (id, name, capacity, created_at, updated_at)"
// @Param order query string false "Sort order (asc or desc)" default(asc)
// @Success 200 {object} common.PaginatedResponse[ClassDTO] "Paginated list of classes"
// @Failure 400 {object} common.ErrorResponse "Invalid filter or sort field"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /classes [get]
func (h *Handler) GetAll(c *gin.Context) {
	params := common.GetPaginationParams(c)
	entities, pagination, err := h.service.GetAllClasses(c.Request.Context(), params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve "+h.GetEntityName()+"s", err.Error())
		return
	}

	dtos := make([]ClassDTO, len(entities))
	for i, entity := range entities {
		dtos[i] = *h.mapper.FromDomain(&entity)
	}

	response := common.PaginatedResponse[ClassDTO]{
		Data:       dtos,
		Pagination: pagination,
	}
	c.JSON(http.StatusOK, response)
}

// GetByID retrieves a class by ID
// @Summary Get class by ID
// @Description Retrieve a specific class by its unique identifier
// @Tags classes
// @Accept json
// @Produce json
// @Param id path int true "Class ID" minimum(1)
// @Success 200 {object} ClassDTO "Class details"
// @Failure 400 {object} common.ErrorResponse "Invalid class ID"
// @Failure 404 {object} common.ErrorResponse "Class not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /classes/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := common.ParseIDFromPath(c, h.GetEntityName())
	if err != nil {
		return
	}

	entity, err := h.service.GetClass(c.Request.Context(), id)
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

// Create creates a new class
// @Summary Create a new class
// @Description Create a new class with the provided name and capacity
// @Tags classes
// @Accept json
// @Produce json
// @Param class body CreateClassDTO true "Class creation data"
// @Success 201 {object} ClassDTO "Created class"
// @Failure 400 {object} common.ErrorResponse "Invalid input data"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /classes [post]
func (h *Handler) Create(c *gin.Context) {
	createDTO, err := h.BindCreateJSON(c)
	if err != nil {
		return
	}

	entity := h.mapper.ToDomain(createDTO)
	if err := h.service.CreateClass(c.Request.Context(), entity); err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to create "+h.GetEntityName(), err.Error())
		return
	}

	createdDTO := h.mapper.FromDomain(entity)
	c.JSON(http.StatusCreated, createdDTO)
}

// Update updates an existing class
// @Summary Update an existing class
// @Description Update an existing class's name and capacity by ID
// @Tags classes
// @Accept json
// @Produce json
// @Param id path int true "Class ID" minimum(1)
// @Param class body UpdateClassDTO true "Class update data"
// @Success 200 {object} ClassDTO "Updated class"
// @Failure 400 {object} common.ErrorResponse "Invalid input data"
// @Failure 404 {object} common.ErrorResponse "Class not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /classes/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, updateDTO, err := h.ParseIDAndBindJSON(c)
	if err != nil {
		return
	}

	entity := h.mapper.ToDomainWithID(updateDTO, id)
	if err := h.service.UpdateClass(c.Request.Context(), entity); err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to update "+h.GetEntityName(), err.Error())
		return
	}

	updatedDTO := h.mapper.FromDomain(entity)
	c.JSON(http.StatusOK, updatedDTO)
}

// Delete removes a class
// @Summary Delete a class
// @Description Delete a class by its ID
// @Tags classes
// @Accept json
// @Produce json
// @Param id path int true "Class ID" minimum(1)
// @Success 200 {object} common.SuccessResponse "Class deleted successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid class ID"
// @Failure 404 {object} common.ErrorResponse "Class not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /classes/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := common.ParseIDFromPath(c, h.GetEntityName())
	if err != nil {
		return
	}

	if err := h.service.DeleteClass(c.Request.Context(), id); err != nil {
		common.HandleError(c, err, "Failed to delete "+h.GetEntityName())
		return
	}

	common.RespondWithSuccess(c, http.StatusOK, h.GetEntityName()+" deleted successfully")
}
