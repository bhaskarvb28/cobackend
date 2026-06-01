package academy

import (
	"context"
	"strings"

	"cobackend/internal/districtAdmin"
	"cobackend/internal/shared"
	"cobackend/internal/academyAdmin"
)

func CreateAcademyService(
	ctx context.Context,
	userID string,
	input CreateAcademyInput,
) (*AcademyResponse, error) {

	// ----------------------------------------------------------
	// Normalize Input
	// ----------------------------------------------------------

	input.Name = strings.TrimSpace(input.Name)
	input.Address = strings.TrimSpace(input.Address)

	// ----------------------------------------------------------
	// Validate Input
	// ----------------------------------------------------------

	if input.Name == "" {
		return nil, shared.ErrAcademyNameRequired
	}

	if input.Address == "" {
		return nil, shared.ErrAddressRequired
	}

	// ----------------------------------------------------------
	// Get District Admin Region
	// ----------------------------------------------------------

	districtAdminRegion, err := districtAdmin.GetDistrictAdminRegion(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	// ----------------------------------------------------------
	// Create Academy
	// ----------------------------------------------------------

	academy, err := CreateAcademyRepository(
		ctx,
		districtAdminRegion.DistrictID,
		input,
	)

	if err != nil {
		return nil, err
	}

	return academy, nil
}

// GetAcademiesService retrieves a paginated list of academies
// based on state, district, search, and sorting criteria.
func GetAcademiesService(
	ctx context.Context,
	query GetAcademiesQuery,
) (PaginatedAcademies, error) {

	// ----------------------------------------------------------
	// Fetch Academies
	// ----------------------------------------------------------
	// Query repository layer directly to fetch rows and pagination total counts
	return GetAcademiesRepository(
		ctx,
		query,
	)
}

// GetDistrictAdminAcademiesService retrieves
// academies belonging to the authenticated
// district admin.
func GetDistrictAdminAcademiesService(
	ctx context.Context,
	userID string,
	query GetAcademiesQuery,
) (PaginatedAcademies, error) {

	// ----------------------------------------------------------
	// Fetch District Academies
	// ----------------------------------------------------------

	return GetDistrictAdminAcademiesRepository(
		ctx,
		userID,
		query,
	)
}



func CreateAcademyBuildingService(
	ctx context.Context,
	userID string,
	input CreateAcademyBuildingInput,
) (*AcademyBuildingResponse, error) {

	// ----------------------------------------------------------
	// Normalize Input
	// ----------------------------------------------------------

	input.BuildingName = strings.TrimSpace(
		input.BuildingName,
	)

	// ----------------------------------------------------------
	// Validate Input
	// ----------------------------------------------------------

	if input.BuildingName == "" {
		return nil, shared.ErrAcademyBuildingNameRequired
	}

	// ----------------------------------------------------------
	// Get Academy Admin Academy
	// ----------------------------------------------------------

	academyID, err := academyAdmin.GetAcademyAdminAcademyID(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	// ----------------------------------------------------------
	// Create Building
	// ----------------------------------------------------------

	academyBuilding, err := CreateAcademyBuildingRepository(
		ctx,
		academyID,
		input,
	)

	if err != nil {
		return nil, err
	}

	return academyBuilding, nil
}

func AddAcademyBuildingDisciplineService(
	ctx context.Context,
	userID string,
	buildingID int64,
	input AddAcademyBuildingDisciplineInput,
) (*AcademyBuildingDisciplineResponse, error) {

	// ----------------------------------------------------------
	// Validate Input
	// ----------------------------------------------------------

	if input.DisciplineID <= 0 {
		return nil, shared.ErrInvalidDisciplineID
	}

	// ----------------------------------------------------------
	// Get Academy Admin Academy
	// ----------------------------------------------------------

	academyID, err := academyAdmin.GetAcademyAdminAcademyID(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	// ----------------------------------------------------------
	// Validate Building Ownership
	// ----------------------------------------------------------

	isOwned, err := CheckAcademyBuildingOwnershipRepository(
		ctx,
		buildingID,
		academyID,
	)

	if err != nil {
		return nil, err
	}

	if !isOwned {
		return nil, shared.ErrUnauthorizedBuildingAccess
	}

	// ----------------------------------------------------------
	// Add Discipline
	// ----------------------------------------------------------

	response, err := AddAcademyBuildingDisciplineRepository(
		ctx,
		buildingID,
		input.DisciplineID,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// ============================================================================
// service.go
// ============================================================================

func AddAcademyBuildingEventService(
	ctx context.Context,
	userID string,
	buildingID int64,
	input AddAcademyBuildingEventInput,
) (*AcademyBuildingEventResponse, error) {

	// ----------------------------------------------------------
	// Validate Input
	// ----------------------------------------------------------

	if input.ShootingEventID <= 0 {
		return nil, shared.ErrInvalidShootingEventID
	}

	// ----------------------------------------------------------
	// Get Academy Admin Academy
	// ----------------------------------------------------------

	academyID, err := academyAdmin.GetAcademyAdminAcademyID(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	// ----------------------------------------------------------
	// Validate Building Ownership
	// ----------------------------------------------------------

	isOwned, err := CheckAcademyBuildingOwnershipRepository(
		ctx,
		buildingID,
		academyID,
	)

	if err != nil {
		return nil, err
	}

	if !isOwned {
		return nil, shared.ErrUnauthorizedBuildingAccess
	}

	// ----------------------------------------------------------
	// Add Building Event
	// ----------------------------------------------------------

	response, err := AddAcademyBuildingEventRepository(
		ctx,
		buildingID,
		input.ShootingEventID,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}

// ============================================================================
// service.go
// ============================================================================

func GetAcademyBuildingsService(
	ctx context.Context,
	userID string,
) ([]AcademyBuilding, error) {

	// ----------------------------------------------------------
	// Get Academy Admin Academy
	// ----------------------------------------------------------

	academyID, err := academyAdmin.GetAcademyAdminAcademyID(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	// ----------------------------------------------------------
	// Get Buildings
	// ----------------------------------------------------------

	buildings, err := GetAcademyBuildingsRepository(
		ctx,
		academyID,
	)

	if err != nil {
		return nil, err
	}

	return buildings, nil
}

func AddAcademyBuildingLaneService(
	ctx context.Context,
	userID string,
	buildingID int64,
	input AddAcademyBuildingLaneInput,
) (*AcademyBuildingLaneResponse, error) {

	// ----------------------------------------------------------
	// Normalize Input
	// ----------------------------------------------------------

	input.LaneName = strings.TrimSpace(
		input.LaneName,
	)

	// ----------------------------------------------------------
	// Validate Input
	// ----------------------------------------------------------

	if input.LaneName == "" {
		return nil, shared.ErrLaneNameRequired
	}

	// ----------------------------------------------------------
	// Get Academy Admin Academy
	// ----------------------------------------------------------

	academyID, err := academyAdmin.GetAcademyAdminAcademyID(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	// ----------------------------------------------------------
	// Validate Building Ownership
	// ----------------------------------------------------------

	isOwned, err := CheckAcademyBuildingOwnershipRepository(
		ctx,
		buildingID,
		academyID,
	)

	if err != nil {
		return nil, err
	}

	if !isOwned {
		return nil, shared.ErrUnauthorizedBuildingAccess
	}

	// ----------------------------------------------------------
	// Create Lane
	// ----------------------------------------------------------

	response, err := CreateAcademyBuildingLaneRepository(
		ctx,
		buildingID,
		input,
	)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetAvailableLanesService(
	ctx context.Context,
	buildingID int64,
) ([]AvailableLaneResponse, error) {

	// ----------------------------------------------------------
	// Validate Input
	// ----------------------------------------------------------

	if buildingID <= 0 {
		return nil, shared.ErrInvalidBuildingID
	}

	// ----------------------------------------------------------
	// Get Available Lanes
	// ----------------------------------------------------------

	lanes, err := GetAvailableLanesRepository(
		ctx,
		buildingID,
	)

	if err != nil {
		return nil, err
	}

	return lanes, nil
}

