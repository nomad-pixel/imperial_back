package car_usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type updateCarMarkUsecase struct {
	carMarkRepo ports.CarMarkRepository
}

type UpdateCarMarkUsecase interface {
	Execute(ctx context.Context, markID int64, name string) (*entities.CarMark, error)
}

func NewUpdateCarMarkUsecase(carMarkRepo ports.CarMarkRepository) UpdateCarMarkUsecase {
	return &updateCarMarkUsecase{carMarkRepo: carMarkRepo}
}

func (u *updateCarMarkUsecase) Execute(ctx context.Context, markID int64, name string) (*entities.CarMark, error) {
	carMark, err := u.carMarkRepo.UpdateCarMark(ctx, markID, name)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, "failed to update car mark")
	}
	return carMark, nil
}
