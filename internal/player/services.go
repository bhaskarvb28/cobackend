package player

// import (
// 	"context"
// 	"crypto/rand"
// 	"encoding/hex"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"time"

// 	"cobackend/internal/academy"
// 	"cobackend/internal/districts"
// 	"cobackend/internal/invitations"
// 	"cobackend/internal/mail"
// 	"cobackend/internal/profile"
// 	"cobackend/internal/roles"
// 	"cobackend/internal/shared"
// 	"cobackend/internal/validation"
// 	"cobackend/internal/academyCoach"
// )

// func InvitePlayerService(
// 	ctx context.Context,
// 	input InvitePlayerInput,
// 	authUserID string,
// ) (string, error) {

// 	email := strings.ToLower(
// 		strings.TrimSpace(input.Email),
// 	)

// 	if email == "" {
// 		return "", shared.NewAPIError(
// 			http.StatusBadRequest,
// 			"email is required",
// 		)
// 	}

// 	if !validation.IsValidEmail(email) {
// 		return "", shared.NewAPIError(
// 			http.StatusBadRequest,
// 			"invalid email format",
// 		)
// 	}

// 	//------------------------------------------------
// 	// District Validation
// 	//------------------------------------------------

// 	districtExists, err := districts.CheckDistrictExists(
// 		ctx,
// 		input.DistrictID,
// 	)

// 	if err != nil {
// 		return "", err
// 	}

// 	if !districtExists {
// 		return "", shared.NewAPIError(
// 			http.StatusNotFound,
// 			"assigned district does not exist",
// 		)
// 	}

// 	//------------------------------------------------
// 	// Academy Validation
// 	//------------------------------------------------

// 	academyExists, err := academy.CheckAcademyExists(
// 		ctx,
// 		input.AcademyID,
// 	)

// 	if err != nil {
// 		return "", err
// 	}

// 	if !academyExists {
// 		return "", shared.NewAPIError(
// 			http.StatusNotFound,
// 			"assigned academy does not exist",
// 		)
// 	}

// 	//------------------------------------------------
// 	// Academy Belongs To District Validation
// 	//------------------------------------------------

// 	academyBelongsToDistrict, err := academy.CheckAcademyBelongsToDistrict(
// 		ctx,
// 		input.AcademyID,
// 		input.DistrictID,
// 	)

// 	if err != nil {
// 		return "", err
// 	}

// 	if !academyBelongsToDistrict {
// 		return "", shared.NewAPIError(
// 			http.StatusBadRequest,
// 			"academy does not belong to assigned district",
// 		)
// 	}

// 	//------------------------------------------------
// 	// Academy Coach Validation
// 	//------------------------------------------------

// 	academyCoachExists, err := academyCoach.CheckAcademyCoachExists(
// 		ctx,
// 		input.AcademyCoachID,
// 	)

// 	if err != nil {
// 		return "", err
// 	}

// 	if !academyCoachExists {
// 		return "", shared.NewAPIError(
// 			http.StatusNotFound,
// 			"assigned academy coach does not exist",
// 		)
// 	}

// 	//------------------------------------------------
// 	// Academy Coach Belongs To Academy Validation
// 	//------------------------------------------------

// 	coachBelongsToAcademy, err := academyCoach.CheckCoachBelongsToAcademy(
// 		ctx,
// 		input.AcademyCoachID,
// 		input.AcademyID,
// 	)

// 	if err != nil {
// 		return "", err
// 	}

// 	if !coachBelongsToAcademy {
// 		return "", shared.NewAPIError(
// 			http.StatusBadRequest,
// 			"academy coach does not belong to assigned academy",
// 		)
// 	}

// 	//------------------------------------------------
// 	// Existing Profile Validation
// 	//------------------------------------------------

// 	profileExists, err := profile.CheckEmailExists(
// 		ctx,
// 		email,
// 	)

// 	if err != nil {
// 		return "", err
// 	}

// 	if profileExists {
// 		return "", shared.NewAPIError(
// 			http.StatusConflict,
// 			"profile with this email already exists",
// 		)
// 	}

// 	//------------------------------------------------
// 	// Pending Invitation Validation
// 	//------------------------------------------------

// 	pendingInviteExists, err := invitations.ExistsPendingInvitationByEmail(
// 		ctx,
// 		email,
// 	)

// 	if err != nil {
// 		return "", err
// 	}

// 	if pendingInviteExists {
// 		return "", shared.NewAPIError(
// 			http.StatusConflict,
// 			"pending invitation already exists for this email",
// 		)
// 	}

// 	//------------------------------------------------
// 	// Get Player Role ID
// 	//------------------------------------------------

// 	roleID, err := roles.GetRoleIDByName(
// 		ctx,
// 		"player",
// 	)

// 	if err != nil {
// 		return "", err
// 	}

// 	//------------------------------------------------
// 	// Generate Token
// 	//------------------------------------------------

// 	tokenBytes := make([]byte, 32)

// 	_, err = rand.Read(tokenBytes)

// 	if err != nil {
// 		return "", err
// 	}

// 	token := hex.EncodeToString(tokenBytes)

// 	//------------------------------------------------
// 	// Expiry
// 	//------------------------------------------------

// 	expiresAt := time.Now().Add(24 * time.Hour)

// 	//------------------------------------------------
// 	// State ID From District
// 	//------------------------------------------------

// 	stateID, err := districts.GetStateIDByDistrictID(
// 		ctx,
// 		input.DistrictID,
// 	)

// 	if err != nil {
// 		return "", err
// 	}

// 	//------------------------------------------------
// 	// Create Invitation
// 	//------------------------------------------------

// 	err = invitations.CreateInvitationRepository(
// 		ctx,
// 		email,
// 		roleID,
// 		authUserID,
// 		token,
// 		&stateID,
// 		&input.DistrictID,
// 		&input.AcademyID,
// 		&input.AcademyCoachID,
// 		nil,
// 		expiresAt,
// 	)

// 	if err != nil {
// 		return "", err
// 	}

// 	//------------------------------------------------
// 	// Invitation Link
// 	//------------------------------------------------

// 	frontendURL := os.Getenv("FRONTEND_URL")

// 	inviteLink := fmt.Sprintf(
// 		"%s/setup-account?token=%s",
// 		frontendURL,
// 		token,
// 	)

// 	//------------------------------------------------
// 	// Send Email
// 	//------------------------------------------------

// 	err = mail.SendPlayerInvitationEmailBrevo(
// 		email,
// 		inviteLink,
// 	)

// 	if err != nil {

// 		deleteErr := invitations.DeleteInvitationByToken(
// 			ctx,
// 			token,
// 		)

// 		if deleteErr != nil {
// 			return "", deleteErr
// 		}

// 		return "", err
// 	}

// 	return inviteLink, nil
// }