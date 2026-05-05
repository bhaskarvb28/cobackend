package auth

import (
	"cobackend/internal/db"
	"context"

	"github.com/google/uuid"
)

func CreateUser(ctx context.Context, input RegisterInput, hashed string) error {
	_, err := db.DB.Exec(ctx,
		`INSERT INTO profiles 
		(id, first_name, last_name, email, password, role_id, contact_number)
		VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		uuid.New(),
		input.FirstName,
		input.LastName,
		input.Email,
		hashed,
		input.RoleID,
		input.ContactNumber,
	)

	return err
}