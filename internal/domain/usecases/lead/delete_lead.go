package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type deleteLeadUsecase struct {
	leadRepo ports.LeadRepository
}

type DeleteLeadUsecase interface {
	Execute(ctx context.Context, id int64) error
}

func NewDeleteLeadUsecase(leadRepo ports.LeadRepository) DeleteLeadUsecase {
	return &deleteLeadUsecase{leadRepo: leadRepo}
}

func (u *deleteLeadUsecase) Execute(ctx context.Context, id int64) error {
	return u.leadRepo.DeleteLead(ctx, id)
}
