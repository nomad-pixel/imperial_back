package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type getLeadByIdUsecase struct {
	leadRepo ports.LeadRepository
}

type GetLeadByIdUsecase interface {
	Execute(ctx context.Context, id int64) (*entities.Lead, error)
}

func NewGetLeadByIdUsecase(leadRepo ports.LeadRepository) GetLeadByIdUsecase {
	return &getLeadByIdUsecase{leadRepo: leadRepo}
}

func (u *getLeadByIdUsecase) Execute(ctx context.Context, id int64) (*entities.Lead, error) {
	return u.leadRepo.GetLeadByID(ctx, id)
}
