package car_tag

import (
	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/middleware"
)

func RegisterRoutes(router gin.IRouter, handler *CarTagHandler, tokenSvc ports.TokenService) {
	api := router.Group("/v1/cars/car-tags")

	api.GET("", handler.GetCarTags)
	api.GET("/:id", handler.GetCarTag)

	api.Use(middleware.AuthMiddleware(tokenSvc))
	{
		api.POST("", handler.CreateCarTag)
		api.PUT("/:id", handler.UpdateCarTag)
		api.DELETE("/:id", handler.DeleteCarTag)
	}
}
