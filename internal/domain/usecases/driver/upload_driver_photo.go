package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type uploadDriverPhotoUsecase struct {
	driverRepo   ports.DriverRepository
	imageService ports.ImageService
}

type UploadDriverPhotoUsecase interface {
	Execute(ctx context.Context, driverID int64, imageData []byte, fileName string) (*entities.Driver, error)
}

func NewUploadDriverPhotoUsecase(driverRepo ports.DriverRepository, imageService ports.ImageService) UploadDriverPhotoUsecase {
	return &uploadDriverPhotoUsecase{
		driverRepo:   driverRepo,
		imageService: imageService,
	}
}

func (u *uploadDriverPhotoUsecase) Execute(ctx context.Context, driverID int64, fileData []byte, fileName string) (*entities.Driver, error) {
	if len(fileData) == 0 {
		return nil, apperrors.New(apperrors.ErrCodeValidation, "photo file is empty")
	}

	driver, err := u.driverRepo.GetDriverByID(ctx, driverID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeNotFound, "driver not found")
	}

	if driver.PhotoURL != "" {
		err = u.imageService.DeleteImage(driver.PhotoURL)
		if err != nil {
			return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to delete old driver photo")
		}
	}

	imagePath, err := u.imageService.SaveImage(fileData, "drivers", fileName)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to save photo")
	}

	driver.SetPhotoURL(imagePath)

	if err := u.driverRepo.UpdateDriver(ctx, driver); err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to update driver")
	}

	return driver, nil
}
