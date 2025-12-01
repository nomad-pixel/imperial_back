package ports

import "github.com/nomad-pixel/imperial/internal/domain/entities"

type CelebrityRepository interface {
	CreateCelebrity(name string, imagePath string) (*entities.Celebrity, error)
	UploadImage(imageData []byte, fileName string) (string, error)
	UpdateCelebrity(id int64, name string, imagePath string) (*entities.Celebrity, error)
	GetCelebrityByID(id int64) (*entities.Celebrity, error)
	DeleteCelebrity(id int64) error
	ListCelebrities(offset int64, limit int64) (int64, []*entities.Celebrity, error)
}
