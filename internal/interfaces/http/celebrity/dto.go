package celebrity

import "github.com/nomad-pixel/imperial/internal/domain/entities"

type MessageResponse struct {
	Message string `json:"message"`
}

type ListCelebritiesResponse struct {
	Total int64                 `json:"total"`
	Data  []*entities.Celebrity `json:"data"`
}

type CreateCelebrityRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=255" example:"John Doe"`
	Image string `json:"image" example:"https://example.com/celebrity.jpg"`
}

type UpdateCelebrityRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=255" example:"John Doe"`
	Image string `json:"image" example:"https://example.com/celebrity.jpg"`
}
