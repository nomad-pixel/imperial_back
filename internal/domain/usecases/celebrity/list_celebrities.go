package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type listCelebritiesUsecase struct {
	celebrityRepo ports.CelebrityRepository
}

type ListCelebritiesUsecase interface {
	Execute(ctx context.Context, offset int64, limit int64) (int64, []*entities.Celebrity, error)
}

func NewListCelebritiesUsecase(celebrityRepo ports.CelebrityRepository) ListCelebritiesUsecase {
	return &listCelebritiesUsecase{
		celebrityRepo: celebrityRepo,
	}
}

func (u *listCelebritiesUsecase) Execute(ctx context.Context, offset int64, limit int64) (int64, []*entities.Celebrity, error) {
	return u.celebrityRepo.ListCelebrities(ctx, offset, limit)
}
