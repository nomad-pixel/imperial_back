package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type getCarMarkUsecase struct {
	carMarkRepo ports.CarMarkRepository
}

type GetCarMarkUsecase interface {
	Execute(ctx context.Context, markID int64) (*entities.CarMark, error)
}

func NewGetCarMarkUsecase(carMarkRepo ports.CarMarkRepository) GetCarMarkUsecase {
	return &getCarMarkUsecase{carMarkRepo: carMarkRepo}
}

func (u *getCarMarkUsecase) Execute(ctx context.Context, markID int64) (*entities.CarMark, error) {
	carMark, err := u.carMarkRepo.GetCarMarkByID(ctx, markID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeNotFound, "car mark not found")
	}
	return carMark, nil
}
