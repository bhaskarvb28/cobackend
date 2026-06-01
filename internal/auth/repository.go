package auth

import (
	"context"
	"errors"
	"fmt"

	"cobackend/internal/db"
	"cobackend/internal/shared"

	"github.com/jackc/pgx/v5"
)


// GetUserByEmail returns authentication credentials
// and role information associated with an email address.
func GetUserByEmail(
	ctx context.Context,
	email string,
) (AuthUser, error) {

	var user AuthUser

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT 
			u.id,
			u.first_name,
			u.last_name,
			u.email,
			u.contact_number,
			u.password_hash,
			r.id,
			r.code,
			r.display_name
		FROM users u
		JOIN roles r 
			ON u.role_id = r.id
		WHERE u.email = $1
		`,
		email,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.ContactNumber,
		&user.PasswordHash,
		&user.Role.ID,
		&user.Role.Code,
		&user.Role.DisplayName,
	)

	if err != nil {
		return AuthUser{}, err
	}

	return user, nil
}

// ----------------------------------------------------------------------------------------------------------

func GetUserRoleCodeByID(
	ctx context.Context,
	userID string,
) (string, error) {

	query := `
		SELECT r.code
		FROM users u
		INNER JOIN roles r
			ON r.id = u.role_id
		WHERE u.id = $1
	`

	var roleCode string

	err := db.DB.QueryRow(
		ctx,
		query,
		userID,
	).Scan(&roleCode)

	if err != nil {

		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return "", shared.ErrUserNotFound
		}

		return "", err
	}

	return roleCode, nil
}

func CheckUserExistsByEmail(
	ctx context.Context,
	email string,
) (bool, error) {

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE email = $1
		)
	`

	var exists bool

	err := db.DB.QueryRow(
		ctx,
		query,
		email,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

// CreateUserTx creates a new user within
// an existing database transaction.
func CreateUserTx(
	ctx context.Context,
	tx pgx.Tx,
	input CreateUserInput,
) (*AuthUser, error) {

	var user AuthUser

	err := tx.QueryRow(
		ctx,
		`
		WITH inserted_user AS (

			INSERT INTO users (
				first_name,
				last_name,
				email,
				password_hash,
				contact_number,
				role_id
			)
			VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				(
					SELECT id
					FROM roles
					WHERE code = $6
				)
			)
			RETURNING
				id,
				first_name,
				last_name,
				email,
				contact_number,
				role_id
		)

		SELECT
			u.id,
			u.first_name,
			u.last_name,
			u.email,
			u.contact_number,
			r.id,
			r.code,
			r.display_name
		FROM inserted_user u
		JOIN roles r
			ON r.id = u.role_id
		`,
		input.FirstName,
		input.LastName,
		input.Email,
		input.PasswordHash,
		input.ContactNumber,
		input.Role,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.ContactNumber,
		&user.Role.ID,
		&user.Role.Code,
		&user.Role.DisplayName,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(
	ctx context.Context,
	userID string,
) (UserResponse, error) {

	var user UserResponse

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT
			u.id,
			u.first_name,
			u.last_name,
			u.email,
			u.contact_number,

			r.id,
			r.code,
			r.display_name

		FROM users u

		INNER JOIN roles r
			ON r.id = u.role_id

		WHERE u.id = $1
		`,
		userID,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.ContactNumber,

		&user.Role.ID,
		&user.Role.Code,
		&user.Role.DisplayName,
	)

	if err != nil {

		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return UserResponse{}, shared.ErrUserNotFound
		}

		return UserResponse{}, err
	}

	return user, nil
}

func GetProfileCompletedStatus(
	ctx context.Context,
	userID string,
	roleName string,
) (bool, error) {

	var query string

	switch roleName {

	case "super_admin":
		return true, nil

	case "state_admin":

		query = `
			SELECT profile_completed
			FROM state_admins
			WHERE user_id = $1
		`

	case "district_admin":

		query = `
			SELECT profile_completed
			FROM district_admins
			WHERE user_id = $1
		`

	case "district_coach":

		query = `
			SELECT profile_completed
			FROM district_coaches
			WHERE user_id = $1
		`

	case "academy_admin":

		query = `
			SELECT profile_completed
			FROM academy_admins
			WHERE user_id = $1
		`

	case "academy_coach":

		query = `
			SELECT profile_completed
			FROM academy_coaches
			WHERE user_id = $1
		`

	case "player":

		query = `
			SELECT profile_completed
			FROM players
			WHERE user_id = $1
		`

	default:
		return false, shared.ErrInvalidRole
	}

	var profileCompleted bool

	err := db.DB.QueryRow(
		ctx,
		query,
		userID,
	).Scan(&profileCompleted)

	if err != nil {
		fmt.Print("Profile Complete Middleware : ", err)
		return false, err
	}

	return profileCompleted, nil
}