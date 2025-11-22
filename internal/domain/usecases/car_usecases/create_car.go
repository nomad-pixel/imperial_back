package car_usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type createCarUsecase struct {
	carRepo ports.CarRepository
}

type CreateCarUsecase interface {
	Execute(ctx context.Context, car *entities.Car) (*entities.Car, error)
}

func NewCreateCarUsecase(carRepo ports.CarRepository) CreateCarUsecase {
	return &createCarUsecase{carRepo: carRepo}
}

func (u *createCarUsecase) Execute(ctx context.Context, car *entities.Car) (*entities.Car, error) {
	err := u.carRepo.CreateCar(ctx, car)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, "failed to create car")
	}
	return car, nil
}
