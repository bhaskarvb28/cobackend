package auth

import (
	"cobackend/internal/db"
	"context"
)

func GetUserByEmail(ctx context.Context, email string) (AuthUser, error) {
	var user AuthUser
    err := db.DB.QueryRow(ctx,
		`SELECT p.id, p.email, p.password, p.role_id, r.role_name
		 FROM profiles p
		 JOIN roles r ON p.role_id = r.role_id
		 WHERE p.email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.RoleID,
		&user.Role, // <-- store role name here
	)

	if err != nil {
		return AuthUser{}, err
	}

	return user, nil
}