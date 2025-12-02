package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type driverRepository struct {
	db *pgxpool.Pool
}

func NewDriverRepository(db *pgxpool.Pool) ports.DriverRepository {
	return &driverRepository{db: db}
}

func (r *driverRepository) CreateDriver(ctx context.Context, driver *entities.Driver) error {
	query := `
		INSERT INTO drivers (full_name, about, photo_url, experience_years, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query,
		driver.FullName,
		driver.About,
		driver.PhotoURL,
		driver.ExperienceYears,
		driver.CreatedAt,
		driver.UpdatedAt,
	).Scan(&driver.ID)
}

func (r *driverRepository) GetDriverByID(ctx context.Context, id int64) (*entities.Driver, error) {
	query := `
		SELECT id, full_name, about, photo_url, experience_years, created_at, updated_at
		FROM drivers
		WHERE id = $1
	`
	driver := &entities.Driver{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&driver.ID,
		&driver.FullName,
		&driver.About,
		&driver.PhotoURL,
		&driver.ExperienceYears,
		&driver.CreatedAt,
		&driver.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return driver, nil
}

func (r *driverRepository) ListDrivers(ctx context.Context, offset, limit int64) (int64, []*entities.Driver, error) {
	countQuery := `SELECT COUNT(*) FROM drivers`
	var total int64
	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return 0, nil, err
	}

	query := `
		SELECT id, full_name, about, photo_url, experience_years, created_at, updated_at
		FROM drivers
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	drivers := make([]*entities.Driver, 0)
	for rows.Next() {
		driver := &entities.Driver{}
		if err := rows.Scan(
			&driver.ID,
			&driver.FullName,
			&driver.About,
			&driver.PhotoURL,
			&driver.ExperienceYears,
			&driver.CreatedAt,
			&driver.UpdatedAt,
		); err != nil {
			return 0, nil, err
		}
		drivers = append(drivers, driver)
	}

	if err := rows.Err(); err != nil {
		return 0, nil, err
	}

	return total, drivers, nil
}

func (r *driverRepository) UpdateDriver(ctx context.Context, driver *entities.Driver) error {
	query := `
		UPDATE drivers
		SET full_name = $1, about = $2, photo_url = $3, experience_years = $4, updated_at = $5
		WHERE id = $6
	`
	result, err := r.db.Exec(ctx, query,
		driver.FullName,
		driver.About,
		driver.PhotoURL,
		driver.ExperienceYears,
		driver.UpdatedAt,
		driver.ID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("driver not found")
	}

	return nil
}

func (r *driverRepository) DeleteDriver(ctx context.Context, id int64) error {
	query := `DELETE FROM drivers WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("driver not found")
	}

	return nil
}
