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
	existingUser, err := u.userRepo.GetUserByEmail(ctx, email)
	if existingUser != nil && err == nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, "User with this email already exists")
	}

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "Ошибка при регистрации")
	}

	// Use entity factory with validation
	user, err := entities.NewUser(email, string(passwordHash))
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, err.Error())
	}

	// Create user in repository
	createdUser, err := u.userRepo.CreateUser(ctx, user.Email, user.PasswordHash)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "Failed to create user")
	}

	return createdUser, nil
}
