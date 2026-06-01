package academy

import "time"

// CreateAcademyInput represents the request payload used to create a new academy.
type CreateAcademyInput struct {
	Name       string `json:"name"`        // Unique academy name
	Address    string `json:"address"`     // Full physical address
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

	Address  string `json:"address"`
	IsActive bool   `json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAcademiesQuery represents parsed academy listing filters.
type GetAcademiesQuery struct {
	Page       int    `json:"page"`         // Current page
	Limit      int    `json:"limit"`        // Pagination limit
	Search     string `json:"search"`       // Name/address search

	StateID    int    `json:"state_id"`     // Filter by state
	DistrictID int    `json:"district_id"`  // Filter by district
	PincodeID  int    `json:"pincode_id"`   // Filter by pincode

	SortBy     string `json:"sort_by"`      // Sort column
	OrderBy    string `json:"order_by"`     // ASC/DESC
}

// PaginatedAcademies contains a list of Academy objects along with standard metadata pagination metrics.
type PaginatedAcademies struct {
	Items       []AcademyResponse `json:"items"`        // Slice of retrieved academy records
	Page        int       `json:"page"`         // Current page number
	Limit       int       `json:"limit"`        // Limit setting applied
	Total       int       `json:"total"`        // Total count matching filters
	TotalPages  int       `json:"total_pages"`  // Total page count calculated dynamically
	HasNext     bool      `json:"has_next"`     // Boolean indicating if next page is available
	HasPrevious bool      `json:"has_previous"` // Boolean indicating if previous page is available
}

type CreateAcademyBuildingInput struct {
	BuildingName string `json:"building_name"`
}

type AcademyBuildingResponse struct {
	ID           int64  `json:"id"`
	AcademyID    string `json:"academy_id"`
	BuildingName string `json:"building_name"`
	IsActive     bool   `json:"is_active"`
}

type AddAcademyBuildingDisciplineInput struct {
	DisciplineID int `json:"discipline_id"`
}

type AcademyBuildingDisciplineResponse struct {
	AcademyBuildingID int `json:"academy_building_id"`
	DisciplineID      int `json:"discipline_id"`
}

type AddAcademyBuildingEventInput struct {
	ShootingEventID int `json:"shooting_event_id"`
}

type AcademyBuildingEventResponse struct {
	AcademyBuildingID int64 `json:"academy_building_id"`
	ShootingEventID   int   `json:"shooting_event_id"`
}

type BuildingDiscipline struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	DisplayName string `json:"display_name"`
}

type BuildingEvent struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	DisplayName string `json:"display_name"`
}

type AcademyBuilding struct {
	ID           int64                `json:"id"`
	AcademyID    string               `json:"academy_id"`
	BuildingName string               `json:"building_name"`
	IsActive     bool                 `json:"is_active"`
	Disciplines  []BuildingDiscipline `json:"disciplines"`
	Events       []BuildingEvent      `json:"events"`
}

type AddAcademyBuildingLaneInput struct {
	LaneName string `json:"lane_name" validate:"required,max=50"`
}

type AcademyBuildingLaneResponse struct {
	ID                 int64  `json:"id"`
	AcademyBuildingID  int64  `json:"academy_building_id"`
	LaneName           string `json:"lane_name"`
	IsUnderMaintenance bool   `json:"is_under_maintenance"`
	IsOccupied         bool   `json:"is_occupied"`
}

type AvailableLaneResponse struct {
	ID       int64  `json:"id"`
	LaneName string `json:"lane_name"`
}