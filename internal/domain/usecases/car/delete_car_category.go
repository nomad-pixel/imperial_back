package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type deleteCarCategoryUsecase struct {
	carCategoryRepo ports.CarCategoryRepository
}

type DeleteCarCategoryUsecase interface {
	Execute(ctx context.Context, categoryID int64) error
}

func NewDeleteCarCategoryUsecase(carCategoryRepo ports.CarCategoryRepository) DeleteCarCategoryUsecase {
	return &deleteCarCategoryUsecase{carCategoryRepo: carCategoryRepo}
}

func (u *deleteCarCategoryUsecase) Execute(ctx context.Context, categoryID int64) error {
	err := u.carCategoryRepo.DeleteCarCategory(ctx, categoryID)
	if err != nil {
		return apperrors.New(apperrors.ErrCodeBadRequest, "failed to delete car category")
	}
	return nil
}
