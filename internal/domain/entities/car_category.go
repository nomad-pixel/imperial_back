package entities

import (
	"errors"
	"strings"
	"time"
)

type CarCategory struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCarCategory(name string) (*CarCategory, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New("car category name cannot be empty")
	}

	if len(name) < 1 {
		return nil, errors.New("car category name must be at least 1 character")
	}

	if len(name) > 100 {
		return nil, errors.New("car category name cannot exceed 100 characters")
	}

	now := time.Now()
	return &CarCategory{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (cc *CarCategory) Validate() error {
	if cc.ID <= 0 {
		return errors.New("invalid car category ID")
	}

	name := strings.TrimSpace(cc.Name)
	if name == "" {
		return errors.New("car category name cannot be empty")
	}

	if len(name) < 1 || len(name) > 100 {
		return errors.New("car category name must be between 1 and 100 characters")
	}

	return nil
}

func (cc *CarCategory) SetName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New("car category name cannot be empty")
	}

	if len(name) < 1 || len(name) > 100 {
		return errors.New("car category name must be between 1 and 100 characters")
	}

	cc.Name = name
	cc.UpdatedAt = time.Now()
	return nil
}
