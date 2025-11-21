package car_image

import (
	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/middleware"
)

func RegisterRoutes(router gin.IRouter, handler *CarImageHandler, tokenSvc ports.TokenService) {
	api := router.Group("/v1/cars/images")
	api.Use(middleware.AuthMiddleware(tokenSvc))
	{
		api.POST("", handler.CreateCarImage)
		api.DELETE("/:id", handler.DeleteCarImage)
		api.GET("", handler.GetCarImagesList)
	}
}
