package academyAdmin

type InviteAcademyAdminInput struct {
	Email 		string `json:"email"`
	StateID		int `json:"state_id"`
	DistrictID	int `json:"district_id"`
	AcademyID   int `json:"academy_id"`
}

