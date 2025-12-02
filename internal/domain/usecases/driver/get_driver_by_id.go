package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type getDriverByIdUsecase struct {
	driverRepo ports.DriverRepository
}

type GetDriverByIdUsecase interface {
	Execute(ctx context.Context, id int64) (*entities.Driver, error)
}

func NewGetDriverByIdUsecase(driverRepo ports.DriverRepository) GetDriverByIdUsecase {
	return &getDriverByIdUsecase{driverRepo: driverRepo}
}

func (u *getDriverByIdUsecase) Execute(ctx context.Context, id int64) (*entities.Driver, error) {
	return u.driverRepo.GetDriverByID(ctx, id)
}
