package ports

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (*entities.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserById(ctx context.Context, id int64) (*entities.User, error)
	ConfirmEmailVerification(ctx context.Context, email string) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error)
	DeleteUser(ctx context.Context, id int64) error
}
