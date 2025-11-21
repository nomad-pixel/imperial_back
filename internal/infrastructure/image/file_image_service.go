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

func (s *FileImageService) SaveCarImage(fileData []byte, fileName string) (string, error) {
	carDir := filepath.Join(s.storagePath, "car")
	if err := os.MkdirAll(carDir, 0755); err != nil {
		return "", errors.Wrap(err, errors.ErrCodeInternal, "failed to create car directory")
	}
	ext := filepath.Ext(fileName)
	timestamp := time.Now().Unix()
	newFileName := fmt.Sprintf("%d_%s%s", timestamp, strings.TrimSuffix(fileName, ext), ext)
	filePath := filepath.Join(carDir, newFileName)
	if err := os.WriteFile(filePath, fileData, 0644); err != nil {
		return "", errors.Wrap(err, errors.ErrCodeInternal, "failed to save image file")
	}
	relativePath := filepath.Join("car", newFileName)
	return relativePath, nil
}

func (s *FileImageService) DeleteCarImage(imagePath string) error {
	fullPath := filepath.Join(s.storagePath, imagePath)
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return errors.New(errors.ErrCodeNotFound, "image not found")
		}
		return errors.Wrap(err, errors.ErrCodeInternal, "failed to delete image")
	}
	return nil
}
