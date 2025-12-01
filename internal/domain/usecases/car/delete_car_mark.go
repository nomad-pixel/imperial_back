package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type deleteCarMarkUsecase struct {
	carMarkRepo ports.CarMarkRepository
}

type DeleteCarMarkUsecase interface {
	Execute(ctx context.Context, markID int64) error
}

func NewDeleteCarMarkUsecase(carMarkRepo ports.CarMarkRepository) DeleteCarMarkUsecase {
	return &deleteCarMarkUsecase{carMarkRepo: carMarkRepo}
}

func (u *deleteCarMarkUsecase) Execute(ctx context.Context, markID int64) error {
	err := u.carMarkRepo.DeleteCarMark(ctx, markID)
	if err != nil {
		return apperrors.New(apperrors.ErrCodeBadRequest, "failed to delete car mark")
	}
	return nil
}
