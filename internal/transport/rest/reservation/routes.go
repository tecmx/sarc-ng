package reservation

import (
	"sarc-ng/internal/domain/reservation"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the reservation routes
func RegisterRoutes(rg *gin.RouterGroup, service reservation.Usecase) {
	handler := NewHandler(service)

	reservations := rg.Group("/reservations")
	{
		reservations.GET("", handler.GetAll)
		reservations.POST("", handler.Create)
		reservations.GET("/:id", handler.GetByID)
		reservations.PUT("/:id", handler.Update)
		reservations.DELETE("/:id", handler.Delete)
	}
}
