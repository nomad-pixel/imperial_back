package car_usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type createCarMarkUsecase struct {
	carMarkRepo ports.CarMarkRepository
}

type CreateCarMarkUsecase interface {
	Execute(ctx context.Context, name string) (*entities.CarMark, error)
}

func NewCreateCarMarkUsecase(carMarkRepo ports.CarMarkRepository) CreateCarMarkUsecase {
	return &createCarMarkUsecase{carMarkRepo: carMarkRepo}
}

func (u *createCarMarkUsecase) Execute(ctx context.Context, name string) (*entities.CarMark, error) {
	carMark, err := u.carMarkRepo.CreateCarMark(ctx, name)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, "failed to create car mark")
	}

	return carMark, nil
}
