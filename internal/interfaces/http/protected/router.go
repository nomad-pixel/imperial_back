package protected

import (
	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/middleware"
)

func RegisterRoutes(router gin.IRouter, handler *ProtectedHandler, tokenSvc ports.TokenService) {
	api := router.Group("/v1/protected")
	api.Use(middleware.AuthMiddleware(tokenSvc))
	{
		api.GET("/me", handler.Me)
	}
}
