package entities

import "time"

type CarTag struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
