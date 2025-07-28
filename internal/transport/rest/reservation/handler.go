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

// GetAll retrieves all reservations
// @Summary Get all reservations
// @Description Retrieve a list of all reservations in the system
// @Tags reservations
// @Accept json
// @Produce json
// @Success 200 {array} ReservationDTO "List of reservations"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /reservations [get]
func (h *Handler) GetAll(c *gin.Context) {
	entities, err := h.service.GetAllReservations()
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve "+h.GetEntityName()+"s", err.Error())
		return
	}

	dtos := make([]ReservationDTO, len(entities))
	for i, entity := range entities {
		dtos[i] = *h.mapper.FromDomain(&entity)
	}
	c.JSON(http.StatusOK, dtos)
}

// GetByID retrieves a reservation by ID
// @Summary Get reservation by ID
// @Description Retrieve a specific reservation by its unique identifier
// @Tags reservations
// @Accept json
// @Produce json
// @Param id path int true "Reservation ID" minimum(1)
// @Success 200 {object} ReservationDTO "Reservation details"
// @Failure 400 {object} common.ErrorResponse "Invalid reservation ID"
// @Failure 404 {object} common.ErrorResponse "Reservation not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /reservations/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := common.ParseIDFromPath(c, h.GetEntityName())
	if err != nil {
		return // Error already handled by ParseIDFromPath
	}

	entity, err := h.service.GetReservation(id)
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
// @Param reservation body CreateReservationDTO true "Reservation creation data"
// @Success 201 {object} ReservationDTO "Created reservation"
// @Failure 400 {object} common.ErrorResponse "Invalid input data"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /reservations [post]
func (h *Handler) Create(c *gin.Context) {
	createDTO, err := h.BindCreateJSON(c)
	if err != nil {
		return // Error already handled by BindCreateJSON
	}

	entity := h.mapper.ToDomain(createDTO)
	if err := h.service.CreateReservation(entity); err != nil {
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
// @Param id path int true "Reservation ID" minimum(1)
// @Param reservation body UpdateReservationDTO true "Reservation update data"
// @Success 200 {object} ReservationDTO "Updated reservation"
// @Failure 400 {object} common.ErrorResponse "Invalid input data"
// @Failure 404 {object} common.ErrorResponse "Reservation not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /reservations/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, updateDTO, err := h.ParseIDAndBindJSON(c)
	if err != nil {
		return // Error already handled by ParseIDAndBindJSON
	}

	entity := h.mapper.ToDomainWithID(updateDTO, id)
	if err := h.service.UpdateReservation(entity); err != nil {
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
// @Param id path int true "Reservation ID" minimum(1)
// @Success 200 {object} common.SuccessResponse "Reservation deleted successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid reservation ID"
// @Failure 404 {object} common.ErrorResponse "Reservation not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /reservations/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := common.ParseIDFromPath(c, h.GetEntityName())
	if err != nil {
		return // Error already handled by ParseIDFromPath
	}

	if err := h.service.DeleteReservation(id); err != nil {
		common.HandleError(c, err, "Failed to delete "+h.GetEntityName())
		return
	}

	common.RespondWithSuccess(c, http.StatusOK, h.GetEntityName()+" deleted successfully")
}
