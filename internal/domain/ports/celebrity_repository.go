package ports

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

type CelebrityRepository interface {
	CreateCelebrity(ctx context.Context, celebrity *entities.Celebrity) error
	UploadImage(ctx context.Context, id int64, imageData []byte, fileName string) (*entities.Celebrity, error)
	UpdateCelebrity(ctx context.Context, celebrity *entities.Celebrity) error
	GetCelebrityByID(ctx context.Context, id int64) (*entities.Celebrity, error)
	DeleteCelebrity(ctx context.Context, id int64) error
	ListCelebrities(ctx context.Context, offset int64, limit int64) (int64, []*entities.Celebrity, error)
}
