package ports

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

type CarRepository interface {
	CreateCar(ctx context.Context, car *entities.Car) (*entities.Car, error)
	GetCarsByMark(ctx context.Context, markId int64) (int64, []*entities.Car, error)
	GetCarsByName(ctx context.Context, name string) (int64, []*entities.Car, error)
	GetCarByID(ctx context.Context, id int64) (*entities.Car, error)
	UpdateCar(ctx context.Context, car *entities.Car) (*entities.Car, error)
	DeleteCar(ctx context.Context, id int64) error
	ListCars(ctx context.Context, offset int64, limit int64) (int64, []*entities.Car, error)
}
