package auth

import (
	"context"

	"cobackend/internal/db"
)

func GetUserByEmail(ctx context.Context, email string) (AuthUser, error) {
	var user AuthUser

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT 
			p.id,
			p.email,
			p.password,
			p.role_id,
			r.role_name
		FROM profiles p
		JOIN roles r ON p.role_id = r.role_id
		WHERE p.email = $1
		`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.RoleID,
		&user.Role,
	)

	if err != nil {
		return AuthUser{}, err
	}

	return user, nil
}

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