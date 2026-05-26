package districtAdmin

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"math"
	"strconv"

	"cobackend/internal/db"
	"cobackend/internal/shared"

)

func CreateDistrictAdminTx(
	ctx context.Context,
	tx pgx.Tx,
	input CreateDistrictAdminInput,
) error {

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO district_admins (
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

func CheckDistrictAdminExists(
	ctx context.Context,
	profileID string,
) (bool, error) {

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM district_admins
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

func UpdateDistrictAdminRepository(
	ctx context.Context,
	profileID string,
	input UpdateDistrictAdminInput,
) error {

	query := `
		UPDATE district_admins
		SET
			district_id = $1
		WHERE profile_id = $2
	`

	commandTag, err := db.DB.Exec(
		ctx,
		query,
		input.DistrictID,
		profileID,
	)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return shared.ErrDistrictAdminNotFound
	}

	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// NEW — GetDistrictAdminStateID
// Fetches which state the district admin belongs to.
// Used to verify the state admin is deleting someone in their own state.
// ─────────────────────────────────────────────────────────────────────────────

func GetDistrictAdminStateID(
	ctx context.Context,
	profileID string,
) (int, error) {

	var stateID int

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT state_id
		FROM district_admins
		WHERE profile_id = $1
		`,
		profileID,
	).Scan(&stateID)

	if err != nil {
		return 0, err
	}

	return stateID, nil
}



// ─────────────────────────────────────────────────────────────────────────────
// NEW — DeleteDistrictAdminRepository
// Deletes the profile row. Because district_admins has
// ON DELETE CASCADE → profiles, the district_admins row
// is removed automatically.
// ─────────────────────────────────────────────────────────────────────────────

func DeleteDistrictAdminRepository(
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

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			switch pgErr.Code {

			// invalid UUID format passed to Postgres
			case "22P02":
				return shared.ErrInvalidUUID

			// table missing (should never happen in prod)
			case "42P01":
				return errors.New("required database table does not exist")

			default:
				return errors.New("failed to delete district admin")
			}
		}

		return errors.New("database operation failed")
	}

	// 0 rows affected means no district admin with this profileID exists
	if commandTag.RowsAffected() == 0 {
		return shared.ErrDistrictAdminNotFound
	}

	return nil
}



func GetDistrictAdminDistrictID(
	ctx context.Context,
	profileID string,
) (int, error) {

	query := `
		SELECT district_id
		FROM district_admins
		WHERE profile_id = $1
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