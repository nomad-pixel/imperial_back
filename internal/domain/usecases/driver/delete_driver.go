package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type deleteDriverUsecase struct {
	driverRepo ports.DriverRepository
}

type DeleteDriverUsecase interface {
	Execute(ctx context.Context, id int64) error
}

func NewDeleteDriverUsecase(driverRepo ports.DriverRepository) DeleteDriverUsecase {
	return &deleteDriverUsecase{driverRepo: driverRepo}
}

func (u *deleteDriverUsecase) Execute(ctx context.Context, id int64) error {
	return u.driverRepo.DeleteDriver(ctx, id)
}
