package entities

import (
	"errors"
	"strings"
	"time"
)

type CarTag struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCarTag(name string) (*CarTag, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New("car tag name cannot be empty")
	}

	if len(name) < 1 {
		return nil, errors.New("car tag name must be at least 1 character")
	}

	if len(name) > 100 {
		return nil, errors.New("car tag name cannot exceed 100 characters")
	}

	now := time.Now()
	return &CarTag{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (ct *CarTag) Validate() error {
	if ct.ID <= 0 {
		return errors.New("invalid car tag ID")
	}

	name := strings.TrimSpace(ct.Name)
	if name == "" {
		return errors.New("car tag name cannot be empty")
	}

	if len(name) < 1 || len(name) > 100 {
		return errors.New("car tag name must be between 1 and 100 characters")
	}

	return nil
}

func (ct *CarTag) SetName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New("car tag name cannot be empty")
	}

	if len(name) < 1 || len(name) > 100 {
		return errors.New("car tag name must be between 1 and 100 characters")
	}

	ct.Name = name
	ct.UpdatedAt = time.Now()
	return nil
}
