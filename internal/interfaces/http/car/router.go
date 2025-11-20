package car

import (
	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/middleware"
)

func RegisterRoutes(router gin.IRouter, handler *CarHandler, tokenSvc ports.TokenService) {
	api := router.Group("/v1/cars")
	api.Use(middleware.AuthMiddleware(tokenSvc))

	{
		api.POST("", handler.CreateCar)
	}
}
