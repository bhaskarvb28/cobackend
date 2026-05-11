package districts

import (
	"context"

	"cobackend/internal/db"

	"fmt"

	"strings"
)

func GetDistrictsRepository(ctx context.Context) ([]District, error) {
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

	var districts [] District

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

		districts = append(districts, district)
	}

	return districts, nil
}

func GetDistrictsByStateIDRepository(
	ctx context.Context,
	stateID string,
	queryParams GetDistrictQueryParams,
) ([]DistrictResponse, error) {

	query := `
		SELECT id, name
		FROM districts
		WHERE state_id = $1
	`

	args := []interface{}{stateID}
	argPos := 2

	// Search filter
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

	districts := []DistrictResponse{}

	for rows.Next() {
		var district DistrictResponse

		err := rows.Scan(
			&district.ID,
			&district.Name,
		)

		if err != nil {
			return nil, err
		}

		districts = append(districts, district)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return districts, nil
}
