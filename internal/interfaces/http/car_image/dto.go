package car_image

import "github.com/nomad-pixel/imperial/internal/domain/entities"

type CarImageResponse = entities.CarImage

type MessageResponse struct {
	Message string `json:"message"`
}

type ListCarImagesResponse struct {
	Total int64                `json:"total"`
	Data  []*entities.CarImage `json:"data"`
}
