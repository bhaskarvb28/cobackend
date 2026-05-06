package auth

import (
	"cobackend/internal/db"
	"context"
)

// func CreateUser(ctx context.Context, input RegisterInput, hashed string) error {
// 	_, err := db.DB.Exec(ctx,
// 		`INSERT INTO profiles 
// 		(id, first_name, last_name, email, password, contact_number)
// 		VALUES ($1,$2,$3,$4,$5,$6)`,
// 		uuid.New(),
// 		input.FirstName,
// 		input.LastName,
// 		input.Email,
// 		hashed,
// 		input.ContactNumber,
// 	)

// 	return err
// }

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