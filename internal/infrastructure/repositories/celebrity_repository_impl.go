package repositories

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

func (r *CelebrityRepositoryImpl) UploadImage(ctx context.Context, id int64, imageData []byte, fileName string) (*entities.Celebrity, error) {
	// Implementation goes here
	return nil, nil
}

func (r *CelebrityRepositoryImpl) UpdateCelebrity(ctx context.Context, celebrity *entities.Celebrity) error {
	// Implementation goes here
	return nil
}

func (r *CelebrityRepositoryImpl) GetCelebrityByID(ctx context.Context, id int64) (*entities.Celebrity, error) {
	// Implementation goes here
	return nil, nil
}

func (r *CelebrityRepositoryImpl) DeleteCelebrity(ctx context.Context, id int64) error {
	// Implementation goes here
	return nil
}

func (r *CelebrityRepositoryImpl) ListCelebrities(ctx context.Context, offset int64, limit int64) (int64, []*entities.Celebrity, error) {
	// Implementation goes here
	return 0, nil, nil
}
