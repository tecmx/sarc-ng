package resource

import (
	"sarc-ng/internal/domain/resource"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up the resource routes
func RegisterRoutes(rg *gin.RouterGroup, service resource.Usecase) {
	handler := NewHandler(service)

	resources := rg.Group("/resources")
	{
		resources.GET("", handler.GetAll)
		resources.POST("", handler.Create)
		resources.GET("/:id", handler.GetByID)
		resources.PUT("/:id", handler.Update)
		resources.DELETE("/:id", handler.Delete)
	}
}
