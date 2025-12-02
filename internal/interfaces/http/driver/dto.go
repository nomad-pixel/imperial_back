package driver

type CreateDriverRequest struct {
	FullName        string `json:"full_name" binding:"required"`
	About           string `json:"about" binding:"required"`
	ExperienceYears string `json:"experience_years" binding:"required"`
}

type UpdateDriverRequest struct {
	FullName        string `json:"full_name" binding:"required"`
	About           string `json:"about" binding:"required"`
	ExperienceYears string `json:"experience_years" binding:"required"`
}

type ListDriversResponse struct {
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}
