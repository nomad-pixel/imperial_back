package car

import (
	"time"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

// swagger:model CreateCarRequest
// CreateCarRequest represents payload to create a car
type CreateCarRequest struct {
	Name           string  `json:"name" binding:"required,name" example:"Toyota"`
	ImageUrl       string  `json:"image_url" binding:"required,image_url" example:"https://example.com/car.jpg"`
	PricePerDay    int64   `json:"price_per_day" binding:"required" example:"100"`
	OnlyWithDriver bool    `json:"only_with_driver" binding:"required" example:"false"`
	MarkId         int64   `json:"mark_id" binding:"required" example:"1"`
	CategoryId     int64   `json:"category_id" binding:"required" example:"2"`
	TagsIds        []int64 `json:"tags_ids" binding:"required"`
}

type CarResponse struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	ImageUrl       string    `json:"imageUrl"`
	OnlyWithDriver bool      `json:"onlyWithDriver"`
	PricePerDay    int64     `json:"pricePerDay"`
	MarkId         int64     `json:"markId"`
	CategoryId     int64     `json:"categoryId"`
	TagsIds        []int64   `json:"tagsIds"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func NewCarResponse(car *entities.Car) CarResponse {
	resp := CarResponse{
		ID:             car.ID,
		Name:           car.Name,
		ImageUrl:       car.ImageUrl,
		OnlyWithDriver: car.OnlyWithDriver,
		PricePerDay:    car.PricePerDay,
		CreatedAt:      car.CreatedAt,
		UpdatedAt:      car.UpdatedAt,
	}

	if car.Mark != nil {
		resp.MarkId = car.Mark.ID
	}
	if car.Category != nil {
		resp.CategoryId = car.Category.ID
	}
	if len(car.Tags) > 0 {
		resp.TagsIds = make([]int64, 0, len(car.Tags))
		for _, t := range car.Tags {
			if t == nil {
				continue
			}
			resp.TagsIds = append(resp.TagsIds, t.ID)
		}
	}

	return resp
}
