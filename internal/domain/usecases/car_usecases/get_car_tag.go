package car_usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type getCarTagUsecase struct {
	carTagRepo ports.CarTagRepository
}

type GetCarTagUsecase interface {
	Execute(ctx context.Context, tagID int64) (*entities.CarTag, error)
}

func NewGetCarTagUsecase(carTagRepo ports.CarTagRepository) GetCarTagUsecase {
	return &getCarTagUsecase{carTagRepo: carTagRepo}
}

func (u *getCarTagUsecase) Execute(ctx context.Context, tagID int64) (*entities.CarTag, error) {
	carTag, err := u.carTagRepo.GetCarTagById(ctx, tagID)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeNotFound, "car tag not found")
	}
	return carTag, nil
}
