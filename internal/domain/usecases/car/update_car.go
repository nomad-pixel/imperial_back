package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type updateCarUsecase struct {
	carRepo ports.CarRepository
}

type UpdateCarUsecase interface {
	Execute(ctx context.Context, car *entities.Car) (*entities.Car, error)
}

func NewUpdateCarUsecase(carRepo ports.CarRepository) UpdateCarUsecase {
	return &updateCarUsecase{carRepo: carRepo}
}
func (u *updateCarUsecase) Execute(ctx context.Context, car *entities.Car) (*entities.Car, error) {
	err := u.carRepo.UpdateCar(ctx, car)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, "failed to update car")
	}
	return car, nil
}
