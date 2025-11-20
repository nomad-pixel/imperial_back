package car_usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type updateCarTagUsecase struct {
	carTagRepo ports.CarTagRepository
}

type UpdateCarTagUsecase interface {
	Execute(ctx context.Context, carID int64, name string) (*entities.CarTag, error)
}

func NewUpdateCarTagUsecase(carTagRepo ports.CarTagRepository) UpdateCarTagUsecase {
	return &updateCarTagUsecase{carTagRepo: carTagRepo}
}

func (u *updateCarTagUsecase) Execute(ctx context.Context, carID int64, name string) (*entities.CarTag, error) {
	car, err := u.carTagRepo.UpdateCarTag(ctx, carID, name)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, "failed to update car")
	}
	return car, nil
}
