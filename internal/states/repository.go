package states

import (
	"cobackend/internal/db"
	"context"
)

func GetStatesRepository(ctx context.Context) ([]State, error) {

	rows, err := db.DB.Query(
		ctx,
		`
		SELECT id, state_name
		FROM states
		ORDER BY state_name ASC
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var states []State

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

	return states, nil
}

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