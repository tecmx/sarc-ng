package class

import (
	"sarc-ng/internal/domain/class"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the class routes
func RegisterRoutes(rg *gin.RouterGroup, service class.Usecase) {
	handler := NewHandler(service)

	classes := rg.Group("/classes")
	{
		classes.GET("", handler.GetAll)
		classes.POST("", handler.Create)
		classes.GET("/:id", handler.GetByID)
		classes.PUT("/:id", handler.Update)
		classes.DELETE("/:id", handler.Delete)
	}
}
