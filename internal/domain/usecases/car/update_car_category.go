package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type updateCarCategoryUsecase struct {
	carCategoryRepo ports.CarCategoryRepository
}

type UpdateCarCategoryUsecase interface {
	Execute(ctx context.Context, categoryID int64, name string) (*entities.CarCategory, error)
}

func NewUpdateCarCategoryUsecase(carCategoryRepo ports.CarCategoryRepository) UpdateCarCategoryUsecase {
	return &updateCarCategoryUsecase{carCategoryRepo: carCategoryRepo}
}

func (u *updateCarCategoryUsecase) Execute(ctx context.Context, categoryID int64, name string) (*entities.CarCategory, error) {
	carCategory, err := u.carCategoryRepo.UpdateCarCategory(ctx, categoryID, name)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, "failed to update car category")
	}
	return carCategory, nil
}
