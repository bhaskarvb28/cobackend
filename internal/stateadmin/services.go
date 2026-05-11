package stateadmin

import (
	"context"
	"encoding/hex"

	// "golang.org/x/crypto/bcrypt"

	// "cobackend/internal/auth"
	"cobackend/internal/invitations"
	"cobackend/internal/mail"
	"cobackend/internal/profiles"
	"cobackend/internal/roles"
	"cobackend/internal/shared"
	"cobackend/internal/states"
	"cobackend/internal/validation"

	// "cobackend/internal/utils"

	"strings"

	"crypto/rand"

	"fmt"
	"net/http"

	"time"
	"os"
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
		return "", shared.NewAPIError(
			http.StatusBadRequest,
			"email is required",
		)
	}

	if !validation.IsValidEmail(email) {
		return "", shared.NewAPIError(
			http.StatusBadRequest,
			"invalid email format",
		)
	}

	if input.AssignedStateID == "" {
		return "", shared.NewAPIError(
			http.StatusBadRequest,
			"assigned_state_id is required",
		)
	}

	if !validation.IsValidUUID(input.AssignedStateID) {
		return "", shared.NewAPIError(
			http.StatusBadRequest,
			"invalid state ID",
		)
	}

	stateExists, err := states.CheckStateExists(ctx, input.AssignedStateID)

	if err != nil {
		return "", err
	}

	if !stateExists {
		return "", shared.NewAPIError(
			http.StatusNotFound,
			"assigned state does not exist",
		)
	}


	profileExists, err := profiles.CheckEmailExists(ctx, email)

	if err != nil {
		return "", err
	}

	if profileExists {
		return "", shared.NewAPIError(
			http.StatusConflict,
			"profile with this email already exists",
		)
	}

	pendingInviteExists, err := invitations.ExistsPendingInvitationByEmail(ctx, email)

	if err != nil {
		return "", err
	}

	if pendingInviteExists {
		return "", shared.NewAPIError(
			http.StatusConflict,
			"pending invitation already exists for this email",
		)
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

	frontendURL := os.Getenv("FRONTEND_URL")

	inviteLink := fmt.Sprintf(
		"%s/setup-account?token=%s",
		frontendURL, 
		token,
	)

	err = mail.SendStateAdminInvitationEmailBrevo(
		email,
		inviteLink,
	)

	if err != nil {

		fmt.Print("Brevo error")

		deleteErr := invitations.DeleteInvitationByToken(
			ctx,
			token,
		)

		if deleteErr != nil {
			return "", deleteErr
		}

		return "", err
	}

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