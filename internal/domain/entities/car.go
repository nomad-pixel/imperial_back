package entities

import "time"

type Car struct {
	ID             int64
	Name           string
	ImageUrl       string
	Tags           []*CarTag
	Mark           *CarMark
	Category       *CarCategory
	OnlyWithDriver bool
	PricePerDay    int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
