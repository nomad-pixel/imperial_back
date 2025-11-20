package auth

import (
	"time"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"password123"`
}

type SignUpResponse struct {
	ID         int64     `json:"id" example:"123"`
	Email      string    `json:"email" example:"user@example.com"`
	IsVerified bool      `json:"is_verified" example:"false"`
	CreatedAt  time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt  time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

func ToSignUpResponse(user *entities.User) SignUpResponse {
	return SignUpResponse{
		ID:         user.ID,
		Email:      user.Email,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

type ConfirmEmailRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
	Code  string `json:"code" binding:"required" example:"123456"`
}

type VerifyEmailResponse struct {
	Message string `json:"message" example:"Email успешно отправлен"`
}

type ConfirmEmailResponse struct {
	Message string `json:"message" example:"Email успешно подтвержден"`
}

type ErrorResponse struct {
	Code    string                 `json:"code" example:"VALIDATION_ERROR"`
	Message string                 `json:"message" example:"Неверный формат данных"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type SignInRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type SignInResponse struct {
	User   SignUpResponse  `json:"user"`
	Tokens entities.Tokens `json:"tokens"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"<refresh_token>"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token" example:"<access_token>"`
	RefreshToken string `json:"refresh_token,omitempty" example:"<refresh_token>"`
}
