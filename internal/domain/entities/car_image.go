package entities

import "time"

type CarImage struct {
	ID        int64     `json:"id"`
	CarID     int64     `json:"car_id"`
	ImagePath string    `json:"image_path"`
	CreatedAt time.Time `json:"created_at"`
}
