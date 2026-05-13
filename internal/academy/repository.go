package academy

import (
	"context"

	"cobackend/internal/db"
)

func CreateAcademyRepository(
	ctx context.Context,
	input CreateAcademyInput,
) error {

	_, err := db.DB.Exec(
		ctx,
		`
		INSERT INTO academies (
			name,
			district_id,
			address
		)
		VALUES ($1, $2, $3)
		`,
		input.Name,
		input.DistrictID,
		input.Address,
	)

	return err
}