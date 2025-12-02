package entities

import (
	"errors"
	"strings"
	"time"
)

type Lead struct {
	ID        int64     `json:"id"`
	FullName  string    `json:"full_name"`
	Phone     string    `json:"phone"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
}

func NewLead(fullName, phone string, startDate, endDate time.Time) (*Lead, error) {
	fullName = strings.TrimSpace(fullName)

	if fullName == "" {
		return nil, errors.New("full name cannot be empty")
	}

	if len(fullName) < 2 {
		return nil, errors.New("full name must be at least 2 characters")
	}

	if len(fullName) > 100 {
		return nil, errors.New("full name cannot exceed 100 characters")
	}

	phone = strings.TrimSpace(phone)
	if phone == "" {
		return nil, errors.New("phone cannot be empty")
	}

	if len(phone) > 32 {
		return nil, errors.New("phone cannot exceed 32 characters")
	}

	if startDate.IsZero() {
		return nil, errors.New("start date cannot be zero")
	}

	if endDate.IsZero() {
		return nil, errors.New("end date cannot be zero")
	}

	if endDate.Before(startDate) {
		return nil, errors.New("end date must be after start date")
	}

	return &Lead{
		FullName:  fullName,
		Phone:     phone,
		StartDate: startDate,
		EndDate:   endDate,
		CreatedAt: time.Now(),
	}, nil
}

func (l *Lead) Validate() error {
	if l.ID <= 0 {
		return errors.New("invalid lead ID")
	}

	fullName := strings.TrimSpace(l.FullName)
	if fullName == "" || len(fullName) < 2 || len(fullName) > 100 {
		return errors.New("full name must be between 2 and 100 characters")
	}

	phone := strings.TrimSpace(l.Phone)
	if phone == "" || len(phone) > 32 {
		return errors.New("phone is invalid")
	}

	if l.StartDate.IsZero() || l.EndDate.IsZero() {
		return errors.New("dates cannot be zero")
	}

	if l.EndDate.Before(l.StartDate) {
		return errors.New("end date must be after start date")
	}

	return nil
}
