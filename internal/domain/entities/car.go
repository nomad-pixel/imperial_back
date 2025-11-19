package entities

import "time"

type Car struct {
	ID             int64
	Name           string
	ImageUrl       string
	Tags           []CarTag
	Mark           CarMark
	CarCategory    CarCategory
	OnlyWithDriver bool
	PricePerDay    float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
