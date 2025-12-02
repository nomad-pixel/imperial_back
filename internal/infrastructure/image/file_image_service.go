package image

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type FileImageService struct {
	storagePath string
	baseURL     string
}

func NewFileImageService(storagePath, baseURL string) (ports.ImageService, error) {
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, errors.Wrap(err, errors.ErrCodeInternal, "failed to create storage directory")
	}
	return &FileImageService{
		storagePath: storagePath,
		baseURL:     baseURL,
	}, nil
}

func (s *FileImageService) SaveImage(fileData []byte, folderName, fileName string) (string, error) {
	carDir := filepath.Join(s.storagePath, folderName)
	if err := os.MkdirAll(carDir, 0755); err != nil {
		return "", errors.Wrap(err, errors.ErrCodeInternal, "failed to create"+folderName+" directory")
	}
	ext := filepath.Ext(fileName)
	timestamp := time.Now().Unix()
	newFileName := fmt.Sprintf("%d_%s%s", timestamp, strings.TrimSuffix(fileName, ext), ext)
	filePath := filepath.Join(carDir, newFileName)
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return "", errors.Wrap(err, errors.ErrCodeInternal, "failed to save image file")
	}
	relativePath := filepath.Join(folderName, newFileName)
	return relativePath, nil
}

func (s *FileImageService) DeleteImage(imagePath string) error {
	fullPath := filepath.Join(s.storagePath, imagePath)
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return errors.New(errors.ErrCodeNotFound, "image not found")
		}
		return errors.Wrap(err, errors.ErrCodeInternal, "failed to delete image")
	}
	return nil
}

func (s *FileImageService) GetFullImagePath(imagePath string) string {
	// imagePath это "cars/filename.png"
	// return "images/cars/filename.png"
	return filepath.Join("images", imagePath)
}
