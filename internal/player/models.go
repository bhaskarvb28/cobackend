package player

import (
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

