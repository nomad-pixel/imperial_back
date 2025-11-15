package auth

import (
	"time"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

// SignUpRequest представляет запрос на регистрацию
type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
}

// SignUpResponse представляет ответ на регистрацию
type SignUpResponse struct {
	ID         int64     `json:"id" example:"123"`
	Email      string    `json:"email" example:"user@example.com"`
	VerifiedAt bool      `json:"verified_at" example:"false"`
	CreatedAt  time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt  time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// ToSignUpResponse преобразует Entity в DTO
func ToSignUpResponse(user *entities.User) SignUpResponse {
	return SignUpResponse{
		ID:         user.ID,
		Email:      user.Email,
		VerifiedAt: user.VerifiedAt,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Code    string                 `json:"code" example:"VALIDATION_ERROR"`
	Message string                 `json:"message" example:"Неверный формат данных"`
	Details map[string]interface{} `json:"details,omitempty"`
}
