package districts

import (
	"context"

	"cobackend/internal/db"
)

func GetDistrictsRepository(ctx context.Context) ([]District, error) {
	rows, err := db.DB.Query(
		ctx,
		`
		SELECT id, state_id, district_name
		FROM districts
		ORDER BY district_name ASC
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
	stateID int,
) ([]DistrictResponse, error) {

	rows, err := db.DB.Query(
		ctx,
		`
		SELECT id, district_name
		FROM districts
		WHERE state_id = $1
		ORDER BY district_name ASC
		`,
		stateID,
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

