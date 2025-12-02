package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type deleteCelebrityUsecase struct {
	celebrityRepo ports.CelebrityRepository
	imageService  ports.ImageService
}

type DeleteCelebrityUsecase interface {
	Execute(ctx context.Context, id int64) error
}

func NewDeleteCelebrityUsecase(celebrityRepo ports.CelebrityRepository, imageService ports.ImageService) DeleteCelebrityUsecase {
	return &deleteCelebrityUsecase{
		celebrityRepo: celebrityRepo,
		imageService:  imageService,
	}
}

func (u *deleteCelebrityUsecase) Execute(ctx context.Context, id int64) error {
	celebrity, err := u.celebrityRepo.GetCelebrityByID(ctx, id)
	if err != nil {
		return apperrors.New(apperrors.ErrCodeNotFound, "celebrity not found")
	}
	if celebrity.Image != "" {
		err = u.imageService.DeleteImage(celebrity.Image)
		if err != nil {
			return apperrors.New(apperrors.ErrCodeInternal, "failed to delete celebrity image")
		}
	}
	err = u.celebrityRepo.DeleteCelebrity(ctx, id)
	if err != nil {
		return apperrors.New(apperrors.ErrCodeBadRequest, "failed to delete celebrity")
	}
	return nil
}
