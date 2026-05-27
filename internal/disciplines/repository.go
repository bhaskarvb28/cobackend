package disciplines

import (
	"context"

	"cobackend/internal/db"
)

func GetDisciplinesRepository(ctx context.Context) ([]Disciplines, error) {

	rows, err := db.DB.Query(
		ctx,
		`
		SELECT 
		id,
		code,
		display_name
		FROM disciplines
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	disciplines := []Disciplines{}

	for rows.Next() {

		var discipline Disciplines

		err := rows.Scan(
			&discipline.ID,
			&discipline.Code,
			&discipline.DisplayName,
		)

		if err != nil {
			return nil, err
		}

		disciplines = append(disciplines, discipline)
	}

	return disciplines, nil

}