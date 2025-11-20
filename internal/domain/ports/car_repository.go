package ports

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

type CarRepository interface {
	CreateCar(ctx context.Context, car *entities.Car) (*entities.Car, error)
	GetCarByID(ctx context.Context, id int64) (*entities.Car, error)
	UpdateCar(ctx context.Context, car *entities.Car) (*entities.Car, error)
	DeleteCar(ctx context.Context, id int64) error
	ListCars(ctx context.Context, offset, limit int64, name string, markID int64) (int64, []*entities.Car, error)
}
