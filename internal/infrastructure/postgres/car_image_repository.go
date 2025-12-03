package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type carImageRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewCarImageRepositoryImpl(db *pgxpool.Pool) ports.CarImageRepository {
	return &carImageRepositoryImpl{db: db}
}

func (r *carImageRepositoryImpl) Save(ctx context.Context, carID int64, imageUrl string) (*entities.CarImage, error) {
	query := `
		INSERT INTO car_images (car_id, image_path)
		VALUES ($1, $2)
		RETURNING id, car_id, image_path, created_at
	`
	var image entities.CarImage
	err := r.db.QueryRow(ctx, query, carID, imageUrl).
		Scan(&image.ID, &image.CarID, &image.ImagePath, &image.CreatedAt)

	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrCodeInternal, "failed to save car image")
	}

	return &image, nil
}

func (r *carImageRepositoryImpl) Delete(ctx context.Context, imageID int64) error {
	query := `DELETE FROM car_images WHERE id = $1`

	result, err := r.db.Exec(ctx, query, imageID)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeInternal, "failed to delete car image")
	}

	if result.RowsAffected() == 0 {
		return apperrors.New(apperrors.ErrCodeNotFound, "car image not found")
	}

	return nil
}

func (r *carImageRepositoryImpl) GetByID(ctx context.Context, imageID int64) (*entities.CarImage, error) {
	query := `
		SELECT id, car_id, image_path, created_at
		FROM car_images
		WHERE id = $1
	`

	var image entities.CarImage
	err := r.db.QueryRow(ctx, query, imageID).
		Scan(&image.ID, &image.CarID, &image.ImagePath, &image.CreatedAt)

	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrCodeNotFound, "car image not found")
	}

	return &image, nil
}

func (r *carImageRepositoryImpl) GetList(ctx context.Context, carID int64, offset, limit int64) (int64, []*entities.CarImage, error) {
	queryCount := `SELECT COUNT(*) FROM car_images WHERE car_id = $1`

	var totalCount int64
	err := r.db.QueryRow(ctx, queryCount, carID).Scan(&totalCount)
	if err != nil {
		return 0, nil, apperrors.Wrap(err, apperrors.ErrCodeInternal, "failed to count car images")
	}
	query := `
		SELECT id, car_id, image_path, created_at
		FROM car_images
		WHERE car_id = $1
		ORDER BY created_at ASC
		OFFSET $2
		LIMIT $3
	`

	rows, err := r.db.Query(ctx, query, carID, offset, limit)
	if err != nil {
		return 0, nil, apperrors.Wrap(err, apperrors.ErrCodeInternal, "failed to list car images")
	}
	defer rows.Close()

	var images []*entities.CarImage
	for rows.Next() {
		var image entities.CarImage
		if err := rows.Scan(&image.ID, &image.CarID, &image.ImagePath, &image.CreatedAt); err != nil {
			return 0, nil, apperrors.Wrap(err, apperrors.ErrCodeInternal, "failed to scan car image")
		}
		images = append(images, &image)
	}

	if err := rows.Err(); err != nil {
		return 0, nil, apperrors.Wrap(err, apperrors.ErrCodeInternal, "error occurred during rows iteration")
	}

	return totalCount, images, nil
}

func (r *carImageRepositoryImpl) GetListByCar(ctx context.Context, carID int64) ([]*entities.CarImage, error) {
	query := `
		SELECT id, car_id, image_path, created_at
		FROM car_images
		WHERE car_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, carID)
	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrCodeInternal, "failed to list car images")
	}
	defer rows.Close()

	var images []*entities.CarImage
	for rows.Next() {
		var image entities.CarImage
		if err := rows.Scan(&image.ID, &image.CarID, &image.ImagePath, &image.CreatedAt); err != nil {
			return nil, apperrors.Wrap(err, apperrors.ErrCodeInternal, "failed to scan car image")
		}
		images = append(images, &image)
	}

	if err := rows.Err(); err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrCodeInternal, "error occurred during rows iteration")
	}

	return images, nil
}
