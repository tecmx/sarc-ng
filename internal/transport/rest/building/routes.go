package building

import (
	"sarc-ng/internal/domain/building"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the building routes
func RegisterRoutes(rg *gin.RouterGroup, service building.Usecase) {
	handler := NewHandler(service)

	buildings := rg.Group("/buildings")
	{
		buildings.GET("", handler.GetAll)
		buildings.POST("", handler.Create)
		buildings.GET("/:id", handler.GetByID)
		buildings.PUT("/:id", handler.Update)
		buildings.DELETE("/:id", handler.Delete)
	}
}
