package academyCoach

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"cobackend/internal/districts"
	"cobackend/internal/invitations"
	"cobackend/internal/mail"
	"cobackend/internal/profiles"
	"cobackend/internal/roles"
	"cobackend/internal/shared"
	"cobackend/internal/validation"
)

func InviteAcademyCoachService(
	ctx context.Context,
	input InviteAcademyCoachInput,
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

	roleID, err := roles.GetRoleIDByName(ctx, "academy_coach")

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

	StateID, err := districts.GetStateIDByDistrictID(
		ctx,
		input.DistrictID,
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
		&StateID,
		&input.DistrictID,
		&input.AcademyID,
		input.DisciplinesSpecialized,
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

	err = mail.SendAcademyCoachInvitationEmailBrevo(
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