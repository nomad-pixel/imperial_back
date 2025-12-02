package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type updateDriverUsecase struct {
	driverRepo ports.DriverRepository
}

type UpdateDriverUsecase interface {
	Execute(ctx context.Context, id int64, fullName, about, experienceYears string) (*entities.Driver, error)
}

func NewUpdateDriverUsecase(driverRepo ports.DriverRepository) UpdateDriverUsecase {
	return &updateDriverUsecase{driverRepo: driverRepo}
}

func (u *updateDriverUsecase) Execute(ctx context.Context, id int64, fullName, about, experienceYears string) (*entities.Driver, error) {
	driver, err := u.driverRepo.GetDriverByID(ctx, id)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeNotFound, "driver not found")
	}

	if err := driver.SetFullName(fullName); err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, err.Error())
	}

	if err := driver.SetAbout(about); err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, err.Error())
	}

	if err := driver.SetExperienceYears(experienceYears); err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, err.Error())
	}

	if err := driver.Validate(); err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, err.Error())
	}

	err = u.driverRepo.UpdateDriver(ctx, driver)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, "failed to update driver")
	}

	return driver, nil
}
