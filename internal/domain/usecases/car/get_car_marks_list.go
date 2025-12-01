package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type getCarMarksListUsecase struct {
	carMarkRepo ports.CarMarkRepository
}

type GetCarMarksListUsecase interface {
	Execute(ctx context.Context, offset, limit int64) (int64, []*entities.CarMark, error)
}

func NewGetCarMarksListUsecase(carMarkRepo ports.CarMarkRepository) GetCarMarksListUsecase {
	return &getCarMarksListUsecase{carMarkRepo: carMarkRepo}
}

func (u *getCarMarksListUsecase) Execute(ctx context.Context, offset, limit int64) (int64, []*entities.CarMark, error) {
	total, carMarks, err := u.carMarkRepo.ListCarMarks(ctx, offset, limit)
	if err != nil {
		return 0, nil, apperrors.New(apperrors.ErrCodeInternal, "failed to get car marks list")
	}
	return total, carMarks, nil
}
