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
			p.password_hash,
			p.role_id,
			r.name
		FROM profiles p
		JOIN roles r ON p.role_id = r.id
		WHERE p.email = $1
		`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.RoleID,
		&user.Role,
	)

	if err != nil {
		return AuthUser{}, err
	}

	return user, nil
}

