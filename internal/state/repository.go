package state

import (
	"context"
	"fmt"
	"strings"

	"cobackend/internal/db"
)

// GetStatesRepository fetches all states
// based on the provided query parameters.
func GetStatesRepository(
	ctx context.Context,
	queryParams GetStatesQueryParams,
) ([]State, error) {

	query := `
		SELECT id, name
		FROM states
	`

	args := []interface{}{}
	argPos := 1

	// ----------------------------------------------------------
	// Apply Search Filter
	// ----------------------------------------------------------

	if queryParams.Search != "" {

		query += fmt.Sprintf(
			" WHERE name ILIKE $%d",
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

	states := []State{}

	// ----------------------------------------------------------
	// Scan Rows
	// ----------------------------------------------------------

	for rows.Next() {

		var state State

		err := rows.Scan(
			&state.ID,
			&state.Name,
		)

		if err != nil {
			return nil, err
		}

		states = append(
			states,
			state,
		)
	}

	// ----------------------------------------------------------
	// Check Row Iteration Errors
	// ----------------------------------------------------------

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return states, nil
}

// CheckStateExists checks whether a state
// exists for the provided state ID.
func CheckStateExists(
	ctx context.Context,
	stateID int,
) (bool, error) {

	var exists bool

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT EXISTS (
			SELECT 1
			FROM states
			WHERE id = $1
		)
		`,
		stateID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}