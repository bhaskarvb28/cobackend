package academy

import "time"

type CreateAcademyInput struct {
	Name       string `json:"name"`
	DistrictID int    `json:"district_id"`
	Address    string `json:"address"`
}

type AcademyResponse struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	DistrictID int       `json:"district_id"`
	Address    string    `json:"address"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type GetAcademiesQuery struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Search     string `json:"search"`
	StateID    int    `json:"state_id"`
	DistrictID int    `json:"district_id"`
	SortBy     string `json:"sort_by"`
	OrderBy    string `json:"order_by"`
}

type PaginatedAcademies struct {
	Items       []AcademyResponse `json:"items"`
	Page        int               `json:"page"`
	Limit       int               `json:"limit"`
	Total       int               `json:"total"`
	TotalPages  int               `json:"total_pages"`
	HasNext     bool              `json:"has_next"`
	HasPrevious bool              `json:"has_previous"`
}