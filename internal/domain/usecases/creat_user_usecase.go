package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type CreateUserUsecase struct {
	userRepository ports.UserRepository
}

func NewCreateUserUsecase(userRepository ports.UserRepository) *CreateUserUsecase {
	return &CreateUserUsecase{userRepository: userRepository}
}

func (uc *CreateUserUsecase) Execute(ctx context.Context, email, password string) (*entities.User, error) {
	user, err := uc.userRepository.CreateUser(ctx, email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
