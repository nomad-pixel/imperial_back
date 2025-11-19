package ports

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

type CarMarkRepository interface {
	CreateCarMark(ctx context.Context, name string) (*entities.CarMark, error)
	GetCarMarkByID(ctx context.Context, id int64) (*entities.CarMark, error)
	UpdateCarMark(ctx context.Context, id int64, name string) (*entities.CarMark, error)
	DeleteCarMark(ctx context.Context, id int64) error
	ListCarMarks(ctx context.Context, offset int64, limit int64) (int64, []*entities.CarMark, error)
}
