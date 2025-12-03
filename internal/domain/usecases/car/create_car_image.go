package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type createCarImageUsecase struct {
	carImageRepo ports.CarImageRepository
	imageService ports.ImageService
}

type CreateCarImageUsecase interface {
	Execute(ctx context.Context, carID int64, fileData []byte, fileName string) (*entities.CarImage, error)
}

func NewCreateCarImageUsecase(carImageRepo ports.CarImageRepository, imageService ports.ImageService) CreateCarImageUsecase {
	return &createCarImageUsecase{
		carImageRepo: carImageRepo,
		imageService: imageService,
	}
}

func (u *createCarImageUsecase) Execute(ctx context.Context, carID int64, fileData []byte, fileName string) (*entities.CarImage, error) {

	if len(fileData) == 0 {
		return nil, apperrors.New(apperrors.ErrCodeValidation, "image file is empty")
	}

	imagePath, err := u.imageService.SaveImage(fileData, "cars", fileName)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to save image: %v")
	}
	var carImage *entities.CarImage
	carImage, err = u.carImageRepo.Save(ctx, carID, imagePath)

	if err != nil {
		err = u.imageService.DeleteImage(imagePath)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.ErrCodeInternal, "failed to cleanup image after DB failure")
		}
		return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to create car image record: %v")
	}

	return carImage, nil
}
