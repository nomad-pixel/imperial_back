package entities

import "time"

type Car struct {
	ID             int64
	Name           string
	ImageUrl       string
	OnlyWithDriver bool
	PricePerDay    int64
	Tags           []*CarTag
	Mark           *CarMark
	Category       *CarCategory
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
