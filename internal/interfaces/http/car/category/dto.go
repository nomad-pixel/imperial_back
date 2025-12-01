package car

import "github.com/nomad-pixel/imperial/internal/domain/entities"

type CreateCarCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateCarCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type CarCategoryResponse = entities.CarCategory

type ListCarCategoriesResponse struct {
	Total int64                   `json:"total"`
	Data  []*entities.CarCategory `json:"data"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
