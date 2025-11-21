package car_usecases

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
	Execute(ctx context.Context, fileData []byte, fileName string) (*entities.CarImage, error)
}

func NewCreateCarImageUsecase(carImageRepo ports.CarImageRepository, imageService ports.ImageService) CreateCarImageUsecase {
	return &createCarImageUsecase{
		carImageRepo: carImageRepo,
		imageService: imageService,
	}
}

func (u *createCarImageUsecase) Execute(ctx context.Context, fileData []byte, fileName string) (*entities.CarImage, error) {

	if len(fileData) == 0 {
		return nil, apperrors.New(apperrors.ErrCodeValidation, "image file is empty")
	}

	imagePath, err := u.imageService.SaveCarImage(fileData, fileName)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to save image: %v")
	}
	var carImage *entities.CarImage
	carImage, err = u.carImageRepo.Save(ctx, imagePath)

	if err != nil {
		err = u.imageService.DeleteCarImage(imagePath)
		if err != nil {
			return nil, apperrors.Wrap(err, apperrors.ErrCodeInternal, "failed to cleanup image after DB failure")
		}
		return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to create car image record: %v")
	}

	// Преобразуем путь для ответа пользователю: cars/file.png -> images/cars/file.png
	carImage.ImagePath = u.imageService.GetFullImagePath(carImage.ImagePath)

	return carImage, nil
}
