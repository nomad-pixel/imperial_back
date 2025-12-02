package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type listDriversUsecase struct {
	driverRepo ports.DriverRepository
}

type ListDriversUsecase interface {
	Execute(ctx context.Context, offset, limit int64) (int64, []*entities.Driver, error)
}

func NewListDriversUsecase(driverRepo ports.DriverRepository) ListDriversUsecase {
	return &listDriversUsecase{driverRepo: driverRepo}
}

func (u *listDriversUsecase) Execute(ctx context.Context, offset, limit int64) (int64, []*entities.Driver, error) {
	return u.driverRepo.ListDrivers(ctx, offset, limit)
}
