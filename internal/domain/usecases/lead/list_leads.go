package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type listLeadsUsecase struct {
	leadRepo ports.LeadRepository
}

type ListLeadsUsecase interface {
	Execute(ctx context.Context, offset, limit int64) (int64, []*entities.Lead, error)
}

func NewListLeadsUsecase(leadRepo ports.LeadRepository) ListLeadsUsecase {
	return &listLeadsUsecase{leadRepo: leadRepo}
}

func (u *listLeadsUsecase) Execute(ctx context.Context, offset, limit int64) (int64, []*entities.Lead, error) {
	return u.leadRepo.ListLeads(ctx, offset, limit)
}
