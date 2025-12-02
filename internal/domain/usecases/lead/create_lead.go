package usecases

import (
	"context"
	"time"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type createLeadUsecase struct {
	leadRepo ports.LeadRepository
}

type CreateLeadUsecase interface {
	Execute(ctx context.Context, fullName, phone string, startDate, endDate time.Time) (*entities.Lead, error)
}

func NewCreateLeadUsecase(leadRepo ports.LeadRepository) CreateLeadUsecase {
	return &createLeadUsecase{leadRepo: leadRepo}
}

func (u *createLeadUsecase) Execute(ctx context.Context, fullName, phone string, startDate, endDate time.Time) (*entities.Lead, error) {
	lead, err := entities.NewLead(fullName, phone, startDate, endDate)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, err.Error())
	}

	err = u.leadRepo.CreateLead(ctx, lead)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeBadRequest, "failed to create lead")
	}

	return lead, nil
}
