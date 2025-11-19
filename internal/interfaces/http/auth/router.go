package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(router gin.IRouter, handler *AuthHandler) {
	api := router.Group("/v1/auth")
	{
		api.POST("/sign-up", handler.SignUp)
		api.POST("/verify-email", handler.VerifyEmail)
		api.POST("/confirm-email", handler.ConfirmEmail)
		api.POST("/sign-in", handler.SignIn)
	}
}
