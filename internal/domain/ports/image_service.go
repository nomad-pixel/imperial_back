package ports

type ImageService interface {
	SaveImage(fileData []byte, folderName, fileName string) (string, error)
	DeleteImage(imagePath string) error
	GetFullImagePath(imagePath string) string
}
