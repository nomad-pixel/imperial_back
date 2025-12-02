package celebrity

import (
	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/middleware"
)

func RegisterRoutes(router gin.IRouter, handler *CelebrityHandler, tokenSvc ports.TokenService) {
	api := router.Group("/v1/celebrities")

	api.Use(middleware.AuthMiddleware(tokenSvc))

	{
		api.POST("", handler.CreateCelebrity)
		api.GET("", handler.ListCelebrities)
		api.GET("/:id", handler.GetCelebrityByID)
		api.PUT("/:id", handler.UpdateCelebrity)
		api.PUT("/:id/image", handler.UploadCelebrityImage)
		api.DELETE("/:id", handler.DeleteCelebrity)
	}
}
