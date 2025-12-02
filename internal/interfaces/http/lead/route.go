package lead

import (
	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/middleware"
)

func RegisterRoutes(router *gin.RouterGroup, handler *LeadHandler, tokenSvc ports.TokenService) {
	leads := router.Group("/v1/leads")
	leads.Use(middleware.AuthMiddleware(tokenSvc))
	{
		leads.POST("", handler.CreateLead)
		leads.GET("/:id", handler.GetLeadByID)
		leads.GET("", handler.ListLeads)
		leads.DELETE("/:id", handler.DeleteLead)
	}
}
