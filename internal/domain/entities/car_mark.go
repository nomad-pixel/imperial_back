package entities

import (
	"errors"
	"strings"
	"time"
)

type CarMark struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCarMark(name string) (*CarMark, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New("car mark name cannot be empty")
	}

	if len(name) < 1 {
		return nil, errors.New("car mark name must be at least 1 character")
	}

	if len(name) > 100 {
		return nil, errors.New("car mark name cannot exceed 100 characters")
	}

	now := time.Now()
	return &CarMark{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (cm *CarMark) Validate() error {
	if cm.ID <= 0 {
		return errors.New("invalid car mark ID")
	}

	name := strings.TrimSpace(cm.Name)
	if name == "" {
		return errors.New("car mark name cannot be empty")
	}

	if len(name) < 1 || len(name) > 100 {
		return errors.New("car mark name must be between 1 and 100 characters")
	}

	return nil
}

func (cm *CarMark) SetName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New("car mark name cannot be empty")
	}

	if len(name) < 1 || len(name) > 100 {
		return errors.New("car mark name must be between 1 and 100 characters")
	}

	cm.Name = name
	cm.UpdatedAt = time.Now()
	return nil
}
