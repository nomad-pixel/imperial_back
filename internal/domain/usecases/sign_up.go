package usecases

import (
	"context"

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
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return u.userRepo.CreateUser(ctx, email, string(passwordHash))
}
