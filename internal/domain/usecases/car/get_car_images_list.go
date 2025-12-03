package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type getCarImagesListUsecase struct {
	carImageRepo ports.CarImageRepository
}

type GetCarImagesListUsecase interface {
	Execute(ctx context.Context, carID int64, offset, limit int64) (total int64, images []*entities.CarImage, err error)
}

func NewGetCarImagesListUsecase(carImageRepository ports.CarImageRepository) GetCarImagesListUsecase {
	return &getCarImagesListUsecase{
		carImageRepo: carImageRepository,
	}
}

func (u *getCarImagesListUsecase) Execute(ctx context.Context, carID int64, offset, limit int64) (int64, []*entities.CarImage, error) {
	total, carImages, err := u.carImageRepo.GetList(ctx, carID, offset, limit)
	if err != nil {
		return 0, nil, apperrors.New(apperrors.ErrCodeInternal, "failed to get car images list from repository")
	}

	return total, carImages, nil
}
