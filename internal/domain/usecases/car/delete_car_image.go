package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type deleteCarImageUsecase struct {
	carImageRepo ports.CarImageRepository
	imageService ports.ImageService
}

type DeleteCarImageUsecase interface {
	Execute(ctx context.Context, imageID int64) error
}

func NewDeleteCarImageUsecase(carImageRepository ports.CarImageRepository, imageService ports.ImageService) DeleteCarImageUsecase {
	return &deleteCarImageUsecase{
		carImageRepo: carImageRepository,
		imageService: imageService,
	}
}

func (u *deleteCarImageUsecase) Execute(ctx context.Context, imageID int64) error {
	image, err := u.carImageRepo.GetByID(ctx, imageID)
	if err != nil {
		return apperrors.New(apperrors.ErrCodeInternal, "failed to get car image from repository")
	}

	err = u.imageService.DeleteImage(image.ImagePath)

	if err != nil {
		return apperrors.New(apperrors.ErrCodeInternal, "failed to delete car image file")
	}

	err = u.carImageRepo.Delete(ctx, imageID)
	if err != nil {
		return apperrors.New(apperrors.ErrCodeInternal, "failed to delete car image from repository")
	}
	return nil
}
