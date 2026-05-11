package districtadmin

type InviteDistrictAdminInput struct {
	Email 			string `json:"email"`
	AssignedDistrictID string `json:"assigned_district_id"`
}

// type CreateDistrictAdminInput struct {
// 	FirstName  string `json:"first_name"`
// 	LastName   string `json:"last_name"`
// 	Email      string `json:"email"`
// 	Password   string `json:"password"`
// 	ContactNumber string `json:"contact_number"`
// 	StateID    int    `json:"state_id"`
// 	DistrictID int    `json:"district_id"`
// }

// type UpdateDistrictAdminInput struct {
// 	StateID    int    `json:"state_id"`
// 	DistrictID int    `json:"district_id"`
// 	ApprovalStatus string `json:"approval_status"`
// 	ApprovalNotes  string `json:"approval_notes"`
// }

// type DistrictAdmin struct {
// 	ID             string `json:"id"`
// 	FirstName      string `json:"first_name"`
// 	LastName       string `json:"last_name"`
// 	Email          string `json:"email"`
// 	ContactNumber  string `json:"contact_number"`
// 	StateID        int    `json:"state_id"`
// 	DistrictID     int    `json:"district_id"`
// 	DPDPConsent    bool   `json:"dpdp_consent"`
// 	ApprovalStatus string `json:"approval_status"`
// 	ApprovalNotes  *string `json:"approval_notes"`
// }

// type GetDistrictAdminsQuery struct {
// 	Page     int
// 	Limit    int
// 	Search   string
// 	StateID  int
// 	DistrictID int
// 	Status   string
// }
