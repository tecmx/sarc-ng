package lesson

import (
	"sarc-ng/internal/domain/lesson"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the lesson routes
func RegisterRoutes(rg *gin.RouterGroup, service lesson.Usecase) {
	handler := NewHandler(service)

	lessons := rg.Group("/lessons")
	{
		lessons.GET("", handler.GetAll)
		lessons.POST("", handler.Create)
		lessons.GET("/:id", handler.GetByID)
		lessons.PUT("/:id", handler.Update)
		lessons.DELETE("/:id", handler.Delete)
	}
}
