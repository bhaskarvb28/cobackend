package states

import (
	"cobackend/internal/db"
	"context"
	"strings"
	"fmt"
)

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

	// Search filter
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

	// Order validation
	order := "ASC"

	if strings.ToUpper(queryParams.Order) == "DESC" {
		order = "DESC"
	}

	query += fmt.Sprintf(
		" ORDER BY name %s",
		order,
	)

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

	for rows.Next() {
		var state State

		err := rows.Scan(
			&state.ID,
			&state.Name,
		)

		if err != nil {
			return nil, err
		}

		states = append(states, state)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return states, nil
}

func CheckStateExists(
	ctx context.Context,
	stateID string,
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