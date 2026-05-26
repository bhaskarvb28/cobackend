package districtCoach

import (
	"context"
	"github.com/jackc/pgx/v5"

	"math"
	"strconv"

	"cobackend/internal/db"
	"cobackend/internal/shared"

	"errors"

)

func CreateDistrictCoachTx(
	ctx context.Context,
	tx pgx.Tx,
	input CreateDistrictCoachInput,
) error {

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO district_coaches (
			user_id,
			district_id
		)
		VALUES (
			$1,
			$2
		)
		`,
		input.UserID,
		input.DistrictID,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetDistrictCoachesRepository(
	ctx context.Context,
	query GetDistrictCoachesQuery,
) (PaginatedDistrictCoaches, error) {

	offset := 0

	if query.Limit > 0 {
		offset = (query.Page - 1) * query.Limit
	}

	baseQuery := `
	SELECT
		p.id,
		p.first_name,
		p.last_name,
		p.email,
		p.contact_number,
		d.state_id,
		dc.district_id,
		dc.coach_code,
		dc.coaching_certificate_proof,
		dc.dpdp_consent
	FROM profiles p
	INNER JOIN district_coaches dc
		ON p.id = dc.profile_id
	INNER JOIN districts d
		ON d.id = dc.district_id
	WHERE 1=1
	`

	countQuery := `
	SELECT COUNT(*)
	FROM profiles p
	INNER JOIN district_coaches dc
		ON p.id = dc.profile_id
	INNER JOIN districts d
		ON d.id = dc.district_id
	WHERE 1=1
	`

	args := []interface{}{}
	argPos := 1

	if query.Search != "" {

		searchCondition := `
		AND (
			p.first_name ILIKE $` + strconv.Itoa(argPos) + `
			OR p.last_name ILIKE $` + strconv.Itoa(argPos) + `
			OR p.email ILIKE $` + strconv.Itoa(argPos) + `
			OR dc.coach_code ILIKE $` + strconv.Itoa(argPos) + `
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
		AND dc.district_id = $` + strconv.Itoa(argPos)

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
		return PaginatedDistrictCoaches{}, err
	}

	sortColumn := AllowedDistrictCoachSortFields[query.SortBy]

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
		return PaginatedDistrictCoaches{}, err
	}

	defer rows.Close()

	coaches := []DistrictCoach{}

	for rows.Next() {

		var c DistrictCoach

		err := rows.Scan(
			&c.ID,
			&c.FirstName,
			&c.LastName,
			&c.Email,
			&c.ContactNumber,
			&c.StateID,
			&c.DistrictID,
			&c.CoachCode,
			&c.CoachCertificationProof,
			&c.DPDPConsent,
		)

		if err != nil {
			return PaginatedDistrictCoaches{}, err
		}

		coaches = append(coaches, c)
	}

	if err := rows.Err(); err != nil {
		return PaginatedDistrictCoaches{}, err
	}

	totalPages := 1

	if query.Limit > 0 {

		totalPages = int(math.Ceil(
			float64(total) / float64(query.Limit),
		))
	}

	return PaginatedDistrictCoaches{
		Items:       coaches,
		Page:        query.Page,
		Limit:       query.Limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     query.Page < totalPages,
		HasPrevious: query.Page > 1,
	}, nil
}


func GetDistrictCoachByProfileID(
	ctx context.Context,
	profileID string,
) (DistrictCoach, error) {

	var districtCoach DistrictCoach

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT
			dc.profile_id,
			dc.district_id,
			d.state_id
		FROM district_coaches dc
		INNER JOIN districts d
			ON d.id = dc.district_id
		WHERE dc.profile_id = $1
		`,
		profileID,
	).Scan(
		&districtCoach.ID,
		&districtCoach.DistrictID,
		&districtCoach.StateID,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return DistrictCoach{}, shared.ErrDistrictCoachNotFound
		}

		return DistrictCoach{}, err
	}

	return districtCoach, nil
}

func DeleteDistrictCoachRepository(
	ctx context.Context,
	profileID string,
) error {

	commandTag, err := db.DB.Exec(
		ctx,
		`
		DELETE FROM profiles
		WHERE id = $1
		`,
		profileID,
	)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return shared.ErrDistrictCoachNotFound
	}

	return nil
}

func CheckDistrictCoachExists(
	ctx context.Context,
	id string,
) (bool, error) {

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM district_coaches
			WHERE profile_id = $1
		)
	`

	var exists bool

	err := db.DB.QueryRow(
		ctx,
		query,
		id,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func UpdateDistrictCoachRepository(
	ctx context.Context,
	id string,
	input UpdateDistrictCoachInput,
) error {

	query := `
		UPDATE district_coaches
		SET
			district_id = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE profile_id = $2
	`

	commandTag, err := db.DB.Exec(
		ctx,
		query,
		*input.DistrictID,
		id,
	)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return shared.ErrDistrictCoachNotFound
	}

	return nil
}