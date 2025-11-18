package auth

import (
	"time"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

// SignUpRequest представляет запрос на регистрацию
// @Description Данные для регистрации нового пользователя
type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
}

// SignUpResponse представляет ответ после регистрации
// @Description Информация о зарегистрированном пользователе
type SignUpResponse struct {
	ID         int64     `json:"id" example:"123"`
	Email      string    `json:"email" example:"user@example.com"`
	VerifiedAt bool      `json:"verified_at" example:"false"`
	CreatedAt  time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt  time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

func ToSignUpResponse(user *entities.User) SignUpResponse {
	return SignUpResponse{
		ID:         user.ID,
		Email:      user.Email,
		VerifiedAt: user.VerifiedAt,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

// VerifyEmailRequest представляет запрос на отправку кода верификации
// @Description Email адрес для отправки кода верификации
type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// ConfirmEmailRequest представляет запрос на подтверждение email
// @Description Email адрес и код для подтверждения email
type ConfirmEmailRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	Code  string `json:"code" binding:"required" example:"123456"`
}

// VerifyEmailResponse представляет ответ после отправки кода верификации
// @Description Сообщение об успешной отправке email
type VerifyEmailResponse struct {
	Message string `json:"message" example:"Email успешно отправлен"`
}

// ConfirmEmailResponse представляет ответ после подтверждения email
// @Description Сообщение об успешной подтверждений email
type ConfirmEmailResponse struct {
	Message string `json:"message" example:"Email успешно подтвержден"`
}

// ErrorResponse представляет ответ с ошибкой
// @Description Информация об ошибке
type ErrorResponse struct {
	Code    string                 `json:"code" example:"VALIDATION_ERROR"`
	Message string                 `json:"message" example:"Неверный формат данных"`
	Details map[string]interface{} `json:"details,omitempty"`
}
