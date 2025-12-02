package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type deleteCelebrityUsecase struct {
	celebrityRepo ports.CelebrityRepository
}

type DeleteCelebrityUsecase interface {
	Execute(ctx context.Context, id int64) error
}

func NewDeleteCelebrityUsecase(celebrityRepo ports.CelebrityRepository) DeleteCelebrityUsecase {
	return &deleteCelebrityUsecase{celebrityRepo: celebrityRepo}
}

func (u *deleteCelebrityUsecase) Execute(ctx context.Context, id int64) error {
	err := u.celebrityRepo.DeleteCelebrity(ctx, id)
	if err != nil {
		return apperrors.New(apperrors.ErrCodeBadRequest, "failed to delete celebrity")
	}
	return nil
}
