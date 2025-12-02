package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type updateCelebrityUsecase struct {
	celebrityRepo ports.CelebrityRepository
}

type UpdateCelebrityUsecase interface {
	Execute(ctx context.Context, celebrity *entities.Celebrity) (*entities.Celebrity, error)
}

func NewUpdateCelebrityUsecase(celebrityRepo ports.CelebrityRepository) UpdateCelebrityUsecase {
	return &updateCelebrityUsecase{celebrityRepo: celebrityRepo}
}

func (u *updateCelebrityUsecase) Execute(ctx context.Context, celebrity *entities.Celebrity) (*entities.Celebrity, error) {
	err := u.celebrityRepo.UpdateCelebrity(ctx, celebrity)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, "failed to update celebrity")
	}

	updatedCelebrity, err := u.celebrityRepo.GetCelebrityByID(ctx, celebrity.ID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeNotFound, "failed to get updated celebrity")
	}

	return updatedCelebrity, nil
}
