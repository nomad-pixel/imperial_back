package entities

import "time"

type CarMark struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
