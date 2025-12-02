package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type createDriverUsecase struct {
	driverRepo ports.DriverRepository
}

type CreateDriverUsecase interface {
	Execute(ctx context.Context, fullName, about, experienceYears string) (*entities.Driver, error)
}

func NewCreateDriverUsecase(driverRepo ports.DriverRepository) CreateDriverUsecase {
	return &createDriverUsecase{driverRepo: driverRepo}
}

func (u *createDriverUsecase) Execute(ctx context.Context, fullName, about, experienceYears string) (*entities.Driver, error) {
	driver, err := entities.NewDriver(fullName, about, experienceYears)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, err.Error())
	}

	err = u.driverRepo.CreateDriver(ctx, driver)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, "failed to create driver")
	}

	return driver, nil
}
