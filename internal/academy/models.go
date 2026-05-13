package academy

type CreateAcademyInput struct {
	Name       string `json:"name"`
	DistrictID int    `json:"district_id"`
	Address    string `json:"address"`
}

type Academy struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	StateID    int    `json:"state_id"`
	DistrictID int    `json:"district_id"`
	Address    string `json:"address"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type GetAcademiesQuery struct {
	Page       int
	Limit      int
	Search     string
	StateID    int
	DistrictID int
	SortBy     string
	OrderBy    string
}

type PaginatedAcademies struct {
	Items       []Academy `json:"items"`
	Page        int       `json:"page"`
	Limit       int       `json:"limit"`
	Total       int       `json:"total"`
	TotalPages  int       `json:"total_pages"`
	HasNext     bool      `json:"has_next"`
	HasPrevious bool      `json:"has_previous"`
}