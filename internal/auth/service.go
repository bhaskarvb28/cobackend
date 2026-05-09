package auth

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func Login(ctx context.Context, input LoginInput) (*LoginResponse, error) {
	// Find user by email
	user, err := GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	)

	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT
	token, err := GenerateJWT(user.ID, user.RoleID, user.Role)
	if err != nil {
		return nil, err
	}

	// Build response
	response := &LoginResponse{
		Token: token,
		User: UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}

	return response, nil
}