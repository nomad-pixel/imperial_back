package ports

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

type LeadRepository interface {
	CreateLead(ctx context.Context, lead *entities.Lead) error
	GetLeadByID(ctx context.Context, id int64) (*entities.Lead, error)
	ListLeads(ctx context.Context, offset, limit int64) (int64, []*entities.Lead, error)
	DeleteLead(ctx context.Context, id int64) error
}
