package usecases

import (
	"context"

	apperrors "github.com/nomad-pixel/imperial/pkg/errors"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"golang.org/x/crypto/bcrypt"
)

type SignUpUsecase interface {
	Execute(ctx context.Context, email, password string) (*entities.User, error)
}

type signUpUsecase struct {
	userRepo ports.UserRepository
}

func NewSignUpUsecase(userRepo ports.UserRepository) SignUpUsecase {
	return &signUpUsecase{
		userRepo: userRepo,
	}
}

func (u *signUpUsecase) Execute(ctx context.Context, email, password string) (*entities.User, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if user != nil && err == nil {
		return nil, apperrors.ErrUserNotFound
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, "Ошибка при регистрации")
	}
	return u.userRepo.CreateUser(ctx, email, string(passwordHash))
}
