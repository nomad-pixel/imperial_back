package entities

import "time"

type CarCategory struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
