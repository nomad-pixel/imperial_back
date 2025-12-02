package entities

import (
	"errors"
	"strings"
	"time"
)

type Driver struct {
	ID              int64     `json:"id"`
	FullName        string    `json:"full_name"`
	About           string    `json:"about"`
	PhotoURL        string    `json:"photo_url"`
	ExperienceYears string    `json:"experience_years"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func NewDriver(fullName, about, experienceYears string) (*Driver, error) {
	fullName = strings.TrimSpace(fullName)

	if fullName == "" {
		return nil, errors.New("full name cannot be empty")
	}

	if len(fullName) < 2 {
		return nil, errors.New("full name must be at least 2 characters")
	}

	if len(fullName) > 255 {
		return nil, errors.New("full name cannot exceed 255 characters")
	}

	about = strings.TrimSpace(about)
	if about == "" {
		return nil, errors.New("about cannot be empty")
	}

	if len(about) > 1000 {
		return nil, errors.New("about cannot exceed 1000 characters")
	}

	experienceYears = strings.TrimSpace(experienceYears)
	if experienceYears == "" {
		return nil, errors.New("experience years cannot be empty")
	}

	if len(experienceYears) > 300 {
		return nil, errors.New("experience years cannot exceed 300 characters")
	}

	now := time.Now()
	return &Driver{
		FullName:        fullName,
		About:           about,
		ExperienceYears: experienceYears,
		PhotoURL:        "",
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

func (d *Driver) Validate() error {
	if d.ID <= 0 {
		return errors.New("invalid driver ID")
	}

	fullName := strings.TrimSpace(d.FullName)
	if fullName == "" || len(fullName) < 2 || len(fullName) > 255 {
		return errors.New("full name must be between 2 and 255 characters")
	}

	about := strings.TrimSpace(d.About)
	if about == "" || len(about) > 1000 {
		return errors.New("about must not be empty and cannot exceed 1000 characters")
	}

	experienceYears := strings.TrimSpace(d.ExperienceYears)
	if experienceYears == "" || len(experienceYears) > 300 {
		return errors.New("experience years must not be empty and cannot exceed 300 characters")
	}

	return nil
}

func (d *Driver) SetFullName(fullName string) error {
	fullName = strings.TrimSpace(fullName)

	if fullName == "" || len(fullName) < 2 || len(fullName) > 255 {
		return errors.New("full name must be between 2 and 255 characters")
	}

	d.FullName = fullName
	d.UpdatedAt = time.Now()
	return nil
}

func (d *Driver) SetAbout(about string) error {
	about = strings.TrimSpace(about)

	if about == "" || len(about) > 1000 {
		return errors.New("about must not be empty and cannot exceed 1000 characters")
	}

	d.About = about
	d.UpdatedAt = time.Now()
	return nil
}

func (d *Driver) SetExperienceYears(experienceYears string) error {
	experienceYears = strings.TrimSpace(experienceYears)

	if experienceYears == "" || len(experienceYears) > 300 {
		return errors.New("experience years must not be empty and cannot exceed 300 characters")
	}

	d.ExperienceYears = experienceYears
	d.UpdatedAt = time.Now()
	return nil
}

func (d *Driver) SetPhotoURL(photoURL string) {
	d.PhotoURL = photoURL
	d.UpdatedAt = time.Now()
}
