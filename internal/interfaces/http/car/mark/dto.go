package car

import "github.com/nomad-pixel/imperial/internal/domain/entities"

// CreateCarMarkRequest represents the request to create a car mark
type CreateCarMarkRequest struct {
	Name string `json:"name" binding:"required" example:"Toyota"`
}

// UpdateCarMarkRequest represents the request to update a car mark
type UpdateCarMarkRequest struct {
	Name string `json:"name" binding:"required" example:"Toyota"`
}

// CarMarkResponse represents a car mark response
type CarMarkResponse = entities.CarMark

// MessageResponse represents a generic message response
type MessageResponse struct {
	Message string `json:"message"`
}

// ListCarMarksResponse represents the response for listing car marks
type ListCarMarksResponse struct {
	Total int64               `json:"total"`
	Data  []*entities.CarMark `json:"data"`
}
