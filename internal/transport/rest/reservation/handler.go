package reservation

import (
	"net/http"
	"sarc-ng/internal/domain/reservation"
	"sarc-ng/internal/transport/common"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for reservation operations
type Handler struct {
	*common.BaseHandler[reservation.Reservation, CreateReservationDTO, UpdateReservationDTO, ReservationDTO]
	service reservation.Usecase
	mapper  *Mapper
}

// NewHandler creates a new reservation handler
func NewHandler(service reservation.Usecase) *Handler {
	mapper := NewMapper()
	baseHandler := common.NewBaseHandler[reservation.Reservation, CreateReservationDTO, UpdateReservationDTO, ReservationDTO](
		"reservation")
	return &Handler{
		BaseHandler: baseHandler,
		service:     service,
		mapper:      mapper,
	}
}

// GetAll retrieves all reservations with pagination
// @Summary Get all reservations (paginated)
// @Description Retrieve a paginated list of all reservations with support for search, filtering, and sorting
// @Tags reservations
// @Accept json
// @Produce json
// @Security CognitoOAuth
// @Security BearerAuth
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param page_size query int false "Page size (default: 20, max: 100)" minimum(1) maximum(100)
// @Param search query string false "Search across purpose and description fields"
// @Param filter[resource_id] query int false "Filter by resource ID"
// @Param filter[user_id] query int false "Filter by user ID"
// @Param filter[status] query string false "Filter by exact status (pending, confirmed, cancelled, etc.)"
// @Param filter[purpose] query string false "Filter by exact purpose match"
// @Param sort query string false "Sort field (id, resource_id, user_id, start_time, end_time, status, created_at, updated_at)"
// @Param order query string false "Sort order (asc or desc)" default(asc)
// @Success 200 {object} common.PaginatedResponse[ReservationDTO] "Paginated list of reservations"
// @Failure 400 {object} common.ErrorResponse "Invalid filter or sort field"
// @Failure 401 {object} common.ErrorResponse "Unauthorized"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /reservations [get]
func (h *Handler) GetAll(c *gin.Context) {
	params := common.GetPaginationParams(c)
	entities, pagination, err := h.service.GetAllReservations(c.Request.Context(), params)
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve "+h.GetEntityName()+"s", err.Error())
		return
	}

	dtos := make([]ReservationDTO, len(entities))
	for i, entity := range entities {
		dtos[i] = *h.mapper.FromDomain(&entity)
	}

	response := common.PaginatedResponse[ReservationDTO]{
		Data:       dtos,
		Pagination: pagination,
	}
	c.JSON(http.StatusOK, response)
}

// GetByID retrieves a reservation by ID
// @Summary Get reservation by ID
// @Description Retrieve a specific reservation by its unique identifier
// @Tags reservations
// @Accept json
// @Produce json
// @Security CognitoOAuth
// @Security BearerAuth
// @Param id path int true "Reservation ID" minimum(1)
// @Success 200 {object} ReservationDTO "Reservation details"
// @Failure 400 {object} common.ErrorResponse "Invalid reservation ID"
// @Failure 401 {object} common.ErrorResponse "Unauthorized"
// @Failure 404 {object} common.ErrorResponse "Reservation not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /reservations/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := common.ParseIDFromPath(c, h.GetEntityName())
	if err != nil {
		return
	}

	entity, err := h.service.GetReservation(c.Request.Context(), id)
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

// Create creates a new reservation
// @Summary Create a new reservation
// @Description Create a new reservation with resource, user, and time information
// @Tags reservations
// @Accept json
// @Produce json
// @Security CognitoOAuth
// @Security BearerAuth
// @Param reservation body CreateReservationDTO true "Reservation creation data"
// @Success 201 {object} ReservationDTO "Created reservation"
// @Failure 400 {object} common.ErrorResponse "Invalid input data"
// @Failure 401 {object} common.ErrorResponse "Unauthorized"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /reservations [post]
func (h *Handler) Create(c *gin.Context) {
	createDTO, err := h.BindCreateJSON(c)
	if err != nil {
		return
	}

	entity := h.mapper.ToDomain(createDTO)
	if err := h.service.CreateReservation(c.Request.Context(), entity); err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to create "+h.GetEntityName(), err.Error())
		return
	}

	createdDTO := h.mapper.FromDomain(entity)
	c.JSON(http.StatusCreated, createdDTO)
}

// Update updates an existing reservation
// @Summary Update an existing reservation
// @Description Update an existing reservation's resource, user, time, and status information by ID
// @Tags reservations
// @Accept json
// @Produce json
// @Security CognitoOAuth
// @Security BearerAuth
// @Param id path int true "Reservation ID" minimum(1)
// @Param reservation body UpdateReservationDTO true "Reservation update data"
// @Success 200 {object} ReservationDTO "Updated reservation"
// @Failure 400 {object} common.ErrorResponse "Invalid input data"
// @Failure 401 {object} common.ErrorResponse "Unauthorized"
// @Failure 404 {object} common.ErrorResponse "Reservation not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /reservations/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, updateDTO, err := h.ParseIDAndBindJSON(c)
	if err != nil {
		return
	}

	entity := h.mapper.ToDomainWithID(updateDTO, id)
	if err := h.service.UpdateReservation(c.Request.Context(), entity); err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to update "+h.GetEntityName(), err.Error())
		return
	}

	updatedDTO := h.mapper.FromDomain(entity)
	c.JSON(http.StatusOK, updatedDTO)
}

// Delete removes a reservation
// @Summary Delete a reservation
// @Description Delete a reservation by its ID
// @Tags reservations
// @Accept json
// @Produce json
// @Security CognitoOAuth
// @Security BearerAuth
// @Param id path int true "Reservation ID" minimum(1)
// @Success 200 {object} common.SuccessResponse "Reservation deleted successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid reservation ID"
// @Failure 401 {object} common.ErrorResponse "Unauthorized"
// @Failure 404 {object} common.ErrorResponse "Reservation not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /reservations/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := common.ParseIDFromPath(c, h.GetEntityName())
	if err != nil {
		return
	}

	if err := h.service.DeleteReservation(c.Request.Context(), id); err != nil {
		common.HandleError(c, err, "Failed to delete "+h.GetEntityName())
		return
	}

	common.RespondWithSuccess(c, http.StatusOK, h.GetEntityName()+" deleted successfully")
}
