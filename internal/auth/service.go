package auth

import (
	"golang.org/x/crypto/bcrypt"
	"context"
)

func Register(ctx context.Context, input RegisterInput) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Pass to repository
	return CreateUser(ctx, input, string(hashed))
}