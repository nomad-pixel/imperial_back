package ports

type ImageService interface {
	SaveCarImage(fileData []byte, fileName string) (string, error)
	DeleteCarImage(imagePath string) error
}
