package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type getCarByIdUsecase struct {
	carRepo ports.CarRepository
}

type GetCarByIdUsecase interface {
	Execute(ctx context.Context, carID int64) (*entities.Car, error)
}

func NewGetCarByIdUsecase(carRepo ports.CarRepository) GetCarByIdUsecase {
	return &getCarByIdUsecase{carRepo: carRepo}
}

func (u *getCarByIdUsecase) Execute(ctx context.Context, carID int64) (*entities.Car, error) {
	car, err := u.carRepo.GetCarByID(ctx, carID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeNotFound, "car not found")
	}
	return car, nil
}
