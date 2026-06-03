package academyAdmin

type InviteAcademyAdminInput struct {
	Email string `json:"email"`

	StateID int `json:"state_id"`

	DistrictID int `json:"district_id"`

	AcademyID string `json:"academy_id"`
}

type CreateAcademyAdminInput struct {
	UserID string

	AcademyID string
}

type AcademyAdmin struct {
	ID string `json:"id"`

	FirstName string `json:"first_name"`

	LastName string `json:"last_name"`

	Email string `json:"email"`

	ContactNumber string `json:"contact_number"`

	StateID int `json:"state_id"`

	DistrictID int `json:"district_id"`

	AcademyID string `json:"academy_id"`

	GSTIN *string `json:"gstin"`

	RegistrationProof *string `json:"registration_proof"`

	DPDPConsent bool `json:"dpdp_consent"`

	CreatedAt string `json:"created_at"`
}

type GetAcademyAdminsQuery struct {
	Page int

	Limit int

	Search string

	StateID int

	DistrictID int

	AcademyID string

	SortBy string

	OrderBy string
}

type PaginatedAcademyAdmins struct {
	Items []AcademyAdmin `json:"items"`

	Page int `json:"page"`

	Limit int `json:"limit"`

	Total int `json:"total"`

	TotalPages int `json:"total_pages"`

	HasNext bool `json:"has_next"`

	HasPrevious bool `json:"has_previous"`
}

type UpdateAcademyAdminInput struct {
	AcademyID *string `json:"academy_id"`
}