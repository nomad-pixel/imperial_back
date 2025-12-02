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
	Execute(ctx context.Context, id int64, name string) (*entities.Celebrity, error)
}

func NewUpdateCelebrityUsecase(celebrityRepo ports.CelebrityRepository) UpdateCelebrityUsecase {
	return &updateCelebrityUsecase{celebrityRepo: celebrityRepo}
}

func (u *updateCelebrityUsecase) Execute(ctx context.Context, id int64, name string) (*entities.Celebrity, error) {
	celebrity, err := u.celebrityRepo.GetCelebrityByID(ctx, id)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeNotFound, "celebrity not found")
	}
	if err := celebrity.SetName(name); err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, err.Error())
	}
	if err := celebrity.Validate(); err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, err.Error())
	}

	err = u.celebrityRepo.UpdateCelebrity(ctx, celebrity)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, "failed to update celebrity")
	}

	return celebrity, nil
}
