package academy

import (
	"time"
)

// ==========================================================
// Academy
// ==========================================================

// CreateAcademyInput represents the request payload
// used to create a new academy.
type CreateAcademyInput struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type AcademyResponse struct {
	ID string `json:"id"`

	Name string `json:"name"`

	StateID   int    `json:"state_id"`
	StateName string `json:"state_name"`

	DistrictID   int    `json:"district_id"`
	DistrictName string `json:"district_name"`

	PincodeID int    `json:"pincode_id"`
	Pincode   string `json:"pincode"`

	Address string `json:"address"`

	IsActive bool `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ==========================================================
// Academy Listing
// ==========================================================

// GetAcademiesQuery represents parsed academy filters.
type GetAcademiesQuery struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`

	StateID    int `json:"state_id"`
	DistrictID int `json:"district_id"`
	PincodeID  int `json:"pincode_id"`

	SortBy  string `json:"sort_by"`
	OrderBy string `json:"order_by"`
}

type PaginatedAcademies struct {
	Items []AcademyResponse `json:"items"`
	Page        int `json:"page"`
	Limit       int `json:"limit"`
	Total       int `json:"total"`
	TotalPages  int `json:"total_pages"`
	HasNext     bool `json:"has_next"`
	HasPrevious bool `json:"has_previous"`
}

// ==========================================================
// Academy Buildings
// ==========================================================

type CreateAcademyBuildingInput struct {
	BuildingName string `json:"building_name"`
}

type AcademyBuildingResponse struct {
	ID int64 `json:"id"`

	AcademyID string `json:"academy_id"`

	BuildingName string `json:"building_name"`

	IsActive bool `json:"is_active"`
}

type AddAcademyBuildingDisciplineInput struct {
	DisciplineID int `json:"discipline_id"`
}

type AcademyBuildingDisciplineResponse struct {
	AcademyBuildingID int `json:"academy_building_id"`

	DisciplineID int `json:"discipline_id"`
}

type AddAcademyBuildingEventInput struct {
	ShootingEventID int `json:"shooting_event_id"`
}

type AcademyBuildingEventResponse struct {
	AcademyBuildingID int64 `json:"academy_building_id"`

	ShootingEventID int `json:"shooting_event_id"`
}

type BuildingDiscipline struct {
	ID int `json:"id"`

	Code string `json:"code"`

	DisplayName string `json:"display_name"`
}

type BuildingEvent struct {
	ID int `json:"id"`

	Code string `json:"code"`

	DisplayName string `json:"display_name"`
}

type AcademyBuilding struct {
	ID int64 `json:"id"`

	AcademyID string `json:"academy_id"`

	BuildingName string `json:"building_name"`

	IsActive bool `json:"is_active"`

	Disciplines []BuildingDiscipline `json:"disciplines"`

	Events []BuildingEvent `json:"events"`
}

// ==========================================================
// Academy Lanes
// ==========================================================

type AddAcademyBuildingLaneInput struct {
	LaneName string `json:"lane_name" validate:"required,max=50"`
}

type AcademyBuildingLaneResponse struct {
	ID int64 `json:"id"`

	AcademyBuildingID int64 `json:"academy_building_id"`

	LaneName string `json:"lane_name"`

	IsUnderMaintenance bool `json:"is_under_maintenance"`
}

type AvailableLaneResponse struct {
	ID int64 `json:"id"`

	LaneName string `json:"lane_name"`
}

type AcademyBuildingDetailsResponse struct {
	ID int64 `json:"id"`

	AcademyID string `json:"academy_id"`

	BuildingName string `json:"building_name"`

	IsActive bool `json:"is_active"`

	Disciplines []BuildingDiscipline `json:"disciplines"`

	Events []BuildingEvent `json:"events"`

	Lanes []AcademyBuildingLaneResponse `json:"lanes"`
}

type UpdateAcademyBuildingInput struct {
	BuildingName string `json:"building_name"`

	IsActive *bool `json:"is_active,omitempty"`
}

type UpdateAcademyBuildingLaneInput struct {
	LaneName string `json:"lane_name"`

	IsUnderMaintenance *bool `json:"is_under_maintenance,omitempty"`

}

type EventResponse struct {
	ID int64 `json:"id"`

	Code string `json:"code"`

	DisplayName string `json:"display_name"`
}