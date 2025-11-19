package ports

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

type CarTagRepository interface {
	CreateCarTag(ctx context.Context, name string) (*entities.CarTag, error)
	UpdateCarTag(ctx context.Context, id int64, name string) (*entities.CarTag, error)
	GetCarTagById(ctx context.Context, id int64) (*entities.CarTag, error)
	DeleteCarTag(ctx context.Context, id int64) error
	ListCarTags(ctx context.Context, offset int64, limit int64) (int64, []*entities.CarTag, error)
}
