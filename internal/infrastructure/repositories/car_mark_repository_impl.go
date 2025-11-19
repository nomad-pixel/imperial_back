package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type CarMarkRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewCarMarkRepositoryImpl(db *pgxpool.Pool) ports.CarMarkRepository {
	return &CarMarkRepositoryImpl{db: db}
}

func (r CarMarkRepositoryImpl) ListCarMarks(ctx context.Context, offset int64, limit int64) (int64, []*entities.CarMark, error) {
	if limit == 0 {
		limit = 20
	}

	var total int64
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM car_marks").Scan(&total)
	if err != nil {
		return 0, nil, err
	}

	query := `
		SELECT id, name, created_at, updated_at
		FROM car_marks
		ORDER BY id
		OFFSET $1 LIMIT $2
	`
	rows, err := r.db.Query(ctx, query, offset, limit)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	var carMarks []*entities.CarMark
	for rows.Next() {
		var mark entities.CarMark
		err := rows.Scan(&mark.ID, &mark.Name, &mark.CreatedAt, &mark.UpdatedAt)
		if err != nil {
			return 0, nil, err
		}
		carMarks = append(carMarks, &mark)
	}

	return total, carMarks, nil
}

func (r CarMarkRepositoryImpl) CreateCarMark(ctx context.Context, name string) (*entities.CarMark, error) {
	query := `
		INSERT INTO car_marks (name)
		VALUES ($1)
		RETURNING id, name, created_at, updated_at
	`
	var mark entities.CarMark
	err := r.db.QueryRow(ctx, query, name).Scan(&mark.ID, &mark.Name, &mark.CreatedAt, &mark.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &mark, nil
}

func (r CarMarkRepositoryImpl) GetCarMarkByID(ctx context.Context, id int64) (*entities.CarMark, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM car_marks
		WHERE id = $1
	`
	var mark entities.CarMark
	err := r.db.QueryRow(ctx, query, id).Scan(&mark.ID, &mark.Name, &mark.CreatedAt, &mark.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &mark, nil
}

func (r CarMarkRepositoryImpl) UpdateCarMark(ctx context.Context, id int64, name string) (*entities.CarMark, error) {
	query := `
		UPDATE car_marks
		SET name = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, name, created_at, updated_at
	`
	var mark entities.CarMark
	err := r.db.QueryRow(ctx, query, name, id).Scan(&mark.ID, &mark.Name, &mark.CreatedAt, &mark.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &mark, nil
}

func (r CarMarkRepositoryImpl) DeleteCarMark(ctx context.Context, id int64) error {
	query := `
		DELETE FROM car_marks
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
