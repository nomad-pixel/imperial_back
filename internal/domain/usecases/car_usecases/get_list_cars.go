package car_usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type getListCarsUsecase struct {
	carRepo ports.CarRepository
}

type GetListCarsUsecase interface {
	Execute(ctx context.Context, offset int64, limit int64, name string, markID int64, categoryID int64) (int64, []*entities.Car, error)
}

func NewGetListCarsUsecase(carRepo ports.CarRepository) GetListCarsUsecase {
	return &getListCarsUsecase{carRepo: carRepo}
}

func (u *getListCarsUsecase) Execute(ctx context.Context, offset int64, limit int64, name string, markID int64, categoryID int64) (int64, []*entities.Car, error) {
	total, cars, err := u.carRepo.ListCars(ctx, offset, limit, name, markID, categoryID)
	if err != nil {
		return 0, nil, apperrors.New(apperrors.ErrCodeInternal, "failed to get cars list")
	}
	return total, cars, nil
}
