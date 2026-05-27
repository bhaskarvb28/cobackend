package academy

import (
	"context"

	"cobackend/internal/db"
	"math"
	"strconv"
)

func CreateAcademyRepository(
	ctx context.Context,
	districtID int,
	input CreateAcademyInput,
) (*AcademyResponse, error) {

	var academy AcademyResponse

	err := db.DB.QueryRow(
		ctx,
		`
		INSERT INTO academies (
			name,
			district_id,
			address
		)
		VALUES (
			$1,
			$2,
			$3
		)
		RETURNING
			id,
			name,
			district_id,
			address,
			is_active,
			created_at,
			updated_at
		`,
		input.Name,
		districtID,
		input.Address,
	).Scan(
		&academy.ID,
		&academy.Name,
		&academy.DistrictID,
		&academy.Address,
		&academy.IsActive,
		&academy.CreatedAt,
		&academy.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &academy, nil
}


func CheckAcademyExists(
	ctx context.Context,
	academyID string,
) (bool, error) {

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM academies
			WHERE id = $1
		)
	`

	var exists bool

	err := db.DB.QueryRow(
		ctx,
		query,
		academyID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func CheckAcademyBelongsToDistrict(
	ctx context.Context,
	academyID int,
	districtID int,
) (bool, error) {

	var exists bool

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT EXISTS(
			SELECT 1
			FROM academies
			WHERE id = $1
			AND district_id = $2
		)
		`,
		academyID,
		districtID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetAcademiesRepository queries the database for paginated academy records.
// It uses multi-step SQL queries, first counting matching records to calculate total pagination metrics,
// and then fetching matching rows using custom limits, offsets, sorting, and state/district joins.
func GetAcademiesRepository(
	ctx context.Context,
	query GetAcademiesQuery,
) (PaginatedAcademies, error) {

	// ----------------------------------------------------------
	// Calculate Offset
	// ----------------------------------------------------------
	// Determines how many rows to skip based on page size (limit) and page number.

	offset := 0

	if query.Limit > 0 {
		offset = (query.Page - 1) * query.Limit
	}

	// ----------------------------------------------------------
	// Build Database Queries
	// ----------------------------------------------------------
	// baseQuery fetches detailed records with a JOIN on districts to get the state_id.
	// countQuery executes a matching SELECT COUNT(*) to calculate pagination pages.

	baseQuery := `
	SELECT
		a.id,
		a.name,
		d.state_id,
		a.district_id,
		a.address,
		a.created_at,
		a.updated_at
	FROM academies a
	INNER JOIN districts d
		ON a.district_id = d.id
	WHERE 1=1
	`

	countQuery := `
	SELECT COUNT(*)
	FROM academies a
	INNER JOIN districts d
		ON a.district_id = d.id
	WHERE 1=1
	`

	args := []interface{}{}
	argPos := 1

	// ----------------------------------------------------------
	// Apply Search Filter
	// ----------------------------------------------------------

	if query.Search != "" {

		searchCondition := `
		AND (
			a.name ILIKE $` + strconv.Itoa(argPos) + `
			OR a.address ILIKE $` + strconv.Itoa(argPos) + `
		)
		`

		baseQuery += searchCondition
		countQuery += searchCondition

		args = append(args, "%"+query.Search+"%")
		argPos++
	}

	// ----------------------------------------------------------
	// Apply State Filter
	// ----------------------------------------------------------

	if query.StateID != 0 {

		condition := `
		AND d.state_id = $` + strconv.Itoa(argPos)

		baseQuery += condition
		countQuery += condition

		args = append(args, query.StateID)
		argPos++
	}

	// ----------------------------------------------------------
	// Apply District Filter
	// ----------------------------------------------------------

	if query.DistrictID != 0 {

		condition := `
		AND a.district_id = $` + strconv.Itoa(argPos)

		baseQuery += condition
		countQuery += condition

		args = append(args, query.DistrictID)
		argPos++
	}

	// ----------------------------------------------------------
	// Query Total Count
	// ----------------------------------------------------------

	var total int

	err := db.DB.QueryRow(
		ctx,
		countQuery,
		args...,
	).Scan(&total)

	if err != nil {
		return PaginatedAcademies{}, err
	}

	// ----------------------------------------------------------
	// Query Paginated Records
	// ----------------------------------------------------------

	sortColumn := AllowedAcademySortFields[query.SortBy]

	baseQuery += `
	ORDER BY ` + sortColumn + ` ` + query.OrderBy

	if query.Limit > 0 {

		baseQuery += `
		LIMIT $` + strconv.Itoa(argPos) + `
		OFFSET $` + strconv.Itoa(argPos+1)

		args = append(args, query.Limit, offset)
	}

	rows, err := db.DB.Query(
		ctx,
		baseQuery,
		args...,
	)

	if err != nil {
		return PaginatedAcademies{}, err
	}

	defer rows.Close()

	// ----------------------------------------------------------
	// Scan Records
	// ----------------------------------------------------------

	academies := []AcademyResponse{}

	for rows.Next() {

		var a AcademyResponse

		err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.StateID,
			&a.DistrictID,
			&a.Address,
			&a.CreatedAt,
			&a.UpdatedAt,
		)

		if err != nil {
			return PaginatedAcademies{}, err
		}

		academies = append(academies, a)
	}

	if err := rows.Err(); err != nil {
		return PaginatedAcademies{}, err
	}

	// ----------------------------------------------------------
	// Paginated Result
	// ----------------------------------------------------------

	totalPages := 1

	if query.Limit > 0 {

		totalPages = int(math.Ceil(
			float64(total) / float64(query.Limit),
		))
	}

	return PaginatedAcademies{
		Items:       academies,
		Page:        query.Page,
		Limit:       query.Limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     query.Page < totalPages,
		HasPrevious: query.Page > 1,
	}, nil
}

func CreateAcademyBuildingRepository(
	ctx context.Context,
	academyID string,
	input CreateAcademyBuildingInput,
) (*AcademyBuildingResponse, error) {

	var academyBuilding AcademyBuildingResponse

	err := db.DB.QueryRow(
		ctx,
		`
		INSERT INTO academy_buildings (
			academy_id,
			building_name
		)
		VALUES (
			$1,
			$2
		)
		RETURNING
			id,
			academy_id,
			building_name,
			is_active
		`,
		academyID,
		input.BuildingName,
	).Scan(
		&academyBuilding.ID,
		&academyBuilding.AcademyID,
		&academyBuilding.BuildingName,
		&academyBuilding.IsActive,
	)

	if err != nil {
		return nil, err
	}

	return &academyBuilding, nil
}

func CheckAcademyBuildingOwnershipRepository(
	ctx context.Context,
	buildingID int64,
	academyID string,
) (bool, error) {

	var exists bool

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT EXISTS (
			SELECT 1
			FROM academy_buildings
			WHERE id = $1
			AND academy_id = $2
		)
		`,
		buildingID,
		academyID,
	).Scan(
		&exists,
	)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func AddAcademyBuildingDisciplineRepository(
	ctx context.Context,
	buildingID int64,
	disciplineID int,
) (*AcademyBuildingDisciplineResponse, error) {

	var response AcademyBuildingDisciplineResponse

	err := db.DB.QueryRow(
		ctx,
		`
		INSERT INTO academy_building_disciplines (
			academy_building_id,
			discipline_id
		)
		VALUES (
			$1,
			$2
		)
		RETURNING
			academy_building_id,
			discipline_id
		`,
		buildingID,
		disciplineID,
	).Scan(
		&response.AcademyBuildingID,
		&response.DisciplineID,
	)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ============================================================================
// repository.go
// ============================================================================

func AddAcademyBuildingEventRepository(
	ctx context.Context,
	buildingID int64,
	shootingEventID int,
) (*AcademyBuildingEventResponse, error) {

	var response AcademyBuildingEventResponse

	err := db.DB.QueryRow(
		ctx,
		`
		INSERT INTO academy_building_events (
			academy_building_id,
			shooting_event_id
		)
		VALUES (
			$1,
			$2
		)
		RETURNING
			academy_building_id,
			shooting_event_id
		`,
		buildingID,
		shootingEventID,
	).Scan(
		&response.AcademyBuildingID,
		&response.ShootingEventID,
	)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ============================================================================
// repository.go
// ============================================================================

func GetAcademyBuildingsRepository(
	ctx context.Context,
	academyID string,
) ([]AcademyBuilding, error) {

	// ----------------------------------------------------------
	// Get Buildings
	// ----------------------------------------------------------

	rows, err := db.DB.Query(
		ctx,
		`
		SELECT
			id,
			academy_id,
			building_name,
			is_active
		FROM academy_buildings
		WHERE academy_id = $1
		ORDER BY id DESC
		`,
		academyID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	buildings := []AcademyBuilding{}

	for rows.Next() {

		var building AcademyBuilding

		err := rows.Scan(
			&building.ID,
			&building.AcademyID,
			&building.BuildingName,
			&building.IsActive,
		)

		if err != nil {
			return nil, err
		}

		// ----------------------------------------------------------
		// Get Building Disciplines
		// ----------------------------------------------------------

		disciplines, err := GetAcademyBuildingDisciplinesRepository(
			ctx,
			building.ID,
		)

		if err != nil {
			return nil, err
		}

		building.Disciplines = disciplines

		// ----------------------------------------------------------
		// Get Building Events
		// ----------------------------------------------------------

		events, err := GetAcademyBuildingEventsRepository(
			ctx,
			building.ID,
		)

		if err != nil {
			return nil, err
		}

		building.Events = events

		buildings = append(
			buildings,
			building,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return buildings, nil
}

func GetAcademyBuildingDisciplinesRepository(
	ctx context.Context,
	buildingID int64,
) ([]BuildingDiscipline, error) {

	rows, err := db.DB.Query(
		ctx,
		`
		SELECT
			d.id,
			d.code,
			d.display_name
		FROM academy_building_disciplines abd

		INNER JOIN disciplines d
			ON d.id = abd.discipline_id

		WHERE abd.academy_building_id = $1

		ORDER BY d.display_name
		`,
		buildingID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	disciplines := []BuildingDiscipline{}

	for rows.Next() {

		var discipline BuildingDiscipline

		err := rows.Scan(
			&discipline.ID,
			&discipline.Code,
			&discipline.DisplayName,
		)

		if err != nil {
			return nil, err
		}

		disciplines = append(
			disciplines,
			discipline,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return disciplines, nil
}

func GetAcademyBuildingEventsRepository(
	ctx context.Context,
	buildingID int64,
) ([]BuildingEvent, error) {

	rows, err := db.DB.Query(
		ctx,
		`
		SELECT
			se.id,
			se.code,
			se.display_name
		FROM academy_building_events abe

		INNER JOIN shooting_events se
			ON se.id = abe.shooting_event_id

		WHERE abe.academy_building_id = $1

		ORDER BY se.display_name
		`,
		buildingID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []BuildingEvent{}

	for rows.Next() {

		var event BuildingEvent

		err := rows.Scan(
			&event.ID,
			&event.Code,
			&event.DisplayName,
		)

		if err != nil {
			return nil, err
		}

		events = append(
			events,
			event,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}