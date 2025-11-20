package car_usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type deleteCarTagUsecase struct {
	carTagRepo ports.CarTagRepository
}

type DeleteCarTagUsecase interface {
	Execute(ctx context.Context, tagID int64) error
}

func NewDeleteCarTagUsecase(carTagRepo ports.CarTagRepository) DeleteCarTagUsecase {
	return &deleteCarTagUsecase{carTagRepo: carTagRepo}
}

func (u *deleteCarTagUsecase) Execute(ctx context.Context, tagID int64) error {
	err := u.carTagRepo.DeleteCarTag(ctx, tagID)
	if err != nil {
		return apperrors.New(apperrors.ErrCodeBadRequest, "failed to delete car tag")
	}
	return nil
}
