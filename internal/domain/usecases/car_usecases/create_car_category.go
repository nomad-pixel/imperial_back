package car_usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type createCarCategoryUsecase struct {
	carCategoryRepo ports.CarCategoryRepository
}

type CreateCarCategoryUsecase interface {
	Execute(ctx context.Context, name string) (*entities.CarCategory, error)
}

func NewCreateCarCategoryUsecase(carCategoryRepo ports.CarCategoryRepository) CreateCarCategoryUsecase {
	return &createCarCategoryUsecase{carCategoryRepo: carCategoryRepo}
}

func (u *createCarCategoryUsecase) Execute(ctx context.Context, name string) (*entities.CarCategory, error) {
	result, err := u.carCategoryRepo.CreateCarCategory(ctx, name)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, "failed to create car category")
	}

	return result, nil
}
