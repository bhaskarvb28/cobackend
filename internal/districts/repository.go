package districts

import (
	"context"

	"cobackend/internal/db"
		"fmt"

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
		fmt.Println(err)
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