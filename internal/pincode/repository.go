package pincode

import (
	"cobackend/internal/db"
	"context"
)

func GetPincodesRepository(ctx context.Context) ([]Pincode, error) {
	rows, err := db.DB.Query(
		ctx,
		`SELECT
		id,
		code,
		district_id
		FROM pincodes
		`,
	)

	if err != nil {
		return nil, err
	}

	pincodes := []Pincode{}

	for rows.Next() {
		var pincode Pincode;

		err = rows.Scan(
			&pincode.ID,
			&pincode.Code,
			&pincode.DistrictID,
		)

		if err != nil {
			return nil, err
		}

		pincodes = append(pincodes, pincode)
	}

	return pincodes, nil
}