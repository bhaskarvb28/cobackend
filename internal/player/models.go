package player

type InvitePlayerInput struct {
	Email           string `json:"email"`
	StateID         int    `json:"state_id"`
	DistrictID      int    `json:"district_id"`
	AcademyID       int    `json:"academy_id"`
	AcademyCoachID  string `json:"academy_coach_id"`
}

type CreatePlayerInput struct {
	UserID       string
	AcademyID    string
	RegisteredBy string
}