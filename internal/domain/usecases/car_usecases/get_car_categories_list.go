package car_usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type getCarCategoriesListUsecase struct {
	carCategoryRepo ports.CarCategoryRepository
}

type GetCarCategoriesListUsecase interface {
	Execute(ctx context.Context, offset, limit int64) (int64, []*entities.CarCategory, error)
}

func NewGetCarCategoriesListUsecase(carCategoryRepo ports.CarCategoryRepository) GetCarCategoriesListUsecase {
	return &getCarCategoriesListUsecase{carCategoryRepo: carCategoryRepo}
}

func (u *getCarCategoriesListUsecase) Execute(ctx context.Context, offset, limit int64) (int64, []*entities.CarCategory, error) {
	total, carCategories, err := u.carCategoryRepo.ListCarCategories(ctx, offset, limit)
	if err != nil {
		return 0, nil, apperrors.New(apperrors.ErrCodeInternal, "failed to get car categories list")
	}
	return total, carCategories, nil
}
