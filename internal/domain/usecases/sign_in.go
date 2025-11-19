package usecases

import (
	"context"

	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type signInUsecase struct {
	userRepo ports.UserRepository
	tokenSvc ports.TokenService
}

type SignInUsecase interface {
	Execute(ctx context.Context, email, password string) (*entities.User, *entities.Tokens, error)
}

func NewSignInUsecase(userRepo ports.UserRepository, tokenSvc ports.TokenService) SignInUsecase {
	return &signInUsecase{
		userRepo: userRepo,
		tokenSvc: tokenSvc,
	}
}

func (u *signInUsecase) Execute(ctx context.Context, email, password string) (*entities.User, *entities.Tokens, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, nil, apperrors.ErrUserNotFound
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, nil, apperrors.ErrInvalidCredentials
	}
	tokens, err := u.tokenSvc.GenerateTokens(user)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}
