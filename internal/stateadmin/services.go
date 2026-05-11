package stateadmin

import (
	"context"
	"encoding/hex"

	// "golang.org/x/crypto/bcrypt"

	// "cobackend/internal/auth"
	// "cobackend/internal/shared"
	"cobackend/internal/invitations"
	// "cobackend/internal/mail"
	"cobackend/internal/profiles"
	"cobackend/internal/roles"
	"cobackend/internal/shared"
	"cobackend/internal/states"
	"cobackend/internal/validation"

	// "cobackend/internal/utils"

	"errors"
	"strings"

	"crypto/rand"

	"fmt"

	"time"
	// "net/http"
)

// func CreateStateAdminService(
// 	ctx context.Context,
// 	input CreateStateAdminInput,
// ) error {

// 	// check if email already exists
// 	exists, err := auth.CheckEmailExists(ctx, input.Email)
// 	if err != nil {
// 		return err
// 	}

// 	if exists {
// 		return shared.ErrEmailExists
// 	}

// 	// validate assigned state exists
// 	stateExists, err := states.CheckStateExists(
// 		ctx,
// 		input.AssignedState,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	if !stateExists {
// 		return shared.ErrInvalidState
// 	}

// 	// hash password
// 	hashedPassword, err := bcrypt.GenerateFromPassword(
// 		[]byte(input.Password),
// 		bcrypt.DefaultCost,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	return CreateStateAdminRepository(
// 		ctx,
// 		input,
// 		string(hashedPassword),
// 	)
// }

func InviteStateAdminService(
	ctx context.Context,
	input InviteStateAdminInput,
	authUserID string,
) (string, error) {

	email := strings.ToLower(
		strings.TrimSpace(input.Email),
	)

	if email == "" {
		return "", errors.New("email is required")
	}

	if !validation.IsValidEmail(input.Email) {
		return "", shared.ErrInvalidEmailFormat
	}

	if input.AssignedStateID == "" {
		return "", errors.New("assigned_state_id is required")
	}

	if !validation.IsValidUUID(input.AssignedStateID) {
		return "", errors.New("Invalid state ID")
	}

	stateExists, err := states.CheckStateExists(ctx, input.AssignedStateID)

	if err != nil {
		return "", err
	}

	if !stateExists {
		return "", errors.New("assigned state does not exist")
	}


	profileExists, err := profiles.CheckEmailExists(ctx, email)

	if err != nil {
		return "", err
	}

	if profileExists {
		return "", errors.New("profile with this email already exists")
	}

	pendingInviteExists, err := invitations.ExistsPendingInvitationByEmail(ctx, email)

	if err != nil {
		return "", err
	}

	if pendingInviteExists {
		return "", errors.New("pending invitation already exists for this email")
	}

	roleID, err := roles.GetRoleIDByName(ctx, "state_admin")

	if err != nil {
		return "", err
	}

	tokenBytes := make([]byte, 32)

	_, err = rand.Read(tokenBytes)

	if err != nil {
		return "", err
	}

	token := hex.EncodeToString(tokenBytes)

	expiresAt := time.Now().Add(24 * time.Hour)

	err = invitations.CreateInvitationRepository(
		ctx,
		email,
		roleID,
		authUserID,
		token,
		input.AssignedStateID,
		expiresAt,
	)

	if err != nil {
		return "", err
	}

	inviteLink := fmt.Sprintf(
		"http://localhost:3000/setup-account?token=%s",
		token,
	)

	// Uncomment this after choosing proper email service to send emails
	// err = mail.SendStateAdminInvitationEmail(
	// 	email,
	// 	inviteLink,
	// )

	// if err != nil {

	// 	deleteErr := invitations.DeleteInvitationByToken(
	// 		ctx,
	// 		token,
	// 	)

	// 	if deleteErr != nil {
	// 		return "", deleteErr
	// 	}

	// 	return "", err
	// }

	return inviteLink, nil
}

// func GetStateAdminsService(ctx context.Context, query GetStateAdminsQuery) ([]StateAdminResponse, error) {
// 	return GetStateAdminsRepository(ctx, query)
// }


// func UpdateAssignedStateService(
// 	ctx context.Context,
// 	id string,
// 	input UpdateAssignedStateInput,
// ) error {

// 	// check state admin exists
// 	exists, err := CheckStateAdminExists(
// 		ctx,
// 		id,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	if !exists {
// 		return shared.ErrStateAdminNotFound
// 	}

// 	// validate state exists
// 	stateExists, err := states.CheckStateExists(
// 		ctx,
// 		input.AssignedState,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	if !stateExists {
// 		return shared.ErrInvalidState
// 	}

// 	return UpdateAssignedStateRepository(
// 		ctx,
// 		id,
// 		input,
// 	)
// }


// func DeleteStateAdminService(
// 	ctx context.Context,
// 	id string,
// ) error {

// 	// check state admin exists
// 	exists, err := CheckStateAdminExists(
// 		ctx,
// 		id,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	if !exists {
// 		return shared.ErrStateAdminNotFound
// 	}

// 	return DeleteStateAdminRepository(
// 		ctx,
// 		id,
// 	)
// }