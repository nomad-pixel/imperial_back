package driver

import (
	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, handler *DriverHandler, tokenSvc ports.TokenService) {
	drivers := router.Group("/v1/drivers")
	drivers.Use(middleware.AuthMiddleware(tokenSvc))
	{
		drivers.POST("", handler.CreateDriver)
		drivers.GET("/:id", handler.GetDriverByID)
		drivers.GET("", handler.ListDrivers)
		drivers.PUT("/:id", handler.UpdateDriver)
		drivers.DELETE("/:id", handler.DeleteDriver)
		drivers.PUT("/:id/photo", handler.UploadDriverPhoto)
	}
}
