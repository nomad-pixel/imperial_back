package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type uploadCelebrityImageUsecase struct {
	celebrityRepo ports.CelebrityRepository
	imageService  ports.ImageService
}

type UploadCelebrityImageUsecase interface {
	Execute(ctx context.Context, id int64, imageData []byte, fileName string) (*entities.Celebrity, error)
}

func NewUploadCelebrityImageUsecase(celebrityRepo ports.CelebrityRepository, imageService ports.ImageService) UploadCelebrityImageUsecase {
	return &uploadCelebrityImageUsecase{
		celebrityRepo: celebrityRepo,
		imageService:  imageService,
	}
}

func (u *uploadCelebrityImageUsecase) Execute(ctx context.Context, id int64, fileData []byte, fileName string) (*entities.Celebrity, error) {
	if len(fileData) == 0 {
		return nil, apperrors.New(apperrors.ErrCodeValidation, "image file is empty")
	}

	var celebrity *entities.Celebrity
	var err error

	celebrity, err = u.celebrityRepo.GetCelebrityByID(ctx, id)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "celebrity not founds")
	}
	if celebrity.Image != "" {
		err = u.imageService.DeleteImage(celebrity.Image)
		if err != nil {
			return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to delete old celebrity image")
		}
	}

	imagePath, err := u.imageService.SaveImage(fileData, "celebrities", fileName)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to save image")
	}

	celebrity, err = u.celebrityRepo.UploadImage(ctx, id, imagePath)

	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "failed to upload celebrity image")
	}
	return celebrity, nil
}
