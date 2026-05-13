package districtCoach

import (
	"strings"
	"cobackend/internal/shared"
	"net/http"

	"context"

	"cobackend/internal/validation"

	"cobackend/internal/districts"

	"cobackend/internal/profiles"

	"cobackend/internal/invitations"

	"cobackend/internal/roles"

	"crypto/rand"

	"encoding/hex"

	"time"

	"os"

	"fmt"

	"cobackend/internal/mail"
)

func InviteDistrictCoachService(
	ctx context.Context,
	input InviteDistrictCoachInput,
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

	districtExists, err := districts.CheckDistrictExists(ctx, input.DistrictID)

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

	roleID, err := roles.GetRoleIDByName(ctx, "district_coach")

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

	if err != nil {
		return "", err
	}

	err = invitations.CreateInvitationRepository(
		ctx,
		email,
		roleID,
		authUserID,
		token,
		&input.StateID,
		&input.DistrictID,
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

	err = mail.SendDistrictCoachInvitationEmailBrevo(
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

func GetDistrictCoachesService(
	ctx context.Context,
	query GetDistrictCoachesQuery,
) (PaginatedDistrictCoaches, error) {

	return GetDistrictCoachesRepository(
		ctx,
		query,
	)
}