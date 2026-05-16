package academyAdmin

import (
	"context"
	"math"
	"strconv"

	"cobackend/internal/db"

	"github.com/jackc/pgx/v5"
)

func CreateAcademyAdminTx(
	ctx context.Context,
	tx pgx.Tx,
	profileID string,
	academyID *int,
	gstin string,
	registrationProof string,
	dpdpConsent bool,
) error {

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO academY_admins (
			profile_id,
			academy_id,
			gstin,
			registration_proof,
			dpdp_consent
		)
		VALUES ($1, $2, $3, $4, $5)
		`,
		profileID,
		academyID,
		gstin,
		registrationProof,
		dpdpConsent,
	)

	if err != nil {
		return err
	}

	return nil
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
		p.id,
		p.first_name,
		p.last_name,
		p.email,
		p.contact_number,
		aa.state_id,
		aa.district_id,
		aa.academy_id,
		aa.gstin,
		aa.registration_proof,
		aa.dpdp_consent,
		aa.created_at::text
	FROM profiles p
	INNER JOIN academy_admins aa
		ON p.id = aa.profile_id
	WHERE 1=1
	`

	countQuery := `
	SELECT COUNT(*)
	FROM profiles p
	INNER JOIN academy_admins aa
		ON p.id = aa.profile_id
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
		)
		`

		baseQuery += searchCondition
		countQuery += searchCondition

		args = append(args, "%"+query.Search+"%")
		argPos++
	}

	if query.StateID != 0 {

		condition := `
		AND aa.state_id = $` + strconv.Itoa(argPos)

		baseQuery += condition
		countQuery += condition

		args = append(args, query.StateID)
		argPos++
	}

	if query.DistrictID != 0 {

		condition := `
		AND aa.district_id = $` + strconv.Itoa(argPos)

		baseQuery += condition
		countQuery += condition

		args = append(args, query.DistrictID)
		argPos++
	}

	if query.AcademyID != 0 {

		condition := `
		AND aa.academy_id = $` + strconv.Itoa(argPos)

		baseQuery += condition
		countQuery += condition

		args = append(args, query.AcademyID)
		argPos++
	}

	var total int

	err := db.DB.QueryRow(
		ctx,
		countQuery,
		args...,
	).Scan(&total)

	if err != nil {
		return PaginatedAcademyAdmins{}, err
	}

	sortColumn := AllowedAcademyAdminSortFields[query.SortBy]

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
		return PaginatedAcademyAdmins{}, err
	}

	defer rows.Close()

	admins := []AcademyAdmin{}

	for rows.Next() {

		var a AcademyAdmin

		err := rows.Scan(
			&a.ID,
			&a.FirstName,
			&a.LastName,
			&a.Email,
			&a.ContactNumber,
			&a.StateID,
			&a.DistrictID,
			&a.AcademyID,
			&a.GSTIN,
			&a.RegistrationProof,
			&a.DPDPConsent,
			&a.CreatedAt,
		)

		if err != nil {
			return PaginatedAcademyAdmins{}, err
		}

		admins = append(admins, a)
	}

	if err := rows.Err(); err != nil {
		return PaginatedAcademyAdmins{}, err
	}

	totalPages := 1

	if query.Limit > 0 {

		totalPages = int(math.Ceil(
			float64(total) / float64(query.Limit),
		))
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
	profileID string,
) (AcademyAdmin, error) {

	var a AcademyAdmin

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT
			p.id,
			p.first_name,
			p.last_name,
			p.email,
			p.contact_number,
			aa.state_id,
			aa.district_id,
			aa.academy_id,
			aa.gstin,
			aa.registration_proof,
			aa.dpdp_consent,
			aa.created_at::text
		FROM profiles p
		INNER JOIN academy_admins aa
			ON p.id = aa.profile_id
		WHERE p.id = $1
		`,
		profileID,
	).Scan(
		&a.ID,
		&a.FirstName,
		&a.LastName,
		&a.Email,
		&a.ContactNumber,
		&a.StateID,
		&a.DistrictID,
		&a.AcademyID,
		&a.GSTIN,
		&a.RegistrationProof,
		&a.DPDPConsent,
		&a.CreatedAt,
	)

	return a, err
}