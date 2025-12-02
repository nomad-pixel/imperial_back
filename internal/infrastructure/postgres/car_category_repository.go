package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type CarCategoryRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewCarCategoryRepositoryImpl(db *pgxpool.Pool) ports.CarCategoryRepository {
	return &CarCategoryRepositoryImpl{db: db}
}

func (r CarCategoryRepositoryImpl) CreateCarCategory(ctx context.Context, name string) (*entities.CarCategory, error) {
	query := `
		INSERT INTO car_categories (name)
		VALUES ($1)
		RETURNING id, name, created_at, updated_at
	`
	var category entities.CarCategory
	err := r.db.QueryRow(ctx, query, name).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r CarCategoryRepositoryImpl) GetCarCategoryByID(ctx context.Context, id int64) (*entities.CarCategory, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM car_categories
		WHERE id = $1
	`
	var category entities.CarCategory
	err := r.db.QueryRow(ctx, query, id).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r CarCategoryRepositoryImpl) UpdateCarCategory(ctx context.Context, id int64, name string) (*entities.CarCategory, error) {
	query := `
		UPDATE car_categories
		SET name = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, name, created_at, updated_at
	`
	var category entities.CarCategory
	err := r.db.QueryRow(ctx, query, name, id).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r CarCategoryRepositoryImpl) DeleteCarCategory(ctx context.Context, id int64) error {
	query := `
		DELETE FROM car_categories
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r CarCategoryRepositoryImpl) ListCarCategories(ctx context.Context, offset int64, limit int64) (int64, []*entities.CarCategory, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT id, name, created_at, updated_at
		FROM car_categories
		ORDER BY created_at DESC
		OFFSET $1
		LIMIT $2
	`
	rows, err := r.db.Query(ctx, query, offset, limit)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	var categories []*entities.CarCategory
	for rows.Next() {
		var category entities.CarCategory
		err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return 0, nil, err
		}
		categories = append(categories, &category)
	}

	var total int64
	countQuery := `SELECT COUNT(*) FROM car_categories`
	err = r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return 0, nil, err
	}

	return total, categories, nil
}
