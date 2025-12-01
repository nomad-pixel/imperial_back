package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type createCarTagUsecase struct {
	carTagRepo ports.CarTagRepository
}

type CreateCarTagUsecase interface {
	Execute(ctx context.Context, name string) (*entities.CarTag, error)
}

func NewCreateCarTagUsecase(carTagRepo ports.CarTagRepository) CreateCarTagUsecase {
	return &createCarTagUsecase{carTagRepo: carTagRepo}
}

func (u *createCarTagUsecase) Execute(ctx context.Context, name string) (*entities.CarTag, error) {
	carTag, err := u.carTagRepo.CreateCarTag(ctx, name)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, "failed to create car tag")
	}

	return carTag, nil
}
