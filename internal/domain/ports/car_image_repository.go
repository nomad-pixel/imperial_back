package ports

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

type CarImageRepository interface {
	Save(ctx context.Context, imageUrl string) (*entities.CarImage, error)
	Delete(ctx context.Context, imageID int64) error
	GetByID(ctx context.Context, imageID int64) (*entities.CarImage, error)
	GetList(ctx context.Context, offset, limit int64) (int64, []*entities.CarImage, error)
}
