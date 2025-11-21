package car_mark

import (
	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/middleware"
)

func RegisterRoutes(router gin.IRouter, handler *CarMarkHandler, tokenSvc ports.TokenService) {
	api := router.Group("/v1/cars/car-marks")

	// Public GET endpoints
	api.GET("", handler.GetCarMarks)
	api.GET("/:id", handler.GetCarMark)

	// Protected endpoints
	api.Use(middleware.AuthMiddleware(tokenSvc))
	{
		api.POST("", handler.CreateCarMark)
		api.PUT("/:id", handler.UpdateCarMark)
		api.DELETE("/:id", handler.DeleteCarMark)
	}
}
