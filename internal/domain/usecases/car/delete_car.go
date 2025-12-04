package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type deleteCarUsecase struct {
	carRepo      ports.CarRepository
	carImageRepo ports.CarImageRepository
	imageService ports.ImageService
}

type DeleteCarUsecase interface {
	Execute(ctx context.Context, carID int64) error
}

func NewDeleteCarUsecase(
	carRepo ports.CarRepository,
	carImageRepo ports.CarImageRepository,
	imageService ports.ImageService,
) DeleteCarUsecase {
	return &deleteCarUsecase{
		carRepo:      carRepo,
		carImageRepo: carImageRepo,
		imageService: imageService,
	}
}

func (u *deleteCarUsecase) Execute(ctx context.Context, carID int64) error {
	images, err := u.carImageRepo.GetListByCar(ctx, carID)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeInternal, "failed to get car images")
	}

	for _, image := range images {
		if err := u.imageService.DeleteImage(image.ImagePath); err != nil {
			return apperrors.Wrap(err, apperrors.ErrCodeInternal, "failed to delete car image file")
		}

		if err := u.carImageRepo.Delete(ctx, image.ID); err != nil {
			return apperrors.Wrap(err, apperrors.ErrCodeInternal, "failed to delete car image from database")
		}
	}

	err = u.carRepo.DeleteCar(ctx, carID)
	if err != nil {
		return apperrors.New(apperrors.ErrCodeBadRequest, "failed to delete car")
	}
	return nil
}
