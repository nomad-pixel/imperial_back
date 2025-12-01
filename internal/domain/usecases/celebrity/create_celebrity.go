package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type createCelebrityUsecase struct {
	celebrityRepo ports.CelebrityRepository
}

type CreateCelebrityUsecase interface {
	Execute(ctx context.Context, celebrity *entities.Celebrity) error
}

func NewCreateCelebrityUsecase(celebrityRepo ports.CelebrityRepository) CreateCelebrityUsecase {
	return &createCelebrityUsecase{celebrityRepo: celebrityRepo}
}

func (u *createCelebrityUsecase) Execute(ctx context.Context, celebrity *entities.Celebrity) error {
	err := u.celebrityRepo.CreateCelebrity(ctx, celebrity)
	if err != nil {
		return apperrors.New(apperrors.ErrCodeBadRequest, "failed to create celebrity")
	}
	return nil
}
