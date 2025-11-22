package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type ErrorResponse struct {
	Code    errors.ErrorCode       `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if appErr, ok := errors.AsAppError(err); ok {
				if appErr.StatusCode >= 500 {
					log.Printf("ERROR: %v", appErr)
				}

				c.JSON(appErr.StatusCode, ErrorResponse{
					Code:    appErr.Code,
					Message: appErr.Message,
					Details: appErr.Details,
				})
				return
			}

			log.Printf("Unexpected error: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    errors.ErrCodeInternal,
				Message: "Внутренняя ошибка сервера",
			})
		}
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v", err)
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:    errors.ErrCodeInternal,
					Message: "Внутренняя ошибка сервера",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
