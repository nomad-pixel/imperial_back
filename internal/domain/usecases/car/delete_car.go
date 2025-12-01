package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type deleteCarUsecase struct {
	carRepo ports.CarRepository
}

type DeleteCarUsecase interface {
	Execute(ctx context.Context, carID int64) error
}

func NewDeleteCarUsecase(carRepo ports.CarRepository) DeleteCarUsecase {
	return &deleteCarUsecase{carRepo: carRepo}
}

func (u *deleteCarUsecase) Execute(ctx context.Context, carID int64) error {
	err := u.carRepo.DeleteCar(ctx, carID)
	if err != nil {
		return apperrors.New(apperrors.ErrCodeBadRequest, "failed to delete car")
	}
	return nil
}
