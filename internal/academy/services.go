package academy

import (
	"context"
	"errors"
	"strings"

	"cobackend/internal/academyAdmin"
	"cobackend/internal/academyCoach"
	"cobackend/internal/districtAdmin"
	"cobackend/internal/player"
	"cobackend/internal/profile"
	"cobackend/internal/shared"
)

func CreateAcademyService(
	ctx context.Context,
	userID string,
	input CreateAcademyInput,
) (error) {

	// ----------------------------------------------------------
	// Normalize Input
	// ----------------------------------------------------------

	input.Name = strings.TrimSpace(input.Name)
	input.Address = strings.TrimSpace(input.Address)

	// ----------------------------------------------------------
	// Validate Input
	// ----------------------------------------------------------

	if input.Name == "" {
		return shared.ErrAcademyNameRequired
	}

	if input.Address == "" {
		return shared.ErrAddressRequired
	}

	// ----------------------------------------------------------
	// Get District Admin Region
	// ----------------------------------------------------------

	districtAdminRegion, err := districtAdmin.GetDistrictAdminRegion(
		ctx,
		userID,
	)

	if err != nil {
		return err
	}

	// ----------------------------------------------------------
	// Create Academy
	// ----------------------------------------------------------

	err = CreateAcademyRepository(
		ctx,
		districtAdminRegion.DistrictID,
		input,
	)

	if err != nil {
		return err
	}

	return nil
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

func GetAcademyPlayersService(
	ctx context.Context,
	userID string,
	query player.GetAcademyPlayersQuery,
) (player.PaginatedPlayers, error) {

	// ----------------------------------------------------------
	// Normalize Input
	// ----------------------------------------------------------

	query.Search = strings.TrimSpace(
		query.Search,
	)

	query.Status = strings.TrimSpace(
		query.Status,
	)

	// ----------------------------------------------------------
	// Get Academy Admin
	// ----------------------------------------------------------

	academyAdminProfile, err :=
		academyAdmin.GetAcademyAdminByIDRepository(
			ctx,
			userID,
		)

	if err != nil {

		return player.PaginatedPlayers{},
			err
	}

	// ----------------------------------------------------------
	// Fetch Players
	// ----------------------------------------------------------

	return player.GetAcademyPlayersRepository(
		ctx,
		academyAdminProfile.AcademyID,
		query,
	)
}

func GetAcademyPlayerService(
	ctx context.Context,
	userID string,
	playerID string,
) (profile.PlayerProfileResponse, error) {

	// ----------------------------------------------------------
	// Normalize Input
	// ----------------------------------------------------------

	playerID = strings.TrimSpace(
		playerID,
	)

	// ----------------------------------------------------------
	// Validate Input
	// ----------------------------------------------------------

	if playerID == "" {

		return profile.PlayerProfileResponse{},
			shared.ErrPlayerIDRequired
	}

	// ----------------------------------------------------------
	// Get Academy Admin
	// ----------------------------------------------------------

	academyAdminProfile, err :=
		academyAdmin.GetAcademyAdminByIDRepository(
			ctx,
			userID,
		)

	if err != nil {

		return profile.PlayerProfileResponse{},
			err
	}

	// ----------------------------------------------------------
	// Fetch Player
	// ----------------------------------------------------------

	playerProfile, err :=
		player.GetAcademyPlayerRepository(
			ctx,
			academyAdminProfile.AcademyID,
			playerID,
		)

	if err != nil {

		return profile.PlayerProfileResponse{},
			err
	}

	return playerProfile, nil
}

func GetAcademyCoachesService(
	ctx context.Context,
	authUserID string,
	query academyCoach.GetAcademyCoachesQuery,
) (
	academyCoach.PaginatedAcademyCoachesResponse,
	error,
) {

	// ----------------------------------------------------------
	// Get Academy ID From Admin
	// ----------------------------------------------------------

	academyID, err := academyAdmin.GetAcademyAdminAcademyID(
		ctx,
		authUserID,
	)

	if err != nil {

		return academyCoach.PaginatedAcademyCoachesResponse{},
			err
	}

	// ----------------------------------------------------------
	// Repository
	// ----------------------------------------------------------

	result, err := academyCoach.GetAcademyCoaches(
		ctx,
		academyID,
		query,
	)

	if err != nil {

		return academyCoach.PaginatedAcademyCoachesResponse{},
			err
	}

	return result, nil
}

func GetAcademyCoachService(
	ctx context.Context,
	authUserID string,
	coachID string,
) (
	academyCoach.AcademyCoachProfileResponse,
	error,
) {

	// ----------------------------------------------------------
	// Get Academy ID
	// ----------------------------------------------------------

	academyID, err := academyAdmin.GetAcademyAdminAcademyID(
		ctx,
		authUserID,
	)

	if err != nil {

		return academyCoach.AcademyCoachProfileResponse{},
			err
	}

	// ----------------------------------------------------------
	// Repository
	// ----------------------------------------------------------

	result, err := academyCoach.GetAcademyCoach(
		ctx,
		academyID,
		coachID,
	)

	if err != nil {

		return academyCoach.AcademyCoachProfileResponse{},
			err
	}

	return result, nil
}

func AssignCoachService(
	ctx context.Context,
	authUserID string,
	playerID string,
	coachUserID string,
) error {

	// ----------------------------------------------------------
	// Get Academy ID
	// ----------------------------------------------------------

	academyID, err := academyAdmin.GetAcademyAdminAcademyID(
		ctx,
		authUserID,
	)

	if err != nil {
		return err
	}

	// ----------------------------------------------------------
	// Validate Player Belongs To Academy
	// ----------------------------------------------------------

	playerExists, err := player.ValidateAcademyPlayer(
		ctx,
		academyID,
		playerID,
	)

	if err != nil {
		return err
	}

	if !playerExists {

		return errors.New(
			"player not found in academy",
		)
	}

	// ----------------------------------------------------------
	// Validate Coach Belongs To Academy
	// ----------------------------------------------------------

	coachExists, err := academyCoach.ValidateAcademyCoach(
		ctx,
		academyID,
		coachUserID,
	)

	if err != nil {
		return err
	}

	if !coachExists {

		return errors.New(
			"coach not found in academy",
		)
	}

	// ----------------------------------------------------------
	// Check Existing Assignment
	// ----------------------------------------------------------

	alreadyAssigned, err := player.CheckPlayerCoachAssignment(
		ctx,
		playerID,
		coachUserID,
	)

	if err != nil {
		return err
	}

	if alreadyAssigned {

		return errors.New(
			"coach already assigned to player",
		)
	}

	// ----------------------------------------------------------
	// Assign Coach
	// ----------------------------------------------------------

	err = player.AssignCoach(
		ctx,
		playerID,
		coachUserID,
	)

	if err != nil {
		return err
	}

	return nil
}

func RemoveCoachService(
	ctx context.Context,
	authUserID string,
	playerID string,
) error {

	// ----------------------------------------------------------
	// Get Academy ID
	// ----------------------------------------------------------

	academyID, err := academyAdmin.GetAcademyAdminAcademyID(
		ctx,
		authUserID,
	)

	if err != nil {
		return err
	}

	// ----------------------------------------------------------
	// Validate Player
	// ----------------------------------------------------------

	playerExists, err := player.ValidateAcademyPlayer(
		ctx,
		academyID,
		playerID,
	)

	if err != nil {
		return err
	}

	if !playerExists {

		return errors.New(
			"player not found in academy",
		)
	}

	// ----------------------------------------------------------
	// Remove Coach
	// ----------------------------------------------------------

	err = player.RemoveCoach(
		ctx,
		playerID,
	)

	if err != nil {
		return err
	}

	return nil
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

func GetAcademyBuildingService(
	ctx context.Context,
	userID string,
	buildingID int64,
) (*AcademyBuildingDetailsResponse, error) {

	if buildingID <= 0 {
		return nil, shared.ErrInvalidBuildingID
	}

	academyID, err := academyAdmin.
		GetAcademyAdminAcademyID(
			ctx,
			userID,
		)

	if err != nil {
		return nil, err
	}

	isOwned, err :=
		CheckAcademyBuildingOwnershipRepository(
			ctx,
			buildingID,
			academyID,
		)

	if err != nil {
		return nil, err
	}

	if !isOwned {
		return nil,
			shared.ErrUnauthorizedBuildingAccess
	}

	return GetAcademyBuildingRepository(
		ctx,
		buildingID,
	)
}

func UpdateAcademyBuildingService(
	ctx context.Context,
	userID string,
	buildingID int64,
	input UpdateAcademyBuildingInput,
) error {

	input.BuildingName =
		strings.TrimSpace(
			input.BuildingName,
		)

	if input.BuildingName == "" {
		return shared.ErrAcademyBuildingNameRequired
	}

	academyID, err := academyAdmin.
		GetAcademyAdminAcademyID(
			ctx,
			userID,
		)

	if err != nil {
		return err
	}

	isOwned, err :=
		CheckAcademyBuildingOwnershipRepository(
			ctx,
			buildingID,
			academyID,
		)

	if err != nil {
		return err
	}

	if !isOwned {
		return shared.ErrUnauthorizedBuildingAccess
	}

	return UpdateAcademyBuildingRepository(
		ctx,
		buildingID,
		input,
	)
}

func RemoveAcademyBuildingDisciplineService(
	ctx context.Context,
	userID string,
	buildingID int64,
	disciplineID int,
) error {

	academyID, err := academyAdmin.
		GetAcademyAdminAcademyID(
			ctx,
			userID,
		)

	if err != nil {
		return err
	}

	isOwned, err :=
		CheckAcademyBuildingOwnershipRepository(
			ctx,
			buildingID,
			academyID,
		)

	if err != nil {
		return err
	}

	if !isOwned {
		return shared.ErrUnauthorizedBuildingAccess
	}

	return RemoveAcademyBuildingDisciplineRepository(
		ctx,
		buildingID,
		disciplineID,
	)
}

func RemoveAcademyBuildingEventService(
	ctx context.Context,
	userID string,
	buildingID int64,
	eventID int,
) error {

	academyID, err := academyAdmin.
		GetAcademyAdminAcademyID(
			ctx,
			userID,
		)

	if err != nil {
		return err
	}

	isOwned, err :=
		CheckAcademyBuildingOwnershipRepository(
			ctx,
			buildingID,
			academyID,
		)

	if err != nil {
		return err
	}

	if !isOwned {
		return shared.ErrUnauthorizedBuildingAccess
	}

	return RemoveAcademyBuildingEventRepository(
		ctx,
		buildingID,
		eventID,
	)
}

func GetAcademyBuildingLanesService(
	ctx context.Context,
	userID string,
	buildingID int64,
) ([]AcademyBuildingLaneResponse, error) {

	academyID, err := academyAdmin.
		GetAcademyAdminAcademyID(
			ctx,
			userID,
		)

	if err != nil {
		return nil, err
	}

	isOwned, err :=
		CheckAcademyBuildingOwnershipRepository(
			ctx,
			buildingID,
			academyID,
		)

	if err != nil {
		return nil, err
	}

	if !isOwned {
		return nil,
			shared.ErrUnauthorizedBuildingAccess
	}

	return GetAcademyBuildingLanesRepository(
		ctx,
		buildingID,
	)
}

func UpdateAcademyBuildingLaneService(
	ctx context.Context,
	userID string,
	laneID int64,
	input UpdateAcademyBuildingLaneInput,
) error {

	input.LaneName =
		strings.TrimSpace(
			input.LaneName,
		)

	if input.LaneName == "" {
		return shared.ErrLaneNameRequired
	}

	academyID, err := academyAdmin.
		GetAcademyAdminAcademyID(
			ctx,
			userID,
		)

	if err != nil {
		return err
	}

	isOwned, err :=
		CheckLaneOwnershipRepository(
			ctx,
			laneID,
			academyID,
		)

	if err != nil {
		return err
	}

	if !isOwned {
		return shared.ErrUnauthorizedBuildingAccess
	}

	return UpdateAcademyBuildingLaneRepository(
		ctx,
		laneID,
		input,
	)
}

func DeleteAcademyBuildingLaneService(
	ctx context.Context,
	userID string,
	laneID int64,
) error {

	academyID, err := academyAdmin.
		GetAcademyAdminAcademyID(
			ctx,
			userID,
		)

	if err != nil {
		return err
	}

	isOwned, err :=
		CheckLaneOwnershipRepository(
			ctx,
			laneID,
			academyID,
		)

	if err != nil {
		return err
	}

	if !isOwned {
		return shared.ErrUnauthorizedBuildingAccess
	}

	return DeleteAcademyBuildingLaneRepository(
		ctx,
		laneID,
	)
}

func DeleteAcademyBuildingService(
	ctx context.Context,
	userID string,
	buildingID int64,
) error {

	if buildingID <= 0 {
		return shared.ErrInvalidBuildingID
	}

	academyID, err := academyAdmin.
		GetAcademyAdminAcademyID(
			ctx,
			userID,
		)

	if err != nil {
		return err
	}

	isOwned, err :=
		CheckAcademyBuildingOwnershipRepository(
			ctx,
			buildingID,
			academyID,
		)

	if err != nil {
		return err
	}

	if !isOwned {
		return shared.ErrUnauthorizedBuildingAccess
	}

	return DeleteAcademyBuildingRepository(
		ctx,
		buildingID,
	)
}

func GetAvailableBuildingEventsService(
	ctx context.Context,
	userID string,
	buildingID int64,
) ([]EventResponse, error) {

	academyID, err :=
		academyAdmin.
			GetAcademyAdminAcademyID(
				ctx,
				userID,
			)

	if err != nil {
		return nil, err
	}

	isOwned, err :=
		CheckAcademyBuildingOwnershipRepository(
			ctx,
			buildingID,
			academyID,
		)

	if err != nil {
		return nil, err
	}

	if !isOwned {
		return nil,
			shared.ErrUnauthorizedBuildingAccess
	}

	return GetAvailableBuildingEventsRepository(
		ctx,
		buildingID,
	)
}