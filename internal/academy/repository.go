package academy

import (
	"context"

	"cobackend/internal/db"
	"math"
	"strconv"
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

func GetAcademiesRepository(
	ctx context.Context,
	query GetAcademiesQuery,
) (PaginatedAcademies, error) {

	offset := 0

	if query.Limit > 0 {
		offset = (query.Page - 1) * query.Limit
	}

	baseQuery := `
	SELECT
		a.id,
		a.name,
		d.state_id,
		a.district_id,
		a.address,
		a.created_at::text,
		a.updated_at::text
	FROM academies a
	INNER JOIN districts d
		ON a.district_id = d.id
	WHERE 1=1
	`

	countQuery := `
	SELECT COUNT(*)
	FROM academies a
	INNER JOIN districts d
		ON a.district_id = d.id
	WHERE 1=1
	`

	args := []interface{}{}
	argPos := 1

	if query.Search != "" {

		searchCondition := `
		AND (
			a.name ILIKE $` + strconv.Itoa(argPos) + `
			OR a.address ILIKE $` + strconv.Itoa(argPos) + `
		)
		`

		baseQuery += searchCondition
		countQuery += searchCondition

		args = append(args, "%"+query.Search+"%")
		argPos++
	}

	if query.StateID != 0 {

		condition := `
		AND d.state_id = $` + strconv.Itoa(argPos)

		baseQuery += condition
		countQuery += condition

		args = append(args, query.StateID)
		argPos++
	}

	if query.DistrictID != 0 {

		condition := `
		AND a.district_id = $` + strconv.Itoa(argPos)

		baseQuery += condition
		countQuery += condition

		args = append(args, query.DistrictID)
		argPos++
	}

	var total int

	err := db.DB.QueryRow(
		ctx,
		countQuery,
		args...,
	).Scan(&total)

	if err != nil {
		return PaginatedAcademies{}, err
	}

	sortColumn := AllowedAcademySortFields[query.SortBy]

	baseQuery += `
	ORDER BY ` + sortColumn + ` ` + query.OrderBy

	if query.Limit > 0 {

		baseQuery += `
		LIMIT $` + strconv.Itoa(argPos) + `
		OFFSET $` + strconv.Itoa(argPos+1)

		args = append(args, query.Limit, offset)
	}

	rows, err := db.DB.Query(
		ctx,
		baseQuery,
		args...,
	)

	if err != nil {
		return PaginatedAcademies{}, err
	}

	defer rows.Close()

	academies := []Academy{}

	for rows.Next() {

		var a Academy

		err := rows.Scan(
			&a.ID,
			&a.Name,
			&a.StateID,
			&a.DistrictID,
			&a.Address,
			&a.CreatedAt,
			&a.UpdatedAt,
		)

		if err != nil {
			return PaginatedAcademies{}, err
		}

		academies = append(academies, a)
	}

	if err := rows.Err(); err != nil {
		return PaginatedAcademies{}, err
	}

	totalPages := 1

	if query.Limit > 0 {

		totalPages = int(math.Ceil(
			float64(total) / float64(query.Limit),
		))
	}

	return PaginatedAcademies{
		Items:       academies,
		Page:        query.Page,
		Limit:       query.Limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     query.Page < totalPages,
		HasPrevious: query.Page > 1,
	}, nil
}