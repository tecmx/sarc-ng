package building

import (
	"net/http"
	"sarc-ng/internal/domain/building"
	"sarc-ng/internal/transport/common"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for building operations
type Handler struct {
	*common.BaseHandler[building.Building, CreateBuildingDTO, UpdateBuildingDTO, BuildingDTO]
	service building.Usecase
	mapper  *Mapper
}

// NewHandler creates a new building handler
func NewHandler(service building.Usecase) *Handler {
	mapper := NewMapper()
	baseHandler := common.NewBaseHandler[building.Building, CreateBuildingDTO, UpdateBuildingDTO, BuildingDTO](
		"building")
	return &Handler{
		BaseHandler: baseHandler,
		service:     service,
		mapper:      mapper,
	}
}

// GetAll retrieves all buildings
// @Summary Get all buildings
// @Description Retrieve a list of all buildings in the system
// @Tags buildings
// @Accept json
// @Produce json
// @Success 200 {array} BuildingDTO "List of buildings"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /buildings [get]
func (h *Handler) GetAll(c *gin.Context) {
	entities, err := h.service.GetAllBuildings()
	if err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to retrieve "+h.GetEntityName()+"s", err.Error())
		return
	}

	dtos := make([]BuildingDTO, len(entities))
	for i, entity := range entities {
		dtos[i] = *h.mapper.FromDomain(&entity)
	}
	c.JSON(http.StatusOK, dtos)
}

// GetByID retrieves a building by ID
// @Summary Get building by ID
// @Description Retrieve a specific building by its unique identifier
// @Tags buildings
// @Accept json
// @Produce json
// @Param id path int true "Building ID" minimum(1)
// @Success 200 {object} BuildingDTO "Building details"
// @Failure 400 {object} common.ErrorResponse "Invalid building ID"
// @Failure 404 {object} common.ErrorResponse "Building not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /buildings/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := common.ParseIDFromPath(c, h.GetEntityName())
	if err != nil {
		return
	}

	entity, err := h.service.GetBuilding(id)
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

// Create creates a new building
// @Summary Create a new building
// @Description Create a new building with the provided name and code
// @Tags buildings
// @Accept json
// @Produce json
// @Param building body CreateBuildingDTO true "Building creation data"
// @Success 201 {object} BuildingDTO "Created building"
// @Failure 400 {object} common.ErrorResponse "Invalid input data"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /buildings [post]
func (h *Handler) Create(c *gin.Context) {
	createDTO, err := h.BindCreateJSON(c)
	if err != nil {
		return
	}

	entity := h.mapper.ToDomain(createDTO)
	if err := h.service.CreateBuilding(entity); err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to create "+h.GetEntityName(), err.Error())
		return
	}

	createdDTO := h.mapper.FromDomain(entity)
	c.JSON(http.StatusCreated, createdDTO)
}

// Update updates an existing building
// @Summary Update an existing building
// @Description Update an existing building's name and code by ID
// @Tags buildings
// @Accept json
// @Produce json
// @Param id path int true "Building ID" minimum(1)
// @Param building body UpdateBuildingDTO true "Building update data"
// @Success 200 {object} BuildingDTO "Updated building"
// @Failure 400 {object} common.ErrorResponse "Invalid input data"
// @Failure 404 {object} common.ErrorResponse "Building not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /buildings/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, updateDTO, err := h.ParseIDAndBindJSON(c)
	if err != nil {
		return
	}

	entity := h.mapper.ToDomainWithID(updateDTO, id)
	if err := h.service.UpdateBuilding(entity); err != nil {
		common.RespondWithError(c, http.StatusInternalServerError, "Failed to update "+h.GetEntityName(), err.Error())
		return
	}

	updatedDTO := h.mapper.FromDomain(entity)
	c.JSON(http.StatusOK, updatedDTO)
}

// Delete removes a building
// @Summary Delete a building
// @Description Delete a building by its ID
// @Tags buildings
// @Accept json
// @Produce json
// @Param id path int true "Building ID" minimum(1)
// @Success 200 {object} common.SuccessResponse "Building deleted successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid building ID"
// @Failure 404 {object} common.ErrorResponse "Building not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /buildings/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := common.ParseIDFromPath(c, h.GetEntityName())
	if err != nil {
		return
	}

	if err := h.service.DeleteBuilding(id); err != nil {
		common.HandleError(c, err, "Failed to delete "+h.GetEntityName())
		return
	}

	common.RespondWithSuccess(c, http.StatusOK, h.GetEntityName()+" deleted successfully")
}
