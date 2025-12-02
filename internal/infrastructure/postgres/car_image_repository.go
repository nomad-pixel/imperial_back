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

func (r *carImageRepositoryImpl) Save(ctx context.Context, imageUrl string) (*entities.CarImage, error) {
	query := `
		INSERT INTO car_images (image_path)
		VALUES ($1)
		RETURNING id, image_path, created_at
	`
	var image entities.CarImage
	err := r.db.QueryRow(ctx, query, imageUrl).
		Scan(&image.ID, &image.ImagePath, &image.CreatedAt)

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
		SELECT id, image_path, created_at
		FROM car_images
		WHERE id = $1
	`

	var image entities.CarImage
	err := r.db.QueryRow(ctx, query, imageID).
		Scan(&image.ID, &image.ImagePath, &image.CreatedAt)

	if err != nil {
		return nil, apperrors.Wrap(err, apperrors.ErrCodeNotFound, "car image not found")
	}

	return &image, nil
}

func (r *carImageRepositoryImpl) GetList(ctx context.Context, offset, limit int64) (int64, []*entities.CarImage, error) {
	queryCount := `SELECT COUNT(*) FROM car_images`

	var totalCount int64
	err := r.db.QueryRow(ctx, queryCount).Scan(&totalCount)
	if err != nil {
		return 0, nil, apperrors.Wrap(err, apperrors.ErrCodeInternal, "failed to count car images")
	}
	query := `
		SELECT id, image_path, created_at
		FROM car_images
		ORDER BY created_at DESC
		OFFSET $1
		LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, offset, limit)
	if err != nil {
		return 0, nil, apperrors.Wrap(err, apperrors.ErrCodeInternal, "failed to list car images")
	}
	defer rows.Close()

	var images []*entities.CarImage
	for rows.Next() {
		var image entities.CarImage
		if err := rows.Scan(&image.ID, &image.ImagePath, &image.CreatedAt); err != nil {
			return 0, nil, apperrors.Wrap(err, apperrors.ErrCodeInternal, "failed to scan car image")
		}
		images = append(images, &image)
	}

	if err := rows.Err(); err != nil {
		return 0, nil, apperrors.Wrap(err, apperrors.ErrCodeInternal, "error occurred during rows iteration")
	}

	return totalCount, images, nil
}
