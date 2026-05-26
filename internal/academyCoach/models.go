package academyCoach

type InviteAcademyCoachInput struct {
	Email string `json:"email"`

	StateID int `json:"state_id"`

	DistrictID int `json:"district_id"`

	AcademyID int `json:"academy_id"`

	DisciplinesSpecialized []int `json:"disciplines_specialized"`
}

type CreateAcademyCoachInput struct {
	UserID    string
	AcademyID string
}