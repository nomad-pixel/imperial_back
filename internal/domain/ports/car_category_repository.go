package ports

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

type CarCategoryRepository interface {
	CreateCarCategory(ctx context.Context, name string) (*entities.CarCategory, error)
	GetCarCategoryByID(ctx context.Context, id int64) (*entities.CarCategory, error)
	UpdateCarCategory(ctx context.Context, id int64, name string) (*entities.CarCategory, error)
	DeleteCarCategory(ctx context.Context, id int64) error
	ListCarCategories(ctx context.Context, offset int64, limit int64) (int64, []*entities.CarCategory, error)
}
