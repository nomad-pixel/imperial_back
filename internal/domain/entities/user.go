package entities

import "time"

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	VerifiedAt   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
