package car

import (
	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/middleware"
)

func RegisterRoutes(router gin.IRouter, handler *CarImageHandler, tokenSvc ports.TokenService) {
	api := router.Group("/v1/cars/:id/images")
	api.Use(middleware.AuthMiddleware(tokenSvc))
	{
		api.POST("", handler.CreateCarImage)
		api.GET("", handler.GetCarImagesList)
		api.DELETE("/:image_id", handler.DeleteCarImage)
	}
}
