package session
// model.go

import (
	"time"
)

type StartPracticeSessionInput struct {
	AcademyBuildingLaneID int64 `json:"academy_building_lane_id"`
	ShootingEventID       int16 `json:"shooting_event_id"`
}

type PracticeSessionResponse struct {
	ID                     int64     `json:"id"`
	PlayerUserID           string    `json:"player_user_id"`
	AcademyBuildingLaneID  int64     `json:"academy_building_lane_id"`
	ShootingEventID        int16     `json:"shooting_event_id"`
	Status                 string    `json:"status"`
	TotalScore             string    `json:"total_score"`
	TotalShotCount         int32     `json:"total_shot_count"`
	StartedAt              time.Time `json:"started_at"`
	EndedAt                *time.Time `json:"ended_at,omitempty"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}