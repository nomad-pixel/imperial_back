package ports

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

type DriverRepository interface {
	CreateDriver(ctx context.Context, driver *entities.Driver) error
	GetDriverByID(ctx context.Context, id int64) (*entities.Driver, error)
	ListDrivers(ctx context.Context, offset, limit int64) (int64, []*entities.Driver, error)
	UpdateDriver(ctx context.Context, driver *entities.Driver) error
	DeleteDriver(ctx context.Context, id int64) error
}
