package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type CelebrityRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewCelebrityRepositoryImpl(db *pgxpool.Pool) ports.CelebrityRepository {
	return &CelebrityRepositoryImpl{db: db}
}

func (r *CelebrityRepositoryImpl) CreateCelebrity(ctx context.Context, celebrity *entities.Celebrity) error {
	query := `
		INSERT INTO celebrities (name, image)
		VALUES ($1, $2)
		RETURNING id, name, image, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query, celebrity.Name, celebrity.Image).Scan(
		&celebrity.ID,
		&celebrity.Name,
		&celebrity.Image,
		&celebrity.CreatedAt,
		&celebrity.UpdatedAt,
	)
}

func (r *CelebrityRepositoryImpl) UploadImage(ctx context.Context, id int64, imagePath string) (*entities.Celebrity, error) {
	query := `
		UPDATE celebrities
		SET image = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, name, image, created_at, updated_at
	`
	var celebrity entities.Celebrity
	err := r.db.QueryRow(ctx, query, imagePath, id).Scan(
		&celebrity.ID,
		&celebrity.Name,
		&celebrity.Image,
		&celebrity.CreatedAt,
		&celebrity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &celebrity, nil
}

func (r *CelebrityRepositoryImpl) UpdateCelebrity(ctx context.Context, celebrity *entities.Celebrity) error {
	query := `
		UPDATE celebrities
		SET name = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, name, image, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query, celebrity.Name, celebrity.ID).Scan(
		&celebrity.ID,
		&celebrity.Name,
		&celebrity.Image,
		&celebrity.CreatedAt,
		&celebrity.UpdatedAt,
	)
}

func (r *CelebrityRepositoryImpl) GetCelebrityByID(ctx context.Context, id int64) (*entities.Celebrity, error) {
	query := `
		SELECT id, name, image, created_at, updated_at
		FROM celebrities
		WHERE id = $1
	`
	var celebrity entities.Celebrity
	err := r.db.QueryRow(ctx, query, id).Scan(
		&celebrity.ID,
		&celebrity.Name,
		&celebrity.Image,
		&celebrity.CreatedAt,
		&celebrity.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &celebrity, nil
}

func (r *CelebrityRepositoryImpl) DeleteCelebrity(ctx context.Context, id int64) error {
	query := `DELETE FROM celebrities WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return nil
	}
	return nil
}

func (r *CelebrityRepositoryImpl) ListCelebrities(ctx context.Context, offset int64, limit int64) (int64, []*entities.Celebrity, error) {
	countQuery := `SELECT COUNT(*) FROM celebrities`
	var total int64
	err := r.db.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		return 0, nil, err
	}

	query := `
		SELECT id, name, image, created_at, updated_at
		FROM celebrities
		ORDER BY created_at DESC
		OFFSET $1 LIMIT $2
	`

	rows, err := r.db.Query(ctx, query, offset, limit)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	var celebrities []*entities.Celebrity
	for rows.Next() {
		var celebrity entities.Celebrity
		err := rows.Scan(
			&celebrity.ID,
			&celebrity.Name,
			&celebrity.Image,
			&celebrity.CreatedAt,
			&celebrity.UpdatedAt,
		)
		if err != nil {
			return 0, nil, err
		}
		celebrities = append(celebrities, &celebrity)
	}

	if err = rows.Err(); err != nil {
		return 0, nil, err
	}

	return total, celebrities, nil
}
