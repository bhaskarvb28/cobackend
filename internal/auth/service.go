package auth

import (
	"context"
	"errors"

	"cobackend/internal/districtAdmin"
	"cobackend/internal/invitations"
	"cobackend/internal/jwt"
	"cobackend/internal/profiles"
	"cobackend/internal/roles"
	"cobackend/internal/shared"
	"cobackend/internal/stateAdmin"
	"cobackend/internal/utils"
	"cobackend/internal/validation"
	"cobackend/internal/districtCoach"
	"cobackend/internal/academyAdmin"

	"github.com/jackc/pgx/v5"

	"net/http"

	"cobackend/internal/db"

	"golang.org/x/crypto/bcrypt"

	"fmt"
	"strings"
	"time"
)

func Login(ctx context.Context, input LoginInput) (*LoginResponse, error) {
	// Find user by email
	user, err := GetUserByEmail(ctx, input.Email)
	if err != nil {
		fmt.Print(err)
		return nil, errors.New("invalid email or password")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(input.Password),
	)

	if err != nil {
		fmt.Print(err)
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT
	token, err := jwt.GenerateJWT(user.ID, user.RoleID, user.Role)
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

func AcceptInvitationService(
	ctx context.Context,
	input AcceptInvitationInput,
) error {

	token := strings.TrimSpace(input.Token)

	if token == "" {
		return shared.NewAPIError(
			http.StatusBadRequest,
			"token is required",
		)
	}

	if strings.TrimSpace(input.FirstName) == "" {
		return shared.NewAPIError(
			http.StatusBadRequest,
			"first_name is required",
		)
	}

	if strings.TrimSpace(input.Password) == "" {
		return shared.NewAPIError(
			http.StatusBadRequest,
			"password is required",
		)
	}

	invitation, err := invitations.GetInvitationByToken(
		ctx,
		token,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return shared.NewAPIError(
				http.StatusUnauthorized,
				"invalid invitation token",
			)
		}

		return shared.NewAPIError(
			http.StatusInternalServerError,
			"failed to fetch invitation",
		)
	}

	if invitation.Status == "accepted" {
		return shared.NewAPIError(
			http.StatusForbidden,
			"invitation already used",
		)
	}

	if time.Now().After(invitation.ExpiresAt) {
		return shared.NewAPIError(
			http.StatusForbidden,
			"invitation expired",
		)
	}


	if !validation.IsStrongPassword(input.Password) {
		return shared.NewAPIError(
			http.StatusBadRequest,
			shared.ErrWeakPassword.Error(),
		)
	}

	passwordHash, err := utils.HashPassword(
		input.Password,
	)

	if err != nil {
		return shared.NewAPIError(
			http.StatusInternalServerError,
			"failed to process password",
		)
	}

	if !validation.IsValidIndianPhone(input.ContactNumber) {
		return shared.NewAPIError(
			http.StatusBadRequest,
			shared.ErrInvalidPhoneNumber.Error(),
		)
	}

	tx, err := db.DB.Begin(ctx)

	if err != nil {
		return shared.NewAPIError(
			http.StatusInternalServerError,
			"failed to start transaction",
		)
	}

	defer tx.Rollback(ctx)

	profileID, err := profiles.CreateProfileTx(
		ctx,
		tx,
		profiles.CreateProfileInput{
			FirstName:     input.FirstName,
			LastName:      input.LastName,
			Email:         invitation.Email,
			PasswordHash:  passwordHash,
			ContactNumber: input.ContactNumber,
			RoleID:        invitation.RoleID,
		},
	)

	if err != nil {
		return shared.NewAPIError(
			http.StatusInternalServerError,
			"failed to create profile",
		)
	}

	roleName, err := roles.GetRoleNameByID(
		ctx,
		invitation.RoleID,
	)

	if err != nil {
		return shared.NewAPIError(
			http.StatusInternalServerError,
			"failed to fetch role",
		)
	}

	switch roleName {

	case "state_admin":

		err = stateAdmin.CreateStateAdminTx(
			ctx,
			tx,
			profileID,
			invitation.StateID,
		)

		if err != nil {
			return shared.NewAPIError(
				http.StatusInternalServerError,
				"failed to create state admin",
			)
		}

	case "district_admin":

		if !input.DPDPConsent {
			return shared.NewAPIError(
				http.StatusBadRequest,
				"dpdp consent is required",
			)
		}

		err = districtAdmin.CreateDistrictAdminTx(
			ctx,
			tx,
			profileID,
			invitation.StateID,
			invitation.DistrictID,
			input.DPDPConsent,
		)

		if err != nil {
			fmt.Print(err)
			return shared.NewAPIError(
				http.StatusInternalServerError,
				"failed to create district admin",
			)
		}

	case "district_coach":

		if !input.DPDPConsent {
			return shared.NewAPIError(
				http.StatusBadRequest,
				"dpdp consent is required",
			)
		}

		if input.CoachCode == "" {
			return shared.NewAPIError(
				http.StatusBadRequest,
				"coach code is required",
			)
		}

		if input.CoachingCertificateProof == "" {
			return shared.NewAPIError(
				http.StatusBadRequest,
				"coaching certificate proof is required",
			)
		}

		err = districtCoach.CreateDistrictCoachTx(
			ctx,
			tx,
			profileID,
			invitation.DistrictID,
			input.CoachCode,
			input.CoachingCertificateProof,
			input.DPDPConsent,
		)

	case "academy_admin":
		if !input.DPDPConsent {
			return shared.NewAPIError(
				http.StatusBadRequest,
				"dpdp consent is required",
			)
		}

		if input.RegistrationProof == "" {
			return shared.NewAPIError(
				http.StatusBadRequest,
				"registration proof is required",
			)
		}

		err = academyAdmin.CreateAcademyAdminTx(
			ctx,
			tx,
			profileID,
			invitation.AcademyID,
			input.GSTIN,
			input.RegistrationProof,
			input.DPDPConsent,
		)

	}

	err = invitations.MarkInvitationUsedTx(
		ctx,
		tx,
		invitation.ID,
	)

	if err != nil {
		return shared.NewAPIError(
			http.StatusInternalServerError,
			"failed to update invitation",
		)
	}

	err = tx.Commit(ctx)

	if err != nil {
		return shared.NewAPIError(
			http.StatusInternalServerError,
			"failed to commit transaction",
		)
	}

	return nil
}