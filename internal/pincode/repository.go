package pincode

import (
	"cobackend/internal/db"
	"context"
)

func GetPincodesRepository(ctx context.Context) ([]Pincode, error) {
	rows, err := db.DB.Query(
		ctx,
		`SELECT
		p.id,
		p.code,
		d.id as district_id,
		d.name as district_name,
		s.name as state_name

		FROM pincodes p

		INNER JOIN districts d
			ON d.id = p.district_id

		INNER JOIN states s
			ON s.id = d.state_id
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
			&pincode.DistrictName,
			&pincode.StateName,
		)

		if err != nil {
			return nil, err
		}

		pincodes = append(pincodes, pincode)
	}

	return pincodes, nil
}