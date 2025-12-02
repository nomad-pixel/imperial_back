package entities

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	IsVerified   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewUser(email, passwordHash string) (*User, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	if !emailRegex.MatchString(email) {
		return nil, errors.New("invalid email format")
	}

	if passwordHash == "" {
		return nil, errors.New("password hash cannot be empty")
	}

	if len(passwordHash) < 32 {
		return nil, errors.New("password hash is too short (must be bcrypt hash)")
	}

	now := time.Now()
	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		IsVerified:   false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (u *User) Validate() error {
	if u.ID <= 0 {
		return errors.New("invalid user ID")
	}

	if u.Email == "" {
		return errors.New("email cannot be empty")
	}

	if !emailRegex.MatchString(u.Email) {
		return errors.New("invalid email format")
	}

	if u.PasswordHash == "" {
		return errors.New("password hash cannot be empty")
	}

	return nil
}

func (u *User) SetEmail(email string) error {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return errors.New("email cannot be empty")
	}

	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}

	u.Email = email
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) MarkAsVerified() {
	u.IsVerified = true
	u.UpdatedAt = time.Now()
}
