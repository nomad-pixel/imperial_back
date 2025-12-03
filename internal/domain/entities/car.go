package entities

import (
	"errors"
	"strings"
	"time"
)

type Car struct {
	ID             int64        `json:"id"`
	Name           string       `json:"name"`
	OnlyWithDriver bool         `json:"only_with_driver"`
	PricePerDay    int64        `json:"price_per_day"`
	Tags           []*CarTag    `json:"tags"`
	Mark           *CarMark     `json:"mark"`
	Category       *CarCategory `json:"category"`
	Images         []*CarImage  `json:"images"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

func NewCar(name string, pricePerDay int64, markID, categoryID int64, onlyWithDriver bool) (*Car, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New("car name cannot be empty")
	}

	if len(name) < 2 {
		return nil, errors.New("car name must be at least 2 characters")
	}

	if len(name) > 200 {
		return nil, errors.New("car name cannot exceed 200 characters")
	}

	if pricePerDay < 0 {
		return nil, errors.New("price per day cannot be negative")
	}

	if markID <= 0 {
		return nil, errors.New("mark ID must be positive")
	}

	if categoryID <= 0 {
		return nil, errors.New("category ID must be positive")
	}

	now := time.Now()
	return &Car{
		Name:           name,
		PricePerDay:    pricePerDay,
		OnlyWithDriver: onlyWithDriver,
		Mark:           &CarMark{ID: markID},
		Category:       &CarCategory{ID: categoryID},
		Tags:           make([]*CarTag, 0),
		Images:         make([]*CarImage, 0),
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (c *Car) Validate() error {
	if c.ID <= 0 {
		return errors.New("invalid car ID")
	}

	name := strings.TrimSpace(c.Name)
	if name == "" {
		return errors.New("car name cannot be empty")
	}

	if len(name) < 2 || len(name) > 200 {
		return errors.New("car name must be between 2 and 200 characters")
	}

	if c.PricePerDay < 0 {
		return errors.New("price per day cannot be negative")
	}

	if c.Mark == nil || c.Mark.ID <= 0 {
		return errors.New("car must have a valid mark")
	}

	if c.Category == nil || c.Category.ID <= 0 {
		return errors.New("car must have a valid category")
	}

	return nil
}

func (c *Car) SetName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New("car name cannot be empty")
	}

	if len(name) < 2 || len(name) > 200 {
		return errors.New("car name must be between 2 and 200 characters")
	}

	c.Name = name
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Car) SetPricePerDay(price int64) error {
	if price < 0 {
		return errors.New("price per day cannot be negative")
	}

	c.PricePerDay = price
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Car) SetOnlyWithDriver(onlyWithDriver bool) {
	c.OnlyWithDriver = onlyWithDriver
	c.UpdatedAt = time.Now()
}

func (c *Car) AddTag(tag *CarTag) error {
	if tag == nil || tag.ID <= 0 {
		return errors.New("invalid tag")
	}

	for _, existingTag := range c.Tags {
		if existingTag.ID == tag.ID {
			return errors.New("tag already exists")
		}
	}

	c.Tags = append(c.Tags, tag)
	c.UpdatedAt = time.Now()
	return nil
}
