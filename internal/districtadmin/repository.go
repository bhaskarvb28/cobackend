package districtadmin

import (
	"context"

	"github.com/jackc/pgx/v5"

	"math"
	"strconv"

	"cobackend/internal/db"
)

func CreateDistrictAdminTx(
	ctx context.Context,
	tx pgx.Tx,
	profileID string,
	stateID *string,
	districtID *string,
	dpdpConsent bool,
) error {

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO district_admins (
			profile_id,
			state_id,
			district_id,
			dpdp_consent
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
		`,
		profileID,
		stateID,
		districtID,
		dpdpConsent,
	)

	if err != nil {
		return err
	}

	return nil
}

func GetDistrictAdminsRepository(
	ctx context.Context,
	query GetDistrictAdminsQuery,
) (PaginatedDistrictAdmins, error) {

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
		da.state_id,
		da.district_id,
		da.dpdp_consent
	FROM profiles p
	INNER JOIN district_admins da
		ON p.id = da.profile_id
	WHERE 1=1
	`

	countQuery := `
	SELECT COUNT(*)
	FROM profiles p
	INNER JOIN district_admins da
		ON p.id = da.profile_id
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
		AND da.state_id = $` + strconv.Itoa(argPos)

		baseQuery += condition
		countQuery += condition

		args = append(args, query.StateID)
		argPos++
	}

	if query.DistrictID != 0 {

		condition := `
		AND da.district_id = $` + strconv.Itoa(argPos)

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
		return PaginatedDistrictAdmins{}, err
	}

	baseQuery += `
	ORDER BY ` + query.SortBy + ` ` + query.OrderBy

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
		return PaginatedDistrictAdmins{}, err
	}

	defer rows.Close()

	admins :=  []DistrictAdmin{}

	for rows.Next() {

		var a DistrictAdmin

		err := rows.Scan(
			&a.ID,
			&a.FirstName,
			&a.LastName,
			&a.Email,
			&a.ContactNumber,
			&a.StateID,
			&a.DistrictID,
			&a.DPDPConsent,
		)

		if err != nil {
			return PaginatedDistrictAdmins{}, err
		}

		admins = append(admins, a)
	}

	totalPages := 1

	if query.Limit > 0 {

		totalPages = int(math.Ceil(
			float64(total) / float64(query.Limit),
		))
	}

	return PaginatedDistrictAdmins{
		Items:       admins,
		Page:        query.Page,
		Limit:       query.Limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     query.Page < totalPages,
		HasPrevious: query.Page > 1,
	}, nil
}

// import (
// 	"context"
// 	"errors"
// 	"strconv"

// 	"cobackend/internal/db"

// 	"github.com/google/uuid"
// 	"github.com/jackc/pgx/v5/pgconn"
// )

// func CreateDistrictAdminRepository(
// 	ctx context.Context,
// 	input CreateDistrictAdminInput,
// 	hashedPassword string,
// ) error {

// 	tx, err := db.DB.Begin(ctx)
// 	if err != nil {
// 		return errors.New("failed to start database transaction")
// 	}
// 	defer tx.Rollback(ctx)

// 	var districtAdminRoleID string
// 	err = tx.QueryRow(
// 		ctx,
// 		`SELECT role_id FROM roles WHERE role_name = 'district_admin'`,
// 	).Scan(&districtAdminRoleID)
// 	if err != nil {
// 		return errors.New("failed to fetch district admin role")
// 	}

// 	profileID := uuid.New()

// 	_, err = tx.Exec(
// 		ctx,
// 		`
// 		INSERT INTO profiles (
// 			id,
// 			first_name,
// 			last_name,
// 			email,
// 			password,
// 			contact_number,
// 			role_id
// 		)
// 		VALUES ($1, $2, $3, $4, $5, $6, $7)
// 		`,
// 		profileID,
// 		input.FirstName,
// 		input.LastName,
// 		input.Email,
// 		hashedPassword,
// 		input.ContactNumber,
// 		districtAdminRoleID,
// 	)

// 	if err != nil {
// 		var pgErr *pgconn.PgError
// 		if errors.As(err, &pgErr) {
// 			switch pgErr.Code {
// 			case "23505":
// 				if pgErr.ConstraintName == "profiles_email_key" {
// 					return errors.New("email already exists")
// 				}
// 				return errors.New("duplicate value already exists")
// 			case "23503":
// 				return errors.New("invalid foreign key reference")
// 			case "23502":
// 				return errors.New("required field is missing")
// 			case "42P01":
// 				return errors.New("required database table does not exist")
// 			case "42703":
// 				return errors.New("required database column does not exist")
// 			default:
// 				return errors.New("database operation failed")
// 			}
// 		}
// 		return errors.New("failed to create profile")
// 	}

// 	_, err = tx.Exec(
// 		ctx,
// 		`
// 		INSERT INTO district_admins (
// 			profile_id,
// 			state_id,
// 			district_id
// 		)
// 		VALUES ($1, $2, $3)
// 		`,
// 		profileID,
// 		input.StateID,
// 		input.DistrictID,
// 	)

// 	if err != nil {
// 		var pgErr *pgconn.PgError
// 		if errors.As(err, &pgErr) {
// 			switch pgErr.Code {
// 			case "23503":
// 				return errors.New("invalid state or district id")
// 			case "23505":
// 				return errors.New("district admin already exists")
// 			case "23502":
// 				return errors.New("state and district are required")
// 			case "42P01":
// 				return errors.New("required database table does not exist")
// 			default:
// 				return errors.New("failed to assign district admin")
// 			}
// 		}
// 		return errors.New("failed to create district admin")
// 	}

// 	err = tx.Commit(ctx)
// 	if err != nil {
// 		return errors.New("failed to commit database transaction")
// 	}

// 	return nil
// }



// func UpdateDistrictAdminRepository(
// 	ctx context.Context,
// 	id string,
// 	input UpdateDistrictAdminInput,
// ) error {

// 	_, err := db.DB.Exec(
// 		ctx,
// 		`
// 		UPDATE district_admins
// 		SET
// 			state_id        = $1,
// 			district_id     = $2,
// 			approval_status = $3,
// 			approval_notes  = $4,
// 			updated_at      = CURRENT_TIMESTAMP
// 		WHERE profile_id = $5
// 		`,
// 		input.StateID,
// 		input.DistrictID,
// 		input.ApprovalStatus,
// 		input.ApprovalNotes,
// 		id,
// 	)

// 	if err != nil {
// 		var pgErr *pgconn.PgError
// 		if errors.As(err, &pgErr) {
// 			switch pgErr.Code {
// 			case "23503":
// 				return errors.New("invalid state or district id")
// 			case "23514":
// 				return errors.New("invalid approval status or missing rejection notes")
// 			case "42P01":
// 				return errors.New("required database table does not exist")
// 			default:
// 				return errors.New("failed to update district admin")
// 			}
// 		}
// 		return errors.New("database operation failed")
// 	}

// 	return nil
// }

// func DeleteDistrictAdminRepository(
// 	ctx context.Context,
// 	id string,
// ) error {

// 	commandTag, err := db.DB.Exec(
// 		ctx,
// 		`
// 		DELETE FROM profiles
// 		WHERE id = $1
// 		`,
// 		id,
// 	)

// 	if err != nil {
// 		var pgErr *pgconn.PgError
// 		if errors.As(err, &pgErr) {
// 			switch pgErr.Code {
// 			case "22P02":
// 				return errors.New("invalid district admin id")
// 			case "42P01":
// 				return errors.New("required database table does not exist")
// 			default:
// 				return errors.New("failed to delete district admin")
// 			}
// 		}
// 		return errors.New("database operation failed")
// 	}

// 	if commandTag.RowsAffected() == 0 {
// 		return errors.New("district admin not found")
// 	}

// 	return nil
// }
