package academyAdmin

import (
	"context"
	"math"
	"strconv"

	"cobackend/internal/db"

	"github.com/jackc/pgx/v5"

	"cobackend/internal/shared"
)

func CreateAcademyAdminTx(
	ctx context.Context,
	tx pgx.Tx,
	input CreateAcademyAdminInput,
) error {

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO academy_admins (
			user_id,
			academy_id
		)
		VALUES ($1, $2)
		`,
		input.UserID,
		input.AcademyID,
	)

	return err
}


func GetAcademyAdminsRepository(
	ctx context.Context,
	query GetAcademyAdminsQuery,
) (PaginatedAcademyAdmins, error) {

	offset := 0

	if query.Limit > 0 {

		offset = (query.Page - 1) * query.Limit
	}

	baseQuery := `
	SELECT
		u.id,
		u.first_name,
		u.last_name,
		u.email,
		u.contact_number,

		s.id,
		d.id,

		aa.academy_id,

		aa.gstin,
		aa.registration_proof,
		aa.dpdp_consent,

		aa.created_at::text

	FROM users u

	INNER JOIN academy_admins aa
		ON aa.user_id = u.id

	INNER JOIN academies a
		ON aa.academy_id = a.id

	INNER JOIN pincodes pc
		ON pc.id = a.pincode_id

	INNER JOIN districts d
		ON d.id = pc.district_id

	INNER JOIN states s
		ON s.id = d.state_id

	WHERE 1=1
	`

	countQuery := `
	SELECT COUNT(*)

	FROM users u

	INNER JOIN academy_admins aa
		ON aa.user_id = u.id

	INNER JOIN academies a
		ON aa.academy_id = a.id

	INNER JOIN pincodes pc
		ON pc.id = a.pincode_id

	INNER JOIN districts d
		ON d.id = pc.district_id

	INNER JOIN states s
		ON s.id = d.state_id

	WHERE 1=1
	`

	args := []interface{}{}
	argPos := 1

	// ----------------------------------------------------------
	// Search
	// ----------------------------------------------------------

	if query.Search != "" {

		searchCondition := `
		AND (
			u.first_name ILIKE $` + strconv.Itoa(argPos) + `
			OR u.last_name ILIKE $` + strconv.Itoa(argPos) + `
			OR u.email ILIKE $` + strconv.Itoa(argPos) + `
		)
		`

		baseQuery += searchCondition
		countQuery += searchCondition

		args = append(
			args,
			"%"+query.Search+"%",
		)

		argPos++
	}

	// ----------------------------------------------------------
	// State Filter
	// ----------------------------------------------------------

	if query.StateID != 0 {

		condition := `
		AND s.id = $` + strconv.Itoa(argPos)

		baseQuery += condition
		countQuery += condition

		args = append(
			args,
			query.StateID,
		)

		argPos++
	}

	// ----------------------------------------------------------
	// District Filter
	// ----------------------------------------------------------

	if query.DistrictID != 0 {

		condition := `
		AND d.id = $` + strconv.Itoa(argPos)

		baseQuery += condition
		countQuery += condition

		args = append(
			args,
			query.DistrictID,
		)

		argPos++
	}

	// ----------------------------------------------------------
	// Academy Filter
	// ----------------------------------------------------------

	if query.AcademyID != "" {

		condition := `
		AND aa.academy_id = $` + strconv.Itoa(argPos)

		baseQuery += condition
		countQuery += condition

		args = append(
			args,
			query.AcademyID,
		)

		argPos++
	}

	// ----------------------------------------------------------
	// Total Count
	// ----------------------------------------------------------

	var total int

	err := db.DB.QueryRow(
		ctx,
		countQuery,
		args...,
	).Scan(&total)

	if err != nil {

		return PaginatedAcademyAdmins{}, err
	}

	// ----------------------------------------------------------
	// Sorting
	// ----------------------------------------------------------

	sortColumn :=
		AllowedAcademyAdminSortFields[
			query.SortBy,
		]

	baseQuery += `
	ORDER BY ` + sortColumn + ` ` + query.OrderBy

	// ----------------------------------------------------------
	// Pagination
	// ----------------------------------------------------------

	if query.Limit > 0 {

		baseQuery += `
		LIMIT $` + strconv.Itoa(argPos) + `
		OFFSET $` + strconv.Itoa(argPos+1)

		args = append(
			args,
			query.Limit,
			offset,
		)
	}

	rows, err := db.DB.Query(
		ctx,
		baseQuery,
		args...,
	)

	if err != nil {

		return PaginatedAcademyAdmins{}, err
	}

	defer rows.Close()

	admins := []AcademyAdmin{}

	for rows.Next() {

		var academyAdmin AcademyAdmin

		err := rows.Scan(
			&academyAdmin.ID,
			&academyAdmin.FirstName,
			&academyAdmin.LastName,
			&academyAdmin.Email,
			&academyAdmin.ContactNumber,

			&academyAdmin.StateID,
			&academyAdmin.DistrictID,

			&academyAdmin.AcademyID,

			&academyAdmin.GSTIN,
			&academyAdmin.RegistrationProof,
			&academyAdmin.DPDPConsent,

			&academyAdmin.CreatedAt,
		)

		if err != nil {

			return PaginatedAcademyAdmins{}, err
		}

		admins = append(
			admins,
			academyAdmin,
		)
	}

	if err := rows.Err(); err != nil {

		return PaginatedAcademyAdmins{}, err
	}

	totalPages := 1

	if query.Limit > 0 {

		totalPages = int(
			math.Ceil(
				float64(total) /
					float64(query.Limit),
			),
		)
	}

	return PaginatedAcademyAdmins{
		Items:       admins,
		Page:        query.Page,
		Limit:       query.Limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     query.Page < totalPages,
		HasPrevious: query.Page > 1,
	}, nil
}

func GetAcademyAdminByIDRepository(
	ctx context.Context,
	userID string,
) (AcademyAdmin, error) {

	var academyAdmin AcademyAdmin

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT
			u.id,
			u.first_name,
			u.last_name,
			u.email,
			u.contact_number,

			s.id,
			d.id,

			aa.academy_id,

			aa.gstin,
			aa.registration_proof,
			aa.dpdp_consent,

			aa.created_at::text

		FROM users u

		INNER JOIN academy_admins aa
			ON aa.user_id = u.id

		INNER JOIN academies a
			ON a.id = aa.academy_id

		INNER JOIN pincodes pc
			ON pc.id = a.pincode_id

		INNER JOIN districts d
			ON d.id = pc.district_id

		INNER JOIN states s
			ON s.id = d.state_id

		WHERE u.id = $1
		`,
		userID,
	).Scan(
		&academyAdmin.ID,
		&academyAdmin.FirstName,
		&academyAdmin.LastName,
		&academyAdmin.Email,
		&academyAdmin.ContactNumber,

		&academyAdmin.StateID,
		&academyAdmin.DistrictID,

		&academyAdmin.AcademyID,

		&academyAdmin.GSTIN,
		&academyAdmin.RegistrationProof,
		&academyAdmin.DPDPConsent,

		&academyAdmin.CreatedAt,
	)

	return academyAdmin, err
}


func CheckAcademyAdminExists(
	ctx context.Context,
	profileID string,
) (bool, error) {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM academy_admins
			WHERE profile_id = $1
		)
	`

	var exists bool

	err := db.DB.QueryRow(
		ctx,
		query,
		profileID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func UpdateAcademyAdminRepository(
	ctx context.Context,
	profileID string,
	input UpdateAcademyAdminInput,
) error {

	query := `
		UPDATE academy_admins
		SET
			academy_id = COALESCE($1, academy_id)
		WHERE profile_id = $2
	`

	commandTag, err := db.DB.Exec(
		ctx,
		query,
		input.AcademyID,
		profileID,
	)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return shared.ErrAcademyAdminNotFound
	}

	return nil
}

func GetAcademyAdminDistrictID(
	ctx context.Context,
	profileID string,
) (int, error) {

	query := `
		SELECT a.district_id
		FROM academy_admins aa
		INNER JOIN academies a
			ON aa.academy_id = a.id
		WHERE aa.profile_id = $1
	`

	var districtID int

	err := db.DB.QueryRow(
		ctx,
		query,
		profileID,
	).Scan(&districtID)

	if err != nil {
		return 0, err
	}

	return districtID, nil
}

func DeleteAcademyAdminRepository(
	ctx context.Context,
	profileID string,
) error {

	query := `
		DELETE FROM academy_admins
		WHERE profile_id = $1
	`

	commandTag, err := db.DB.Exec(
		ctx,
		query,
		profileID,
	)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return shared.ErrAcademyAdminNotFound
	}

	return nil
}

func GetAcademyAdminAcademyID(
	ctx context.Context,
	profileID string,
) (string, error) {

	var academyID string

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT 
		academy_id
		FROM academy_admins
		WHERE user_id = $1
		`,
		profileID,
	).Scan(
		&academyID,
	)

	return academyID, err
}