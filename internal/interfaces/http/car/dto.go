package car

import "github.com/nomad-pixel/imperial/internal/domain/entities"

type CreateCarRequest struct {
	Name           string  `json:"name" binding:"required,min=1,max=255" example:"Toyota"`
	PricePerDay    int64   `json:"price_per_day" binding:"required,min=0" example:"100"`
	OnlyWithDriver bool    `json:"only_with_driver" example:"false"`
	MarkId         int64   `json:"mark_id" binding:"required,min=1" example:"1"`
	CategoryId     int64   `json:"category_id" binding:"required,min=1" example:"2"`
	TagsIds        []int64 `json:"tags_ids" binding:"required"`
}

type UpdateCarRequest struct {
	Name           string  `json:"name" binding:"required,min=1,max=255" example:"Toyota"`
	PricePerDay    int64   `json:"price_per_day" binding:"required,min=0" example:"100"`
	OnlyWithDriver bool    `json:"only_with_driver" example:"false"`
	MarkId         int64   `json:"mark_id" binding:"required,min=1" example:"1"`
	CategoryId     int64   `json:"category_id" binding:"required,min=1" example:"2"`
	TagsIds        []int64 `json:"tags_ids" binding:"required"`
}

type CarResponse = entities.Car

type MessageResponse struct {
	Message string `json:"message"`
}

type ListCarsResponse struct {
	Total int64           `json:"total"`
	Data  []*entities.Car `json:"data"`
}
