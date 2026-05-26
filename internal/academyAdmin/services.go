package academyAdmin

// import (
// 	"context"
// 	"crypto/rand"
// 	"encoding/hex"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"time"

// 	"cobackend/internal/districts"
// 	"cobackend/internal/invitations"
// 	"cobackend/internal/mail"
// 	"cobackend/internal/profile"
// 	"cobackend/internal/roles"
// 	"cobackend/internal/shared"
// 	"cobackend/internal/validation"
// 	"cobackend/internal/academy"
// 	"cobackend/internal/districtAdmin"
// )

// func InviteAcademyAdminService(
// 	ctx context.Context,
// 	input InviteAcademyAdminInput,
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
// 	// Get Role ID
// 	//------------------------------------------------

// 	roleID, err := roles.GetRoleIDByName(
// 		ctx,
// 		"academy_admin",
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
// 	// Get State ID From District
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
// 		nil,
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

// 	err = mail.SendAcademyAdminInvitationEmailBrevo(
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

// func GetAcademyAdminsService(
// 	ctx context.Context,
// 	query GetAcademyAdminsQuery,
// ) (PaginatedAcademyAdmins, error) {

// 	return GetAcademyAdminsRepository(ctx, query)
// }

// func GetAcademyAdminByIDService(
// 	ctx context.Context,
// 	profileID string,
// ) (AcademyAdmin, error) {

// 	return GetAcademyAdminByIDRepository(ctx, profileID)
// }

// func UpdateAcademyAdminService(
// 	ctx context.Context,
// 	profileID string,
// 	input UpdateAcademyAdminInput,
// ) error {

// 	// check academy admin exists
// 	exists, err := CheckAcademyAdminExists(
// 		ctx,
// 		profileID,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	if !exists {
// 		return shared.ErrAcademyAdminNotFound
// 	}

// 	// validate academy exists
// 	if input.AcademyID != nil {

// 		academyExists, err := academy.CheckAcademyExists(
// 			ctx,
// 			*input.AcademyID,
// 		)

// 		if err != nil {
// 			return err
// 		}

// 		if !academyExists {
// 			return shared.ErrInvalidAcademy
// 		}
// 	}

// 	return UpdateAcademyAdminRepository(
// 		ctx,
// 		profileID,
// 		input,
// 	)
// }

// func DeleteAcademyAdminService(
// 	ctx context.Context,
// 	authUserID string,
// 	academyAdminProfileID string,
// ) error {

// 	exists, err := CheckAcademyAdminExists(
// 		ctx,
// 		academyAdminProfileID,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	if !exists {
// 		return shared.ErrAcademyAdminNotFound
// 	}

// 	authDistrictID, err := districtAdmin.GetDistrictAdminDistrictID(
// 		ctx,
// 		authUserID,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	academyAdminDistrictID, err := GetAcademyAdminDistrictID(
// 		ctx,
// 		academyAdminProfileID,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	if authDistrictID != academyAdminDistrictID {
// 		return shared.ErrForbidden
// 	}

// 	return DeleteAcademyAdminRepository(
// 		ctx,
// 		academyAdminProfileID,
// 	)
// }

