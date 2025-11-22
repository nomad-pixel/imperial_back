package car_usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type getCarImagesListUsecase struct {
	carImageRepo ports.CarImageRepository
	imageService ports.ImageService
}

type GetCarImagesListUsecase interface {
	Execute(ctx context.Context, offset, limit int64) (total int64, images []*entities.CarImage, err error)
}

func NewGetCarImagesListUsecase(carImageRepository ports.CarImageRepository, imageService ports.ImageService) GetCarImagesListUsecase {
	return &getCarImagesListUsecase{
		carImageRepo: carImageRepository,
		imageService: imageService,
	}
}

func (u *getCarImagesListUsecase) Execute(ctx context.Context, offset, limit int64) (int64, []*entities.CarImage, error) {
	total, carImages, err := u.carImageRepo.GetList(ctx, offset, limit)
	if err != nil {
		return 0, nil, apperrors.New(apperrors.ErrCodeInternal, "failed to get car images list from repository")
	}

	for _, image := range carImages {
		image.ImagePath = u.imageService.GetFullImagePath(image.ImagePath)
	}

	return total, carImages, nil
}
