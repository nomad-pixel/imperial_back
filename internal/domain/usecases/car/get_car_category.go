package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type getCarCategoryUsecase struct {
	carCategoryRepo ports.CarCategoryRepository
}

type GetCarCategoryUsecase interface {
	Execute(ctx context.Context, categoryID int64) (*entities.CarCategory, error)
}

func NewGetCarCategoryUsecase(carCategoryRepo ports.CarCategoryRepository) GetCarCategoryUsecase {
	return &getCarCategoryUsecase{carCategoryRepo: carCategoryRepo}
}

func (u *getCarCategoryUsecase) Execute(ctx context.Context, categoryID int64) (*entities.CarCategory, error) {
	carCategory, err := u.carCategoryRepo.GetCarCategoryByID(ctx, categoryID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeNotFound, "car category not found")
	}
	return carCategory, nil
}
