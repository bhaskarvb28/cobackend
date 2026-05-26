package districtAdmin

type InviteDistrictAdminInput struct {
	Email 		string `json:"email"`
	StateID		int `json:"state_id"`
	DistrictID	int `json:"district_id"`
}

type CreateDistrictAdminInput struct {
	UserID       string
	DistrictID   int32
}

type GetDistrictAdminsQuery struct {
	Page        int
	Limit       int
	Search      string
	StateID     int
	DistrictID  int
	SortBy      string
	OrderBy     string
}

type DistrictAdmin struct {
	ID             string `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Email          string `json:"email"`
	ContactNumber  string `json:"contact_number"`
	StateID        int    `json:"state_id"`
	DistrictID     int    `json:"district_id"`
	DPDPConsent    bool   `json:"dpdp_consent"`
}

type PaginatedDistrictAdmins struct {
	Items        []DistrictAdmin `json:"items"`
	Page         int             `json:"page"`
	Limit        int             `json:"limit"`
	Total        int             `json:"total"`
	TotalPages   int             `json:"total_pages"`
	HasNext      bool            `json:"has_next"`
	HasPrevious  bool            `json:"has_previous"`
}

type UpdateDistrictAdminInput struct {
	DistrictID     *int    `json:"district_id"`
}

