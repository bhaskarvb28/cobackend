package profiles

import (
	"cobackend/internal/db"
	"context"

	"github.com/jackc/pgx/v5"
)

func CheckEmailExists(
	ctx context.Context,
	email string,
) (bool, error) {

	var exists bool

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT EXISTS (
			SELECT 1
			FROM profiles
			WHERE email = $1
		)
		`,
		email,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func CreateProfileTx(
	ctx context.Context,
	tx pgx.Tx,
	input CreateProfileInput,
) (string, error) {

	var profileID string

	err := tx.QueryRow(
		ctx,
		`
		INSERT INTO profiles (
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
			$6
		)
		RETURNING id
		`,
		input.FirstName,
		input.LastName,
		input.Email,
		input.PasswordHash,
		input.ContactNumber,
		input.RoleID,
	).Scan(&profileID)

	if err != nil {
		return "", err
	}

	return profileID, nil
}