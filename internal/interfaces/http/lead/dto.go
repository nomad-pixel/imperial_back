package lead

import "time"

type CreateLeadRequest struct {
	FullName  string    `json:"full_name" binding:"required"`
	Phone     string    `json:"phone" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

type ListLeadsResponse struct {
	Total int64 `json:"total"`
	Data  interface{} `json:"data"`
}
