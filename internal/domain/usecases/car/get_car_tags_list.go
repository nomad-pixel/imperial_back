package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type getCarTagsListUsecase struct {
	carTagRepo ports.CarTagRepository
}

type GetCarTagsListUsecase interface {
	Execute(ctx context.Context, offset, limit int64) (int64, []*entities.CarTag, error)
}

func NewGetCarTagsListUsecase(carTagRepo ports.CarTagRepository) GetCarTagsListUsecase {
	return &getCarTagsListUsecase{carTagRepo: carTagRepo}
}

func (u *getCarTagsListUsecase) Execute(ctx context.Context, offset, limit int64) (int64, []*entities.CarTag, error) {
	total, carTags, err := u.carTagRepo.ListCarTags(ctx, 0, 1000)
	if err != nil {
		return 0, nil, apperrors.New(apperrors.ErrCodeInternal, "failed to get car tags list")
	}
	return total, carTags, nil
}
