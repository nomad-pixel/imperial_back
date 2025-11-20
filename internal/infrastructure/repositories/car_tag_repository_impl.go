package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type CarTagRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewCarTagRepositoryImpl(db *pgxpool.Pool) ports.CarTagRepository {
	return &CarTagRepositoryImpl{db: db}
}

func (r CarTagRepositoryImpl) CreateCarTag(ctx context.Context, name string) (*entities.CarTag, error) {
	query := `
		INSERT INTO car_tags (name)
		VALUES ($1)
		RETURNING id, name, created_at, updated_at
	`

	var tag entities.CarTag
	err := r.db.QueryRow(ctx, query, name).Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r CarTagRepositoryImpl) UpdateCarTag(ctx context.Context, tagId int64, name string) (*entities.CarTag, error) {
	query := `
		UPDATE car_tags
		SET name = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, name, created_at, updated_at
	`
	var tag entities.CarTag
	err := r.db.QueryRow(ctx, query, name, tagId).Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r CarTagRepositoryImpl) GetCarTagById(ctx context.Context, tagId int64) (*entities.CarTag, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM car_tags
		WHERE id = $1
	`
	var tag entities.CarTag
	err := r.db.QueryRow(ctx, query, tagId).Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r CarTagRepositoryImpl) DeleteCarTag(ctx context.Context, id int64) error {
	query := `
		DELETE FROM car_tags
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r CarTagRepositoryImpl) ListCarTags(ctx context.Context, offset int64, limit int64) (int64, []*entities.CarTag, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	var total int64
	countQuery := `SELECT COUNT(*) FROM car_tags`
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return 0, nil, err
	}

	query := `
		SELECT id, name, created_at, updated_at
		FROM car_tags
		ORDER BY created_at DESC
		OFFSET $1 LIMIT $2
	`
	rows, err := r.db.Query(ctx, query, offset, limit)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	var tags []*entities.CarTag
	for rows.Next() {
		var tag entities.CarTag
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
			return 0, nil, err
		}
		tags = append(tags, &tag)
	}

	return total, tags, nil
}
