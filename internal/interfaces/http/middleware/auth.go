package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

func AuthMiddleware(tokenSvc ports.TokenService) gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			_ = c.Error(apperrors.ErrUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "authorization required"})
			return
		}

		var token string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			_ = c.Error(apperrors.ErrUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization header"})
			return
		}

		userID, err := tokenSvc.ValidateAccessToken(token)
		if err != nil {
			_ = c.Error(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
			return
		}

		// set user id into context for handlers
		c.Set("user_id", userID)
		c.Next()
	}
}
