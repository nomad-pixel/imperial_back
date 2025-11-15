package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

// SignUpUsecase интерфейс для регистрации пользователей
type SignUpUsecase interface {
	Execute(ctx context.Context, email, password string) (*entities.User, error)
}
