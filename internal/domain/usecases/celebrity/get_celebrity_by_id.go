package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type getCelebrityByIdUsecase struct {
	celebrityRepo ports.CelebrityRepository
}

type GetCelebrityByIdUsecase interface {
	Execute(ctx context.Context, id int64) (*entities.Celebrity, error)
}

func NewGetCelebrityByIdUsecase(celebrityRepo ports.CelebrityRepository) GetCelebrityByIdUsecase {
	return &getCelebrityByIdUsecase{
		celebrityRepo: celebrityRepo,
	}
}

func (u *getCelebrityByIdUsecase) Execute(ctx context.Context, id int64) (*entities.Celebrity, error) {
	return u.celebrityRepo.GetCelebrityByID(ctx, id)
}
