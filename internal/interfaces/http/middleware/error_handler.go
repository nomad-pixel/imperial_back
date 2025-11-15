package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

// ErrorResponse представляет структуру ответа с ошибкой
type ErrorResponse struct {
	Code    errors.ErrorCode       `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ErrorHandler middleware для централизованной обработки ошибок
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Проверяем наличие ошибок после выполнения handlers
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Пытаемся привести к AppError
			if appErr, ok := errors.AsAppError(err); ok {
				// Логируем ошибку
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

			// Если это не AppError, возвращаем общую ошибку
			log.Printf("Unexpected error: %v", err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Code:    errors.ErrCodeInternal,
				Message: "Внутренняя ошибка сервера",
			})
		}
	}
}

// Recovery middleware для обработки паник
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
