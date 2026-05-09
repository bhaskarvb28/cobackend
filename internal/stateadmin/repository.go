package stateadmin

import (
	"context"
	"errors"

	"cobackend/internal/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"strconv"

	"cobackend/internal/shared"
)

// Super Admin Can Create a State Admin
func CreateStateAdminRepository(
	ctx context.Context,
	input CreateStateAdminInput,
	hashedPassword string,
) error {

	tx, err := db.DB.Begin(ctx)
	if err != nil {
		return errors.New("failed to start database transaction")
	}

	defer tx.Rollback(ctx)

	var stateAdminRoleID string

	err = tx.QueryRow(
		ctx,
		`
		SELECT role_id
		FROM roles
		WHERE role_name = 'state_admin'
		`,
	).Scan(&stateAdminRoleID)

	if err != nil {
		return errors.New("failed to fetch state admin role")
	}

	profileID := uuid.New()

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO profiles (
			id,
			first_name,
			last_name,
			email,
			password,
			contact_number,
			role_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		`,
		profileID,
		input.FirstName,
		input.LastName,
		input.Email,
		hashedPassword,
		input.ContactNumber,
		stateAdminRoleID,
	)

	if err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			switch pgErr.Code {


			// UNIQUE VIOLATION
			case "23505":

				switch pgErr.ConstraintName {

				case "profiles_email_key":
					return errors.New("email already exists")

				default:
					return errors.New("duplicate value already exists")
				}

			// FOREIGN KEY VIOLATION
			case "23503":
				return errors.New("invalid foreign key reference")

			// NOT NULL VIOLATION
			case "23502":
				return errors.New("required field is missing")

			// UNDEFINED TABLE
			case "42P01":
				return errors.New("required database table does not exist")

			// UNDEFINED COLUMN
			case "42703":
				return errors.New("required database column does not exist")

			default:
				return errors.New("database operation failed")
			}
		}

		return errors.New("failed to create profile")
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO state_admins (
			profile_id,
			assigned_state
		)
		VALUES ($1, $2)
		`,
		profileID,
		input.AssignedState,
	)

	if err != nil {

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {

			switch pgErr.Code {

			// FOREIGN KEY VIOLATION
			case "23503":
				return errors.New("invalid assigned state")

			// UNIQUE VIOLATION
			case "23505":
				return errors.New("state admin already exists")

			// NOT NULL VIOLATION
			case "23502":
				return errors.New("assigned state is required")

			// UNDEFINED TABLE
			case "42P01":
				return errors.New("required database table does not exist")

			default:
				return errors.New("failed to assign state admin")
			}
		}

		return errors.New("failed to create state admin")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.New("failed to commit database transaction")
	}

	return nil
}

// Super Admin can get state admins
func GetStateAdminsRepository(
	ctx context.Context,
	query GetStateAdminsQuery,
) ([]StateAdminResponse, error) {

	offset := (query.Page - 1) * query.Limit

	baseQuery := `
	SELECT 
		p.id,
		p.first_name,
		p.last_name,
		p.email,
		p.contact_number,
		sa.assigned_state
	FROM profiles p
	INNER JOIN state_admins sa
		ON p.id = sa.profile_id
	WHERE 1=1
	`

	args := []interface{}{}
	argPos := 1

	// SEARCH
	if query.Search != "" {

		baseQuery += `
		AND (
			p.first_name ILIKE $` + strconv.Itoa(argPos) + `
			OR p.last_name ILIKE $` + strconv.Itoa(argPos) + `
			OR p.email ILIKE $` + strconv.Itoa(argPos) + `
		)
		`

		args = append(
			args,
			"%"+query.Search+"%",
		)

		argPos++
	}

	// STATE FILTER
	if query.AssignedState != 0 {

		baseQuery += `
		AND sa.assigned_state = $` + strconv.Itoa(argPos)

		args = append(
			args,
			query.AssignedState,
		)

		argPos++
	}

	// ORDER + PAGINATION
	baseQuery += `
	ORDER BY p.first_name ASC
	LIMIT $` + strconv.Itoa(argPos) + `
	OFFSET $` + strconv.Itoa(argPos+1)

	args = append(
		args,
		query.Limit,
		offset,
	)

	rows, err := db.DB.Query(
		ctx,
		baseQuery,
		args...,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	stateAdmins := []StateAdminResponse{}

	for rows.Next() {

		var stateAdmin StateAdminResponse

		err := rows.Scan(
			&stateAdmin.ID,
			&stateAdmin.FirstName,
			&stateAdmin.LastName,
			&stateAdmin.Email,
			&stateAdmin.ContactNumber,
			&stateAdmin.AssignedState,
		)

		if err != nil {
			return nil, err
		}

		stateAdmins = append(
			stateAdmins,
			stateAdmin,
		)
	}

	// check iteration errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stateAdmins, nil
}

func CheckStateAdminExists(
	ctx context.Context,
	id string,
) (bool, error) {

	var exists bool

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT EXISTS(
			SELECT 1
			FROM state_admins
			WHERE profile_id = $1
		)
		`,
		id,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func UpdateAssignedStateRepository(
	ctx context.Context,
	id string,
	input UpdateAssignedStateInput,
) error {

	commandTag, err := db.DB.Exec(
		ctx,
		`
		UPDATE state_admins
		SET assigned_state = $1
		WHERE profile_id = $2
		`,
		input.AssignedState,
		id,
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
	id string,
) error {

	commandTag, err := db.DB.Exec(
		ctx,
		`
		DELETE FROM profiles
		WHERE id = $1
		`,
		id,
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