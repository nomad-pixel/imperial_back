package car

import "github.com/nomad-pixel/imperial/internal/domain/entities"

// CreateCarTagRequest represents the request to create a car tag
type CreateCarTagRequest struct {
	Name string `json:"name" binding:"required" example:"Sedan"`
}

// UpdateCarTagRequest represents the request to update a car tag
type UpdateCarTagRequest struct {
	Name string `json:"name" binding:"required" example:"Sedan"`
}

// CarTagResponse represents a car tag response
type CarTagResponse = entities.CarTag

// MessageResponse represents a generic message response
type MessageResponse struct {
	Message string `json:"message"`
}

// ListCarTagsResponse represents the response for listing car tags
type ListCarTagsResponse struct {
	Total int64              `json:"total"`
	Data  []*entities.CarTag `json:"data"`
}
