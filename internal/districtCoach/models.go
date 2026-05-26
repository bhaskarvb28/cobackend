package districtCoach

type InviteDistrictCoachInput struct {
	Email 		string `json:"email"`
	StateID		int `json:"state_id"`
	DistrictID	int `json:"district_id"`
}

type CreateDistrictCoachInput struct {
	UserID     string
	DistrictID int32
}

type GetDistrictCoachesQuery struct {
	Page       int
	Limit      int
	Search     string
	StateID	   int
	DistrictID int
	SortBy     string
	OrderBy    string
}

type DistrictCoach struct {
	ID                       string `json:"id"`
	FirstName                string `json:"first_name"`
	LastName                 string `json:"last_name"`
	Email                    string `json:"email"`
	ContactNumber            string `json:"contact_number"`
	StateID					 int 	`json:"state_id"`
	DistrictID               int    `json:"district_id"`
	CoachCode                string `json:"coach_code"`
	CoachCertificationProof  string `json:"coach_certification_proof"`
	DPDPConsent              bool   `json:"dpdp_consent"`
}

type PaginatedDistrictCoaches struct {
	Items       []DistrictCoach `json:"items"`
	Page        int             `json:"page"`
	Limit       int             `json:"limit"`
	Total       int             `json:"total"`
	TotalPages  int             `json:"total_pages"`
	HasNext     bool            `json:"has_next"`
	HasPrevious bool            `json:"has_previous"`
}



type UpdateDistrictCoachInput struct {
	DistrictID                 *int    `json:"district_id"`
}