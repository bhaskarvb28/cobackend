package auth

import (
	"golang.org/x/crypto/bcrypt"
	"context"

	"errors"
)

// func Register(ctx context.Context, input RegisterInput) error {
// 	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}

// 	// Pass to repository
// 	return CreateUser(ctx, input, string(hashed))
// }

func Login(ctx context.Context, input LoginInput) (string, error) {
	// find user by email 
	user, err := GetUserByEmail(ctx, input.Email)
	if err != nil {
		return "", errors.New("Invalid email or password")
	}

	// compare passwords
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	)

	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// generate jwt
	token, err := GenerateJWT(user.ID, user.RoleID)
	if err != nil {
		return "", err
	}

	return token, nil


}