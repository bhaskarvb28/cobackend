package academyCoach

import "time"

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

// ======================================================
// Core Coach Models
// ======================================================

type DisciplineSummary struct {
	ID int `json:"id"`

	Code string `json:"code"`

	DisplayName string `json:"display_name"`
}

type CoachSummary struct {
	UserID string `json:"user_id"`

	FullName string `json:"full_name"`

	CoachCode string `json:"coach_code"`

	Disciplines []DisciplineSummary `json:"disciplines"`

	AssignedPlayersCount int `json:"assigned_players_count"`
}

// ======================================================
// Paginated Coach Response
// ======================================================

type PaginatedAcademyCoachesResponse struct {
	Items []CoachSummary `json:"items"`

	Page int `json:"page"`

	Limit int `json:"limit"`

	Total int `json:"total"`

	TotalPages int `json:"total_pages"`

	HasNext bool `json:"has_next"`

	HasPrevious bool `json:"has_previous"`
}

// ======================================================
// Coach Profile
// ======================================================

type AcademyCoachProfileResponse struct {
	UserID string `json:"user_id"`

	CoachCode string `json:"coach_code"`

	FirstName string `json:"first_name"`

	LastName string `json:"last_name"`

	FullName string `json:"full_name"`

	Email string `json:"email"`

	Status string `json:"status"`

	JoinedAt time.Time `json:"joined_at"`

	Disciplines []DisciplineSummary `json:"disciplines"`

	AssignedPlayersCount int `json:"assigned_players_count"`
}

// ======================================================
// Query Models
// ======================================================

type GetAcademyCoachesQuery struct {
	Page int `json:"page"`

	Limit int `json:"limit"`

	Search string `json:"search"`

	DisciplineID int `json:"discipline_id"`

	SortBy string `json:"sort_by"`

	OrderBy string `json:"order_by"`
}

// ======================================================
// Assignment Models
// ======================================================

type AssignCoachInput struct {
	PlayerID string `json:"player_id"`

	CoachUserID string `json:"coach_user_id"`
}