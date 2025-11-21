package car_category

import (
	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/middleware"
)

func RegisterRoutes(router gin.IRouter, handler *CarCategoryHandler, tokenSvc ports.TokenService) {
	api := router.Group("/v1/cars/car-categories")

	// Public GET endpoints
	api.GET("", handler.GetCarCategories)
	api.GET("/:id", handler.GetCarCategory)

	// Protected endpoints
	api.Use(middleware.AuthMiddleware(tokenSvc))
	{
		api.POST("", handler.CreateCarCategory)
		api.PUT("/:id", handler.UpdateCarCategory)
		api.DELETE("/:id", handler.DeleteCarCategory)
	}
}
