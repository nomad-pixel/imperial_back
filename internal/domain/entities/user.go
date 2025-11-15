package entities

import "time"

// User представляет доменную модель пользователя
type User struct {
	ID           int64
	Email        string
	PasswordHash string
	VerifiedAt   bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
