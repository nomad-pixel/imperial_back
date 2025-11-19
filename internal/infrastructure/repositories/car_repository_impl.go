package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type CarRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewCarRepositoryImpl(db *pgxpool.Pool) ports.CarRepository {
	return &CarRepositoryImpl{db: db}
}

func (r CarRepositoryImpl) CreateCar(ctx context.Context, car *entities.Car) (*entities.Car, error) {

	return nil, nil
}

func (r CarRepositoryImpl) GetCarsByMark(ctx context.Context, markId int64) (int64, []*entities.Car, error) {
	// Implementation goes here
	return 0, nil, nil
}

func (r CarRepositoryImpl) GetCarsByName(ctx context.Context, name string) (int64, []*entities.Car, error) {
	// Implementation goes here
	return 0, nil, nil
}

func (r CarRepositoryImpl) GetCarByID(ctx context.Context, id int64) (*entities.Car, error) {
	// Implementation goes here
	return nil, nil
}

func (r CarRepositoryImpl) UpdateCar(ctx context.Context, car *entities.Car) (*entities.Car, error) {
	// Implementation goes here
	return nil, nil
}

func (r CarRepositoryImpl) DeleteCar(ctx context.Context, id int64) error {
	// Implementation goes here
	return nil
}

func (r CarRepositoryImpl) ListCars(ctx context.Context, offset int64, limit int64) (int64, []*entities.Car, error) {
	// Implementation goes here
	return 0, nil, nil
}
