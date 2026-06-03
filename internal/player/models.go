package player

import (
	"cobackend/internal/profile"
	"time"
)

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

// ==========================================================
// Academy Players
// ==========================================================

type GetAcademyPlayersQuery struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`

	Search string `json:"search"`

	DisciplineID int `json:"discipline_id"`

	CoachAssigned *bool `json:"coach_assigned,omitempty"`

	Status string `json:"status"`

	SortBy  string `json:"sort_by"`
	OrderBy string `json:"order_by"`
}

type PlayerListItemResponse struct {
	ID string `json:"id"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Email string `json:"email"`

	Status string `json:"status"`

	ProfileCompleted bool `json:"profile_completed"`

	Gender string `json:"gender"`

	PrimaryDiscipline *profile.Discipline `json:"primary_discipline,omitempty"`

	CurrentCoach *profile.CoachSummary `json:"current_coach,omitempty"`

	JoinedAt time.Time `json:"joined_at"`
}

type PaginatedPlayers struct {
	Items []PlayerListItemResponse `json:"items"`

	Page        int `json:"page"`
	Limit       int `json:"limit"`
	Total       int `json:"total"`
	TotalPages  int `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}


type ShootingEventResponse struct {
	ID          int16   `json:"id"`
	Code        string  `json:"code"`
	DisplayName string  `json:"display_name"`
	Distance    *int16  `json:"distance,omitempty"`
}

type CompatibleBuildingResponse struct {
	ID           int64     `json:"id"`
	BuildingName string    `json:"building_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

