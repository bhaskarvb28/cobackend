package stateAdmin

import (
	"context"
	"errors"
	"math"

	"cobackend/internal/db"

	// "github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5"


	"strconv"

	"cobackend/internal/shared"
)

func CreateStateAdminTx(
	ctx context.Context,
	tx pgx.Tx,
	profileID string,
	StateID *int,
) error {

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO state_admins (
			profile_id,
			state_id
		)
		VALUES (
			$1,
			$2
		)
		`,
		profileID,
		StateID,
	)

	if err != nil {
		return err
	}

	return nil
}

func CheckStateAdminExists(
	ctx context.Context,
	profileID string,
) (bool, error) {

	var exists bool

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT EXISTS (
			SELECT 1
			FROM state_admins
			WHERE profile_id = $1
		)
		`,
		profileID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func UpdateAssignedStateRepository(
	ctx context.Context,
	profileID string,
	input UpdateStateInput,
) error {

	commandTag, err := db.DB.Exec(
		ctx,
		`
		UPDATE state_admins
		SET state_id = $1
		WHERE profile_id = $2
		`,
		input.StateID,
		profileID,
	)

	if err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			switch pgErr.Code {

			case "23503":
				return shared.ErrInvalidState

			case "42P01":
				return errors.New(
					"required database table does not exist",
				)

			default:
				return errors.New(
					"failed to update assigned state",
				)
			}
		}

		return errors.New(
			"database operation failed",
		)
	}

	// no rows updated
	if commandTag.RowsAffected() == 0 {
		return shared.ErrStateAdminNotFound
	}

	return nil
}

func DeleteStateAdminRepository(
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

			// invalid uuid format
			case "22P02":
				return shared.ErrInvalidUUID

			// undefined table
			case "42P01":
				return errors.New(
					"required database table does not exist",
				)

			default:
				return errors.New(
					"failed to delete state admin",
				)
			}
		}

		return errors.New(
			"database operation failed",
		)
	}

	// no rows deleted
	if commandTag.RowsAffected() == 0 {
		return shared.ErrStateAdminNotFound
	}

	return nil
}

func GetStateAdminsRepository(
	ctx context.Context,
	query GetStateAdminsQuery,
) (PaginatedStateAdmins, error) {

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
		sa.state_id
	FROM profiles p
	INNER JOIN state_admins sa
		ON p.id = sa.profile_id
	LEFT JOIN states s
		ON s.id = sa.state_id
	WHERE 1=1
	`

	countQuery := `
	SELECT COUNT(*)
	FROM profiles p
	INNER JOIN state_admins sa
		ON p.id = sa.profile_id
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
		AND sa.state_id = $` + strconv.Itoa(argPos)

		baseQuery += condition
		countQuery += condition

		args = append(args, query.StateID)
		argPos++
	}

	var total int

	err := db.DB.QueryRow(
		ctx,
		countQuery,
		args...,
	).Scan(&total)

	if err != nil {
		return PaginatedStateAdmins{}, err
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
		return PaginatedStateAdmins{}, err
	}

	defer rows.Close()

	admins := []StateAdmin{}

	for rows.Next() {

		var a StateAdmin

		err := rows.Scan(
			&a.ID,
			&a.FirstName,
			&a.LastName,
			&a.Email,
			&a.ContactNumber,
			&a.StateID,
		)

		if err != nil {
			return PaginatedStateAdmins{}, err
		}

		admins = append(admins, a)
	}

	if err := rows.Err(); err != nil {
		return PaginatedStateAdmins{}, err
	}

	totalPages := 1

	if query.Limit > 0 {

		totalPages = int(math.Ceil(
			float64(total) / float64(query.Limit),
		))
	}

	return PaginatedStateAdmins{
		Items:       admins,
		Page:        query.Page,
		Limit:       query.Limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     query.Page < totalPages,
		HasPrevious: query.Page > 1,
	}, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// NEW — GetStateAdminStateID
// Fetches which state the logged-in state admin manages.
// e.g. Karnataka state admin → returns Karnataka's state_id
// ─────────────────────────────────────────────────────────────────────────────

func GetStateAdminStateID(
	ctx context.Context,
	profileID string,
) (int, error) {

	var stateID int

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT state_id
		FROM state_admins
		WHERE profile_id = $1
		`,
		profileID,
	).Scan(&stateID)

	if err != nil {
		return 0, err
	}

	return stateID, nil
}


func GetAssignedStateByStateAdmin(
	ctx context.Context,
	profileID string,
) (int, error) {

	var assignedStateID int

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT state_id
		FROM state_admins
		WHERE profile_id = $1
		`,
		profileID,
	).Scan(&assignedStateID)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return 0, shared.ErrStateAdminNotFound
		}

		return 0, err
	}

	return assignedStateID, nil
}