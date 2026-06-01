package invitation

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	// "crypto/rand"
	// "crypto/sha256"
	// "encoding/hex"
	// "fmt"
	// "os"
	// "strconv"
	// "strings"
	// "time"

	// "cobackend/internal/academy"
	// "cobackend/internal/auth"
	// "cobackend/internal/db"
	// "cobackend/internal/districtAdmin"
	// "cobackend/internal/districtCoach"
	// "cobackend/internal/district"
	// "cobackend/internal/jwtToken"
	// "cobackend/internal/mail"
	"cobackend/internal/academy"
	"cobackend/internal/academyAdmin"
	"cobackend/internal/academyCoach"
	"cobackend/internal/auth"
	"cobackend/internal/db"
	"cobackend/internal/district"
	"cobackend/internal/districtAdmin"
	"cobackend/internal/districtCoach"
	"cobackend/internal/jwtToken"
	"cobackend/internal/mail"
	"cobackend/internal/player"
	"cobackend/internal/role"
	"cobackend/internal/shared"
	"cobackend/internal/state"
	"cobackend/internal/stateAdmin"
	"cobackend/internal/utils"
	"cobackend/internal/validation"
	// "cobackend/internal/stateAdmin"
	// "cobackend/internal/state"
	// "cobackend/internal/utils"
	// "cobackend/internal/validation"
	// "cobackend/internal/academyAdmin"
	// "cobackend/internal/academyCoach"
	// "cobackend/internal/player"
)

// GetInvitationsService fetches all invitations
// visible to the authenticated user.
func GetInvitationsService(
	ctx context.Context,
	userID string,
	role string,
	query InvitationsQueryParams,
) (*PaginatedInvitations, error) {

	// ----------------------------------------------------------
	// Invitation Permissions
	// ----------------------------------------------------------


	roles, ok := InvitationPermissions[role]

	if !ok {
		return nil, shared.ErrForbidden
	}

	isSuperAdmin := role == "super_admin"


	// ----------------------------------------------------------
	// Validate Sorting
	// ----------------------------------------------------------

	allowedSortFields := map[string]bool{
		"created_at": true,
		"expires_at": true,
		"email":      true,
		"status":     true,
	}

	if !allowedSortFields[query.SortBy] {
		query.SortBy = "created_at"
	}

	if query.Order != "asc" &&
		query.Order != "desc" {

		query.Order = "desc"
	}

	// ----------------------------------------------------------
	// Fetch Invitations
	// ----------------------------------------------------------

	invitations, err := GetInvitationsRepository(
		ctx,
		userID,
		roles,
		isSuperAdmin,
		query,
	)

	if err != nil {
		return nil, err
	}

	return invitations, nil
}


// CreateInvitationService creates a new invitation
// based on the provided input.
func CreateInvitationService(
	ctx context.Context,
	input CreateInvitationInput,
	authUserID string,
) error {

	//------------------------------------------------
	// Normalize Input
	//------------------------------------------------

	name := strings.TrimSpace(
		input.Name,
	)

	email := strings.ToLower(
		strings.TrimSpace(input.Email),
	)

	roleCode := strings.TrimSpace(
		input.Role,
	)

	scopeType := strings.TrimSpace(
		input.ScopeType,
	)

	scopeID := strings.TrimSpace(
		input.ScopeID,
	)

	//------------------------------------------------
	// Get Inviter Role
	//------------------------------------------------

	inviterRole, err := auth.GetUserRoleCodeByID(
		ctx,
		authUserID,
	)

	if err != nil {
		return err
	}

	//------------------------------------------------
	// Invitation Permissions
	//------------------------------------------------

	targetRoles, exists := InvitationPermissions[inviterRole]

	if !exists {
		return shared.ErrForbidden
	}

	canInvite := false

	for _, allowedRole := range targetRoles {

		if allowedRole == roleCode {
			canInvite = true
			break
		}
	}

	if !canInvite {
		return shared.ErrRoleNotAuthorized
	}

	//------------------------------------------------
	// Scope Validation
	//------------------------------------------------

	switch roleCode {

		//------------------------------------------------
		// State Admin
		//------------------------------------------------

	case "state_admin":

		if scopeType != "state" {
			return shared.ErrInvalidScope
		}

		stateID, err := strconv.Atoi(scopeID)

		if err != nil {
			return shared.ErrInvalidScope
		}

		exists, err := state.CheckStateExists(
			ctx,
			stateID,
		)

		if err != nil {
			return err
		}

		if !exists {
			return shared.ErrStateNotFound
		}

	//------------------------------------------------
	// District Roles
	//------------------------------------------------

	case "district_admin", "district_coach":

		if scopeType != "district" {
			return shared.ErrInvalidScope
		}

		districtID, err := strconv.Atoi(scopeID)

		if err != nil {
			return shared.ErrInvalidScope
		}

		exists, err := district.CheckDistrictExists(
			ctx,
			districtID,
		)

		if err != nil {
			return err
		}

		if !exists {
			return shared.ErrDistrictNotFound
		}

	//------------------------------------------------
	// Academy Roles
	//------------------------------------------------

	case "academy_admin", "academy_coach", "player":

		if scopeType != "academy" {
			return shared.ErrInvalidScope
		}

		if !validation.IsValidUUID(scopeID) {
			return shared.ErrInvalidScope
		}

		exists, err := academy.CheckAcademyExists(
			ctx,
			scopeID,
		)

		if err != nil {
			return err
		}

		if !exists {
			return shared.ErrAcademyNotFound
		}
	}

	//------------------------------------------------
	// Existing User Validation
	//------------------------------------------------

	emailExists, err := auth.CheckUserExistsByEmail(
		ctx,
		email,
	)

	if err != nil {
		return err
	}

	if emailExists {
		return shared.ErrEmailAlreadyExists
	}

	//------------------------------------------------
	// Pending Invitation Validation
	//------------------------------------------------

	pendingExists, err := ExistsPendingInvitationByEmail(
		ctx,
		email,
	)

	if err != nil {
		return err
	}

	if pendingExists {
		return shared.ErrInvitationAlreadyExists
	}

	//------------------------------------------------
	// Resolve Role ID
	//------------------------------------------------

	roleID, err := role.GetRoleIDByCode(
		ctx,
		roleCode,
	)

	if err != nil {
		return err
	}

	//------------------------------------------------
	// Generate Token
	//------------------------------------------------

	tokenBytes := make([]byte, 32)

	_, err = rand.Read(
		tokenBytes,
	)

	if err != nil {
		return err
	}

	rawToken := hex.EncodeToString(
		tokenBytes,
	)

	//------------------------------------------------
	// Hash Token
	//------------------------------------------------

	hash := sha256.Sum256(
		[]byte(rawToken),
	)

	tokenHash := hex.EncodeToString(
		hash[:],
	)

	//------------------------------------------------
	// Expiry
	//------------------------------------------------

	expiresAt := time.Now().
		Add(24 * time.Hour)

	//------------------------------------------------
	// Create Invitation
	//------------------------------------------------

	invitationID, err := CreateInvitationRepository(
		ctx,
		name,
		email,
		roleID,
		authUserID,
		tokenHash,
		scopeType,
		scopeID,
		expiresAt,
	)

	if err != nil {
		return err
	}

	//------------------------------------------------
	// Invitation Link
	//------------------------------------------------

	frontendURL := os.Getenv(
		"FRONTEND_URL",
	)

	inviteLink := fmt.Sprintf(
		"%s/accept-invitation?token=%s",
		frontendURL,
		rawToken,
	)

	//------------------------------------------------
	// Send Email
	//------------------------------------------------

	err = mail.SendInvitationEmail(
		email,
		name,
		inviteLink,
		roleCode,
	)

	if err != nil {

		deleteErr := DeleteInvitationByID(
			ctx,
			invitationID,
		)

		if deleteErr != nil {
			return deleteErr
		}

		return err
	}

	// Remove in production
	fmt.Print("\n")
	fmt.Print(rawToken)
	fmt.Print("\n")

	return nil
}

// DeleteInvitationService deletes an invitation.
func DeleteInvitationService(
	ctx context.Context,
	invitationID int64,
	userID string,
	role string,
) error {

	// ----------------------------------------------------------
	// Invitation Permissions
	// ----------------------------------------------------------

	roles, ok := InvitationPermissions[role]

	if !ok {
		return shared.ErrForbidden
	}

	// ----------------------------------------------------------
	// Fetch Invitation
	// ----------------------------------------------------------

	invitation, err := GetInvitationByIDRepository(
		ctx,
		invitationID,
		userID,
		roles,
	)

	if err != nil {
		return err
	}

	if invitation == nil {
		return shared.ErrInvitationNotFound
	}

	// ----------------------------------------------------------
	// Delete Invitation
	// ----------------------------------------------------------

	err = DeleteInvitationByID(
		ctx,
		invitationID,
	)

	if err != nil {
		return err
	}

	return nil
}



// GetInvitationByTokenService fetches an invitation
// using a public invitation token.
func GetInvitationByTokenService(
	ctx context.Context,
	token string,
) (*InvitationResponse, error) {

	// ----------------------------------------------------------
	// Hash Token
	// ----------------------------------------------------------

	hashedToken := utils.HashSHA256(token)

	// ----------------------------------------------------------
	// Get Invitation
	// ----------------------------------------------------------

	invitation, err := GetInvitationByTokenRepository(
		ctx,
		hashedToken,
	)

	if err != nil {
		return nil, err
	}

	if invitation == nil {
		return nil, shared.ErrInvitationNotFound
	}

	// ----------------------------------------------------------
	// Validate Status
	// ----------------------------------------------------------

	switch invitation.Status {

	case "revoked":
		return nil, shared.ErrInvitationRevoked

	case "accepted":
		return nil, shared.ErrInvitationAlreadyAccepted
	}

	// ----------------------------------------------------------
	// Check Expiration
	// ----------------------------------------------------------

	if time.Now().After(invitation.ExpiresAt) {

		return nil, shared.ErrInvitationExpired
	}

	return invitation, nil
}




//------------------------------------------------------------------------------------------------
//BORDER
//------------------------------------------------------------------------------------------------







// GetInvitationByIDService fetches invitation details
// for the provided invitation ID.
// GetInvitationByIDService fetches invitation details
// for the provided invitation ID.
// func GetInvitationByIDService(
// 	ctx context.Context,
// 	invitationID int64,
// 	userID string,
// 	role string,
// ) (*InvitationResponse, error) {

// 	// ----------------------------------------------------------
// 	// Invitation Permissions
// 	// ----------------------------------------------------------


// 	roles, ok := InvitationPermissions[role]

// 	if !ok {
// 		return nil, shared.ErrForbidden
// 	}

// 	// ----------------------------------------------------------
// 	// Fetch Invitation
// 	// ----------------------------------------------------------

// 	invitation, err := GetInvitationByIDRepository(
// 		ctx,
// 		invitationID,
// 		userID,
// 		roles,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if invitation == nil {
// 		return nil, shared.ErrInvitationNotFound
// 	}

// 	return invitation, nil
// }





func AcceptInvitationService(
	ctx context.Context,
	input AcceptInvitationInput,
) (*auth.LoginResponse, error) {

	// ----------------------------------------------------------
	// Validate Input
	// ----------------------------------------------------------

	input.Token = strings.TrimSpace(input.Token)
	input.FirstName = strings.TrimSpace(input.FirstName)
	input.LastName = strings.TrimSpace(input.LastName)
	input.ContactNumber = strings.TrimSpace(input.ContactNumber)

	if input.Token == "" {
		return nil, shared.ErrTokenRequired
	}

	if input.FirstName == "" {
		return nil, shared.ErrFirstNameRequired
	}

	if input.Password == "" {
		return nil, shared.ErrPasswordRequired
	}

	if !validation.IsStrongPassword(input.Password) {
		return nil, shared.ErrWeakPassword
	}

	if !validation.IsValidIndianPhone(input.ContactNumber) {
		return nil, shared.ErrInvalidPhoneNumber
	}

	// ----------------------------------------------------------
	// Get Invitation
	// ----------------------------------------------------------

	invitation, err := GetInvitationByTokenService(
		ctx,
		input.Token,
	)

	if err != nil {

		fmt.Printf(
			"[GetInvitationByTokenService] token=%s err=%v\n",
			input.Token,
			err,
		)

		return nil, err
	}

	// ----------------------------------------------------------
	// Check Existing User
	// ----------------------------------------------------------

	exists, err := auth.CheckUserExistsByEmail(
		ctx,
		invitation.Email,
	)

	if err != nil {

		fmt.Printf(
			"[CheckUserExistsByEmail] email=%s err=%v\n",
			invitation.Email,
			err,
		)

		return nil, err
	}

	if exists {
		return nil, shared.ErrEmailAlreadyExists
	}

	// ----------------------------------------------------------
	// Hash Password
	// ----------------------------------------------------------

	passwordHash, err := utils.HashPassword(
		input.Password,
	)

	if err != nil {

		fmt.Printf(
			"[HashPassword] err=%v\n",
			err,
		)

		return nil, err
	}

	// ----------------------------------------------------------
	// Begin Transaction
	// ----------------------------------------------------------

	tx, err := db.DB.Begin(ctx)

	if err != nil {

		fmt.Printf(
			"[BeginTx] err=%v\n",
			err,
		)

		return nil, err
	}

	defer tx.Rollback(ctx)

	// ----------------------------------------------------------
	// Create User
	// ----------------------------------------------------------

	user, err := auth.CreateUserTx(
		ctx,
		tx,
		auth.CreateUserInput{
			FirstName:     input.FirstName,
			LastName:      input.LastName,
			Email:         invitation.Email,
			PasswordHash:  passwordHash,
			ContactNumber: input.ContactNumber,
			Role:          invitation.Role.Key,
		},
	)

	if err != nil {

		fmt.Printf(
			"[CreateUserTx] role=%s email=%s err=%v\n",
			invitation.Role,
			invitation.Email,
			err,
		)

		return nil, err
	}

	// ----------------------------------------------------------
	// Create Role Extension
	// ----------------------------------------------------------

	switch invitation.Role.Key {

	case "state_admin":

		if *invitation.ScopeType != "state" {
			return nil, shared.ErrInvalidInvitationScope
		}

		if invitation.ScopeID == nil {
			return nil, shared.ErrInvalidInvitationScope
		}

		stateID, err := strconv.ParseInt(
			*invitation.ScopeID,
			10,
			32,
		)

		if err != nil {

			fmt.Printf(
				"[ParseStateID] scopeID=%s err=%v\n",
				*invitation.ScopeID,
				err,
			)

			return nil, err
		}

		err = stateAdmin.CreateStateAdminTx(
			ctx,
			tx,
			user.ID,
			int16(stateID),
		)

		if err != nil {

			fmt.Printf(
				"[CreateStateAdminTx] userID=%s stateID=%d err=%v\n",
				user.ID,
				stateID,
				err,
			)

			return nil, err
		}

	case "district_admin":

		fmt.Println("[Role] district_admin")

		if *invitation.ScopeType != "district" {
			return nil, shared.ErrInvalidInvitationScope
		}

		if invitation.ScopeID == nil {
			return nil, shared.ErrInvalidInvitationScope
		}

		districtID, err := strconv.ParseInt(
			*invitation.ScopeID,
			10,
			32,
		)

		if err != nil {

			fmt.Printf(
				"[ParseDistrictAdminID] scopeID=%s err=%v\n",
				*invitation.ScopeID,
				err,
			)

			return nil, err
		}

		err = districtAdmin.CreateDistrictAdminTx(
			ctx,
			tx,
			districtAdmin.CreateDistrictAdminInput{
				UserID:     user.ID,
				DistrictID: int32(districtID),
			},
		)

		if err != nil {

			fmt.Printf(
				"[CreateDistrictAdminTx] userID=%s districtID=%d err=%v\n",
				user.ID,
				districtID,
				err,
			)

			return nil, err
		}

	case "district_coach":

		fmt.Println("[Role] district_coach")

		if *invitation.ScopeType != "district" {
			return nil, shared.ErrInvalidInvitationScope
		}

		if invitation.ScopeID == nil {
			return nil, shared.ErrInvalidInvitationScope
		}

		districtID, err := strconv.ParseInt(
			*invitation.ScopeID,
			10,
			32,
		)

		if err != nil {

			fmt.Printf(
				"[ParseDistrictCoachID] scopeID=%s err=%v\n",
				*invitation.ScopeID,
				err,
			)

			return nil, err
		}

		err = districtCoach.CreateDistrictCoachTx(
			ctx,
			tx,
			districtCoach.CreateDistrictCoachInput{
				UserID:     user.ID,
				DistrictID: int32(districtID),
			},
		)

		if err != nil {

			fmt.Printf(
				"[CreateDistrictCoachTx] userID=%s districtID=%d err=%v\n",
				user.ID,
				districtID,
				err,
			)

			return nil, err
		}

	case "academy_admin":

		fmt.Println("[Role] academy_admin")

		if *invitation.ScopeType != "academy" {
			return nil, shared.ErrInvalidInvitationScope
		}

		if invitation.ScopeID == nil {
			return nil, shared.ErrInvalidInvitationScope
		}

		if !validation.IsValidUUID(*invitation.ScopeID) {
			return nil, shared.ErrInvalidInvitationScope
		}

		err = academyAdmin.CreateAcademyAdminTx(
			ctx,
			tx,
			academyAdmin.CreateAcademyAdminInput{
				UserID:    user.ID,
				AcademyID: *invitation.ScopeID,
			},
		)

		if err != nil {

			fmt.Printf(
				"[CreateAcademyAdminTx] userID=%s academyID=%s err=%v\n",
				user.ID,
				*invitation.ScopeID,
				err,
			)

			return nil, err
		}

	case "academy_coach":

		fmt.Println("[Role] academy_coach")

		if *invitation.ScopeType != "academy" {
			return nil, shared.ErrInvalidInvitationScope
		}

		if invitation.ScopeID == nil {
			return nil, shared.ErrInvalidInvitationScope
		}

		if !validation.IsValidUUID(*invitation.ScopeID) {
			return nil, shared.ErrInvalidInvitationScope
		}

		err = academyCoach.CreateAcademyCoachTx(
			ctx,
			tx,
			academyCoach.CreateAcademyCoachInput{
				UserID:    user.ID,
				AcademyID: *invitation.ScopeID,
			},
		)

		if err != nil {

			fmt.Printf(
				"[CreateAcademyCoachTx] userID=%s academyID=%s err=%v\n",
				user.ID,
				*invitation.ScopeID,
				err,
			)

			return nil, err
		}

	case "player":

		fmt.Println("[Role] player")

		if *invitation.ScopeType != "academy" {
			return nil, shared.ErrInvalidInvitationScope
		}

		if invitation.ScopeID == nil {
			return nil, shared.ErrInvalidInvitationScope
		}

		if !validation.IsValidUUID(*invitation.ScopeID) {
			return nil, shared.ErrInvalidInvitationScope
		}

		err = player.CreatePlayerTx(
			ctx,
			tx,
			player.CreatePlayerInput{
				UserID:       user.ID,
				AcademyID:    *invitation.ScopeID,
				RegisteredBy: invitation.InvitedBy,
			},
		)

		if err != nil {

			fmt.Printf(
				"[CreatePlayerTx] userID=%s academyID=%s err=%v\n",
				user.ID,
				*invitation.ScopeID,
				err,
			)

			return nil, err
		}
	}

	// ----------------------------------------------------------
	// Mark Invitation Accepted
	// ----------------------------------------------------------

	err = AcceptInvitationTx(
		ctx,
		tx,
		invitation.ID,
		user.ID,
	)

	if err != nil {

		fmt.Printf(
			"[AcceptInvitationTx] invitationID=%s userID=%s err=%v\n",
			invitation.ID,
			user.ID,
			err,
		)

		return nil, err
	}

	// ----------------------------------------------------------
	// Commit Transaction
	// ----------------------------------------------------------

	err = tx.Commit(ctx)

	if err != nil {

		fmt.Printf(
			"[CommitTx] err=%v\n",
			err,
		)

		return nil, err
	}

	// ----------------------------------------------------------
	// Generate JWT
	// ----------------------------------------------------------

	token, err := jwtToken.GenerateJWT(
		user.ID,
		user.Role.ID,
		user.Role.Code,
	)

	if err != nil {

		fmt.Printf(
			"[GenerateJWT] userID=%s err=%v\n",
			user.ID,
			err,
		)

		return nil, err
	}

	// ----------------------------------------------------------
	// Response
	// ----------------------------------------------------------

	response := &auth.LoginResponse{
		Token: token,
		User: auth.UserResponse{
			ID:            user.ID,
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Email:         user.Email,
			ContactNumber: user.ContactNumber,
			Role:          user.Role,
		},
	}

	return response, nil
}

