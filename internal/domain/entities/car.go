package entities

import "time"

type Car struct {
	ID             int64        `json:"id"`
	Name           string       `json:"name"`
	ImageUrl       string       `json:"image_url"`
	OnlyWithDriver bool         `json:"only_with_driver"`
	PricePerDay    int64        `json:"price_per_day"`
	Tags           []*CarTag    `json:"tags"`
	Mark           *CarMark     `json:"mark"`
	Category       *CarCategory `json:"category"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}
