package districtadmin

import (
	"context"
	"strings"

	"cobackend/internal/shared"
	"cobackend/internal/validation"

	"cobackend/internal/districts"

	"net/http"

	"cobackend/internal/invitations"
	"cobackend/internal/profiles"
	"cobackend/internal/roles"

	"crypto/rand"
	"encoding/hex"
	"time"

	"os"

	"fmt"

	"cobackend/internal/mail"

	// "golang.org/x/crypto/bcrypt"
)

func InviteDistrictAdminService(
	ctx context.Context,
	input InviteDistrictAdminInput,
	authUserID string,
) (string, error) {

	email := strings.ToLower(
		strings.TrimSpace(input.Email),
	)

	assignedDistrictID := strings.TrimSpace(input.AssignedDistrictID)

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

	districtExists, err := districts.CheckDistrictExists(ctx, assignedDistrictID)

	if err != nil {
		return "", err
	}

	if !districtExists {
		return "", shared.NewAPIError(
			http.StatusNotFound,
			"assigned district does not exist",
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

	roleID, err := roles.GetRoleIDByName(ctx, "district_admin")

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

	assignedStateID, err := districts.GetStateIDByDistrictID(
		ctx,
		assignedDistrictID,
	)

	if err != nil {
		return "", err
	}

	err = invitations.CreateInvitationRepository(
		ctx,
		email,
		roleID,
		authUserID,
		token,
		&assignedStateID,
		&assignedDistrictID,
		nil,
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

	err = mail.SendDistrictAdminInvitationEmailBrevo(
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

// func CreateDistrictAdminService(ctx context.Context, input CreateDistrictAdminInput) error {
// 	hashedPassword, err := bcrypt.GenerateFromPassword(
// 		[]byte(input.Password),
// 		bcrypt.DefaultCost,
// 	)
// 	if err != nil {
// 		return err
// 	}

// 	return CreateDistrictAdminRepository(ctx, input, string(hashedPassword))
// }

// func GetDistrictAdminsService(ctx context.Context, query GetDistrictAdminsQuery) ([]DistrictAdmin, error) {
// 	return GetDistrictAdminsRepository(ctx, query)
// }

// func UpdateDistrictAdminService(ctx context.Context, id string, input UpdateDistrictAdminInput) error {
// 	return UpdateDistrictAdminRepository(ctx, id, input)
// }

// func DeleteDistrictAdminService(ctx context.Context, id string) error {
// 	return DeleteDistrictAdminRepository(ctx, id)
// }
