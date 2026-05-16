package stateAdmin

type InviteStateAdminInput struct {
	Email           string `json:"email"`
	StateID int `json:"state_id"`
}

type GetStateAdminsQuery struct {
	Page 		int
	Limit 		int
	Search 		string
	StateID 	int	
	SortBy 		string // Column name: first_name, last_name, email, created_at
	OrderBy 	string // asc or desc
}

type StateAdmin struct {
	ID              string `json:"id"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	ContactNumber   string `json:"contact_number"`
	StateID         int    `json:"state_id"`
}

type PaginatedStateAdmins struct {
	Items       []StateAdmin 		 `json:"items"`
	Page        int                  `json:"page"`
	Limit       int                  `json:"limit"`
	Total       int                  `json:"total"`
	TotalPages  int                  `json:"total_pages"`
	HasNext     bool                 `json:"has_next"`
	HasPrevious bool                 `json:"has_previous"`
}

type UpdateStateInput struct {
	StateID int `json:"state_id"`
}

// type CreateStateAdminInput struct {
// 	FirstName     string `json:"first_name"`
// 	LastName      string `json:"last_name"`
// 	Email         string `json:"email"`
// 	Password      string `json:"password"`
// 	ContactNumber string `json:"contact_number"`
// 	AssignedState string    `json:"assigned_state"`
// }

