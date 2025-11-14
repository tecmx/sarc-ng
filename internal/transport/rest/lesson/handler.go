package lesson

import (
	"net/http"
	"sarc-ng/internal/domain/lesson"
	"sarc-ng/internal/transport/common"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for lesson operations
type Handler struct {
	*common.BaseHandler[lesson.Lesson, CreateLessonDTO, UpdateLessonDTO, LessonDTO]
	service lesson.Usecase
	mapper  *Mapper
}

// NewHandler creates a new lesson handler
func NewHandler(service lesson.Usecase) *Handler {
	mapper := NewMapper()
	baseHandler := common.NewBaseHandler[lesson.Lesson, CreateLessonDTO, UpdateLessonDTO, LessonDTO](
		"lesson")
	return &Handler{
		BaseHandler: baseHandler,
		service:     service,
		mapper:      mapper,
	}
}

// GetAll retrieves all lessons with pagination
// @Summary Get all lessons (paginated)
// @Description Retrieve a paginated list of all lessons with support for search, filtering, and sorting
// @Tags lessons
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param page_size query int false "Page size (default: 20, max: 100)" minimum(1) maximum(100)
// @Param search query string false "Search across title and description fields"
// @Param filter[title] query string false "Filter by exact title match"
// @Param filter[duration] query int false "Filter by exact duration (in minutes)"
// @Param sort query string false "Sort field (id, title, duration, start_time, end_time, created_at, updated_at)"
// @Param order query string false "Sort order (asc or desc)" default(asc)
// @Success 200 {object} common.PaginatedResponse[LessonDTO] "Paginated list of lessons"
// @Failure 400 {object} common.ErrorResponse "Invalid filter or sort field"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /lessons [get]
func (h *Handler) GetAll(c *gin.Context) {
	params := common.GetPaginationParams(c)
	entities, pagination, err := h.service.GetAllLessons(c.Request.Context(), params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve "+h.GetEntityName()+"s", err.Error())
		return
	}

	dtos := make([]LessonDTO, len(entities))
	for i, entity := range entities {
		dtos[i] = *h.mapper.FromDomain(&entity)
	}

	response := common.PaginatedResponse[LessonDTO]{
		Data:       dtos,
		Pagination: pagination,
	}
	c.JSON(http.StatusOK, response)
}

// GetByID retrieves a lesson by ID
// @Summary Get lesson by ID
// @Description Retrieve a specific lesson by its unique identifier
// @Tags lessons
// @Accept json
// @Produce json
// @Param id path int true "Lesson ID" minimum(1)
// @Success 200 {object} LessonDTO "Lesson details"
// @Failure 400 {object} common.ErrorResponse "Invalid lesson ID"
// @Failure 404 {object} common.ErrorResponse "Lesson not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /lessons/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := common.ParseIDFromPath(c, h.GetEntityName())
	if err != nil {
		return
	}

	entity, err := h.service.GetLesson(c.Request.Context(), id)
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

// Create creates a new lesson
// @Summary Create a new lesson
// @Description Create a new lesson with the provided title, duration, and start time
// @Tags lessons
// @Accept json
// @Produce json
// @Param lesson body CreateLessonDTO true "Lesson creation data"
// @Success 201 {object} LessonDTO "Created lesson"
// @Failure 400 {object} common.ErrorResponse "Invalid input data"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /lessons [post]
func (h *Handler) Create(c *gin.Context) {
	createDTO, err := h.BindCreateJSON(c)
	if err != nil {
		return
	}

	entity := h.mapper.ToDomain(createDTO)
	if err := h.service.CreateLesson(c.Request.Context(), entity); err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to create "+h.GetEntityName(), err.Error())
		return
	}

	createdDTO := h.mapper.FromDomain(entity)
	c.JSON(http.StatusCreated, createdDTO)
}

// Update updates an existing lesson
// @Summary Update an existing lesson
// @Description Update an existing lesson's title, duration, and start time by ID
// @Tags lessons
// @Accept json
// @Produce json
// @Param id path int true "Lesson ID" minimum(1)
// @Param lesson body UpdateLessonDTO true "Lesson update data"
// @Success 200 {object} LessonDTO "Updated lesson"
// @Failure 400 {object} common.ErrorResponse "Invalid input data"
// @Failure 404 {object} common.ErrorResponse "Lesson not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /lessons/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, updateDTO, err := h.ParseIDAndBindJSON(c)
	if err != nil {
		return
	}

	entity := h.mapper.ToDomainWithID(updateDTO, id)
	if err := h.service.UpdateLesson(c.Request.Context(), entity); err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to update "+h.GetEntityName(), err.Error())
		return
	}

	updatedDTO := h.mapper.FromDomain(entity)
	c.JSON(http.StatusOK, updatedDTO)
}

// Delete removes a lesson
// @Summary Delete a lesson
// @Description Delete a lesson by its ID
// @Tags lessons
// @Accept json
// @Produce json
// @Param id path int true "Lesson ID" minimum(1)
// @Success 200 {object} common.SuccessResponse "Lesson deleted successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid lesson ID"
// @Failure 404 {object} common.ErrorResponse "Lesson not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /lessons/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := common.ParseIDFromPath(c, h.GetEntityName())
	if err != nil {
		return
	}

	if err := h.service.DeleteLesson(c.Request.Context(), id); err != nil {
		common.HandleError(c, err, "Failed to delete "+h.GetEntityName())
		return
	}

	common.RespondWithSuccess(c, http.StatusOK, h.GetEntityName()+" deleted successfully")
}
