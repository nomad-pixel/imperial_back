package entities

import (
	"errors"
	"strings"
	"time"
)

type Celebrity struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewCelebrity(name string) (*Celebrity, error) {
	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New("celebrity name cannot be empty")
	}

	if len(name) < 2 {
		return nil, errors.New("celebrity name must be at least 2 characters")
	}

	if len(name) > 200 {
		return nil, errors.New("celebrity name cannot exceed 200 characters")
	}

	now := time.Now()
	return &Celebrity{
		Name:      name,
		Image:     "",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (c *Celebrity) Validate() error {
	if c.ID <= 0 {
		return errors.New("invalid celebrity ID")
	}

	name := strings.TrimSpace(c.Name)
	if name == "" {
		return errors.New("celebrity name cannot be empty")
	}

	if len(name) < 2 || len(name) > 200 {
		return errors.New("celebrity name must be between 2 and 200 characters")
	}

	return nil
}

func (c *Celebrity) SetName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New("celebrity name cannot be empty")
	}

	if len(name) < 2 || len(name) > 200 {
		return errors.New("celebrity name must be between 2 and 200 characters")
	}

	c.Name = name
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Celebrity) SetImage(imageURL string) {
	c.Image = imageURL
	c.UpdatedAt = time.Now()
}
