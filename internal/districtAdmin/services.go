package districtAdmin

// import (
// 	"context"
// 	"strings"

// 	"cobackend/internal/shared"
// 	"cobackend/internal/validation"

// 	"cobackend/internal/districts"

// 	"net/http"

// 	"cobackend/internal/invitations"
// 	"cobackend/internal/profile"
// 	"cobackend/internal/roles"
// 	"cobackend/internal/stateAdmin"

// 	"crypto/rand"
// 	"encoding/hex"
// 	"time"

// 	"os"

// 	"fmt"

// 	"cobackend/internal/mail"
// 	"errors"

// 	"github.com/jackc/pgx/v5"

// 	// "golang.org/x/crypto/bcrypt"
// )

// func InviteDistrictAdminService(
// 	ctx context.Context,
// 	input InviteDistrictAdminInput,
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

// 	districtExists, err := districts.CheckDistrictExists(ctx, input.DistrictID)

// 	if err != nil {
// 		return "", err
// 	}

// 	if !districtExists {
// 		return "", shared.NewAPIError(
// 			http.StatusNotFound,
// 			"assigned district does not exist",
// 		)
// 	}

// 	profileExists, err := profile.CheckEmailExists(ctx, email)

// 	if err != nil {
// 		return "", err
// 	}

// 	if profileExists {
// 		return "", shared.NewAPIError(
// 			http.StatusConflict,
// 			"profile with this email already exists",
// 		)
// 	}

// 	pendingInviteExists, err := invitations.ExistsPendingInvitationByEmail(ctx, email)

// 	if err != nil {
// 		return "", err
// 	}

// 	if pendingInviteExists {
// 		return "", shared.NewAPIError(
// 			http.StatusConflict,
// 			"pending invitation already exists for this email",
// 		)
// 	}

// 	roleID, err := roles.GetRoleIDByName(ctx, "district_admin")

// 	if err != nil {
// 		return "", err
// 	}

// 	tokenBytes := make([]byte, 32)

// 	_, err = rand.Read(tokenBytes)

// 	if err != nil {
// 		return "", err
// 	}

// 	token := hex.EncodeToString(tokenBytes)

// 	expiresAt := time.Now().Add(24 * time.Hour)

// 	StateID, err := districts.GetStateIDByDistrictID(
// 		ctx,
// 		input.DistrictID,
// 	)

// 	if err != nil {
// 		return "", err
// 	}

// 	err = invitations.CreateInvitationRepository(
// 		ctx,
// 		email,
// 		roleID,
// 		authUserID,
// 		token,
// 		&StateID,
// 		&input.DistrictID,
// 		nil,
// 		nil,
// 		nil,
// 		expiresAt,
// 	)

// 	if err != nil {
// 		return "", err
// 	}

// 	frontendURL := os.Getenv("FRONTEND_URL")

// 	inviteLink := fmt.Sprintf(
// 		"%s/setup-account?token=%s",
// 		frontendURL, 
// 		token,
// 	)

// 	err = mail.SendDistrictAdminInvitationEmailBrevo(
// 		email,
// 		inviteLink,
// 	)

// 	if err != nil {

// 		fmt.Print("Brevo error")

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


// func GetDistrictAdminsService(
// 	ctx context.Context,
// 	query GetDistrictAdminsQuery,
// ) (PaginatedDistrictAdmins, error) {

// 	return GetDistrictAdminsRepository(
// 		ctx,
// 		query,
// 	)
// }

// // func UpdateDistrictAdminService(ctx context.Context, id string, input UpdateDistrictAdminInput) error {
// // 	return UpdateDistrictAdminRepository(ctx, id, input)
// // }



// func UpdateDistrictAdminService(
// 	ctx context.Context,
// 	profileID string,
// 	input UpdateDistrictAdminInput,
// ) error {

// 	// check district admin exists
// 	exists, err := CheckDistrictAdminExists(
// 		ctx,
// 		profileID,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	if !exists {
// 		return shared.ErrDistrictAdminNotFound
// 	}

// 	// validate district exists
// 	districtExists, err := districts.CheckDistrictExists(
// 		ctx,
// 		*input.DistrictID,
// 	)

// 	if err != nil {
// 		return err
// 	}

// 	if !districtExists {
// 		return shared.ErrInvalidDistrict
// 	}

// 	return UpdateDistrictAdminRepository(
// 		ctx,
// 		profileID,
// 		input,
// 	)
// }

// func DeleteDistrictAdminService(
// 	ctx context.Context,
// 	authUserID string,
// 	profileID string,
// ) error {

// 	// Step 1: Find which state the district admin belongs to
// 	districtAdminStateID, err := GetDistrictAdminStateID(
// 		ctx,
// 		profileID,
// 	)

// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			return shared.ErrDistrictAdminNotFound
// 		}
// 		return shared.NewAPIError(
// 			http.StatusInternalServerError,
// 			"failed to fetch district admin",
// 		)
// 	}

// 	// Step 2: Find which state the logged-in state admin manages
// 	stateAdminStateID, err := stateAdmin.GetStateAdminStateID(
// 		ctx,
// 		authUserID,
// 	)

// 	if err != nil {
// 		return shared.NewAPIError(
// 			http.StatusInternalServerError,
// 			"failed to fetch state admin details",
// 		)
// 	}

// 	// Step 3: State admin can only delete district admins in their own state
// 	if stateAdminStateID != districtAdminStateID {
// 		return shared.NewAPIError(
// 			http.StatusForbidden,
// 			shared.ErrForbiddenDistrict.Error(),
// 		)
// 	}

// 	// Step 4: Delete the profile (district_admins row auto-removed via CASCADE)
// 	return DeleteDistrictAdminRepository(
// 		ctx,
// 		profileID,
// 	)
// }