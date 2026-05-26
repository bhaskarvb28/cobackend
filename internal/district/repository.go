package district

import (
	"context"
	"fmt"
	"strings"

	"cobackend/internal/db"
)

// GetDistrictsRepository fetches all districts.
func GetDistrictsRepository(
	ctx context.Context,
) ([]District, error) {

	// ----------------------------------------------------------
	// Execute Query
	// ----------------------------------------------------------

	rows, err := db.DB.Query(
		ctx,
		`
		SELECT id, state_id, name
		FROM districts
		ORDER BY name ASC
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	districts := []District{}

	// ----------------------------------------------------------
	// Scan Rows
	// ----------------------------------------------------------

	for rows.Next() {

		var district District

		err := rows.Scan(
			&district.ID,
			&district.StateID,
			&district.DistrictName,
		)

		if err != nil {
			return nil, err
		}

		districts = append(
			districts,
			district,
		)
	}

	// ----------------------------------------------------------
	// Check Row Iteration Errors
	// ----------------------------------------------------------

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return districts, nil
}

// GetDistrictsByStateIDRepository fetches all districts
// belonging to a specific state based on
// the provided query parameters.
func GetDistrictsByStateIDRepository(
	ctx context.Context,
	stateID int,
	queryParams GetDistrictQueryParams,
) ([]DistrictResponse, error) {

	query := `
		SELECT id, name
		FROM districts
		WHERE state_id = $1
	`

	args := []interface{}{stateID}
	argPos := 2

	// ----------------------------------------------------------
	// Apply Search Filter
	// ----------------------------------------------------------

	if queryParams.Search != "" {

		query += fmt.Sprintf(
			" AND name ILIKE $%d",
			argPos,
		)

		args = append(
			args,
			"%"+queryParams.Search+"%",
		)

		argPos++
	}

	// ----------------------------------------------------------
	// Apply Sorting
	// ----------------------------------------------------------

	order := "ASC"

	if strings.ToUpper(queryParams.Order) == "DESC" {
		order = "DESC"
	}

	query += fmt.Sprintf(
		" ORDER BY name %s",
		order,
	)

	// ----------------------------------------------------------
	// Execute Query
	// ----------------------------------------------------------

	rows, err := db.DB.Query(
		ctx,
		query,
		args...,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	districts := []DistrictResponse{}

	// ----------------------------------------------------------
	// Scan Rows
	// ----------------------------------------------------------

	for rows.Next() {

		var district DistrictResponse

		err := rows.Scan(
			&district.ID,
			&district.Name,
		)

		if err != nil {
			return nil, err
		}

		districts = append(
			districts,
			district,
		)
	}

	// ----------------------------------------------------------
	// Check Row Iteration Errors
	// ----------------------------------------------------------

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return districts, nil
}

// CheckDistrictExists checks whether a district
// exists for the provided district ID.
func CheckDistrictExists(
	ctx context.Context,
	districtID int,
) (bool, error) {

	var exists bool

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT EXISTS (
			SELECT 1
			FROM districts
			WHERE id = $1
		)
		`,
		districtID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetStateIDByDistrictID fetches the state ID
// associated with the provided district ID.
func GetStateIDByDistrictID(
	ctx context.Context,
	districtID int,
) (int, error) {

	var stateID int

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT state_id
		FROM districts
		WHERE id = $1
		`,
		districtID,
	).Scan(&stateID)

	if err != nil {
		return 0, err
	}

	return stateID, nil
}

