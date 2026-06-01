package invitation

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"

	"cobackend/internal/db"
	"cobackend/internal/shared"
)

func GetInvitationsRepository(
	ctx context.Context,
	userID string,
	roles []string,
	isSuperAdmin bool,
	query InvitationsQueryParams,
) (*PaginatedInvitations, error) {

	// ----------------------------------------------------------
	// Pagination
	// ----------------------------------------------------------

	offset := (query.Page - 1) * query.Limit

	// ----------------------------------------------------------
	// Count Query
	// ----------------------------------------------------------

	var total int

	countQuery := `
		SELECT COUNT(*)
		FROM invitations i
		JOIN roles r
			ON r.id = i.role_id
		WHERE 1=1
	`

	countArgs := []interface{}{}

	// ----------------------------------------------------------
	// RBAC Filtering
	// ----------------------------------------------------------

	if !isSuperAdmin {

		countQuery += `
			AND i.invited_by = $` +
			strconv.Itoa(len(countArgs)+1)

		countArgs = append(
			countArgs,
			userID,
		)

		countQuery += `
			AND r.code = ANY($` +
			strconv.Itoa(len(countArgs)+1) +
			`)`

		countArgs = append(
			countArgs,
			roles,
		)
	}

	// ----------------------------------------------------------
	// Search
	// ----------------------------------------------------------

	if query.Search != "" {

		countQuery += `
			AND i.email ILIKE '%' || $` +
			strconv.Itoa(len(countArgs)+1) +
			` || '%'
		`

		countArgs = append(
			countArgs,
			query.Search,
		)
	}

	// ----------------------------------------------------------
	// Status
	// ----------------------------------------------------------

	if query.Status != "" {

		countQuery += `
			AND i.status = $` +
			strconv.Itoa(len(countArgs)+1)

		countArgs = append(
			countArgs,
			query.Status,
		)
	}

	// ----------------------------------------------------------
	// Role
	// ----------------------------------------------------------

	if query.Role != "" {

		countQuery += `
			AND r.code = $` +
			strconv.Itoa(len(countArgs)+1)

		countArgs = append(
			countArgs,
			query.Role,
		)
	}

	err := db.DB.QueryRow(
		ctx,
		countQuery,
		countArgs...,
	).Scan(&total)

	if err != nil {
		return nil, err
	}

	// ----------------------------------------------------------
	// Data Query
	// ----------------------------------------------------------

	dataQuery := `
		SELECT
			i.id,
			i.email,

			r.code,
			r.display_name,

			i.scope_type,
			i.scope_id,

			COALESCE(
				s.name,
				d.name,
				a.name
			) AS organization_name,

			i.status,

			u.id,
			u.first_name,
			u.last_name,

			i.expires_at,
			i.created_at

		FROM invitations i

		JOIN roles r
			ON r.id = i.role_id

		JOIN users u
			ON u.id = i.invited_by

		LEFT JOIN states s
			ON i.scope_type = 'state'
			AND s.id::text = i.scope_id

		LEFT JOIN districts d
			ON i.scope_type = 'district'
			AND d.id::text = i.scope_id

		LEFT JOIN academies a
			ON i.scope_type = 'academy'
			AND a.id::text = i.scope_id

		WHERE 1=1
	`

	dataArgs := []interface{}{}

	// ----------------------------------------------------------
	// RBAC Filtering
	// ----------------------------------------------------------

	if !isSuperAdmin {

		dataQuery += `
			AND i.invited_by = $` +
			strconv.Itoa(len(dataArgs)+1)

		dataArgs = append(
			dataArgs,
			userID,
		)

		dataQuery += `
			AND r.code = ANY($` +
			strconv.Itoa(len(dataArgs)+1) +
			`)`

		dataArgs = append(
			dataArgs,
			roles,
		)
	}

	// ----------------------------------------------------------
	// Search
	// ----------------------------------------------------------

	if query.Search != "" {

		dataQuery += `
			AND i.email ILIKE '%' || $` +
			strconv.Itoa(len(dataArgs)+1) +
			` || '%'
		`

		dataArgs = append(
			dataArgs,
			query.Search,
		)
	}

	// ----------------------------------------------------------
	// Status
	// ----------------------------------------------------------

	if query.Status != "" {

		dataQuery += `
			AND i.status = $` +
			strconv.Itoa(len(dataArgs)+1)

		dataArgs = append(
			dataArgs,
			query.Status,
		)
	}

	// ----------------------------------------------------------
	// Role
	// ----------------------------------------------------------

	if query.Role != "" {

		dataQuery += `
			AND r.code = $` +
			strconv.Itoa(len(dataArgs)+1)

		dataArgs = append(
			dataArgs,
			query.Role,
		)
	}

	// ----------------------------------------------------------
	// Sorting
	// ----------------------------------------------------------

	dataQuery += `
		ORDER BY ` + query.SortBy + ` ` + query.Order

	// ----------------------------------------------------------
	// Pagination
	// ----------------------------------------------------------

	dataQuery += `
		LIMIT $` + strconv.Itoa(len(dataArgs)+1) + `
		OFFSET $` + strconv.Itoa(len(dataArgs)+2)

	dataArgs = append(
		dataArgs,
		query.Limit,
		offset,
	)

	rows, err := db.DB.Query(
		ctx,
		dataQuery,
		dataArgs...,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	invitations := []InvitationResponse{}

	for rows.Next() {

		var invitation InvitationResponse

		var roleKey string
		var roleLabel string

		var scopeType *string
		var scopeID *string
		var organizationName *string

		var createdByID string
		var firstName string
		var lastName *string

		err := rows.Scan(
			&invitation.ID,
			&invitation.Email,

			&roleKey,
			&roleLabel,

			&scopeType,
			&scopeID,
			&organizationName,

			&invitation.Status,

			&createdByID,
			&firstName,
			&lastName,

			&invitation.ExpiresAt,
			&invitation.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		invitation.Role = RoleResponse{
			Key:   roleKey,
			Label: roleLabel,
		}

		if scopeType != nil &&
			scopeID != nil &&
			organizationName != nil {

			invitation.Organization =
				&OrganizationResponse{
					Type: *scopeType,
					ID:   *scopeID,
					Name: *organizationName,
				}
		}

		fullName := firstName

		if lastName != nil {
			fullName += " " + *lastName
		}

		invitation.CreatedBy = UserSummary{
			ID:   createdByID,
			Name: fullName,
		}

		invitations = append(
			invitations,
			invitation,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	totalPages := 1

	if total > 0 {
		totalPages =
			(total + query.Limit - 1) /
				query.Limit
	}

	return &PaginatedInvitations{
		Items:       invitations,
		Page:        query.Page,
		Limit:       query.Limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     query.Page < totalPages,
		HasPrevious: query.Page > 1,
	}, nil
}


// CreateInvitationRepository creates a new invitation
// and returns the created invitation ID.
func CreateInvitationRepository(
	ctx context.Context,
	name string,
	email string,
	roleID int,
	invitedBy string,
	tokenHash string,
	scopeType string,
	scopeID string,
	expiresAt time.Time,
) (int64, error) {

	query := `
		INSERT INTO invitations (
			name,
			email,
			role_id,
			invited_by,
			token_hash,
			scope_type,
			scope_id,
			expires_at
		)
		VALUES (
			$1, $2, $3, $4,
			$5, $6, $7, $8
		)
		RETURNING id
	`

	var invitationID int64

	err := db.DB.QueryRow(
		ctx,
		query,
		name,
		email,
		roleID,
		invitedBy,
		tokenHash,
		scopeType,
		scopeID,
		expiresAt,
	).Scan(&invitationID)

	if err != nil {
		return 0, err
	}

	return invitationID, nil
}


// GetInvitationByTokenRepository fetches invitation
// details using invitation token hash.
func GetInvitationByTokenRepository(
	ctx context.Context,
	hashedToken string,
) (*InvitationResponse, error) {

	query := `
		SELECT
			i.id,
			i.email,

			r.code,
			r.display_name,

			i.scope_type,
			i.scope_id,

			COALESCE(
				s.name,
				d.name,
				a.name
			) AS organization_name,

			i.status,

			i.invited_by,
			i.accepted_at,
			i.used_by,

			u.id,
			u.first_name,
			u.last_name,

			i.expires_at,
			i.created_at

		FROM invitations i

		JOIN roles r
			ON r.id = i.role_id

		JOIN users u
			ON u.id = i.invited_by

		LEFT JOIN states s
			ON i.scope_type = 'state'
			AND s.id::text = i.scope_id

		LEFT JOIN districts d
			ON i.scope_type = 'district'
			AND d.id::text = i.scope_id

		LEFT JOIN academies a
			ON i.scope_type = 'academy'
			AND a.id::text = i.scope_id

		WHERE i.token_hash = $1

		LIMIT 1
	`

	var invitation InvitationResponse

	var roleKey string
	var roleLabel string

	var scopeType *string
	var scopeID *string
	var organizationName *string

	var invitedBy string

	var acceptedAt *time.Time
	var usedBy *string

	var createdByID string
	var firstName string
	var lastName *string

	err := db.DB.QueryRow(
		ctx,
		query,
		hashedToken,
	).Scan(
		&invitation.ID,
		&invitation.Email,

		&roleKey,
		&roleLabel,

		&scopeType,
		&scopeID,
		&organizationName,

		&invitation.Status,

		&invitedBy,
		&acceptedAt,
		&usedBy,

		&createdByID,
		&firstName,
		&lastName,

		&invitation.ExpiresAt,
		&invitation.CreatedAt,
	)

	if err != nil {

		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return nil, nil
		}

		return nil, err
	}

	// ----------------------------------------------------------
	// Primitive Fields
	// ----------------------------------------------------------

	invitation.RoleCode = roleKey

	invitation.ScopeType = scopeType
	invitation.ScopeID = scopeID

	invitation.InvitedBy = invitedBy

	invitation.AcceptedAt = acceptedAt
	invitation.UsedBy = usedBy

	// ----------------------------------------------------------
	// Role DTO
	// ----------------------------------------------------------

	invitation.Role =
		RoleResponse{
			Key: roleKey,

			Label: roleLabel,
		}

	// ----------------------------------------------------------
	// Organization DTO
	// ----------------------------------------------------------

	if scopeType != nil &&
		scopeID != nil &&
		organizationName != nil {

		invitation.Organization =
			&OrganizationResponse{
				Type: *scopeType,
				ID:   *scopeID,
				Name: *organizationName,
			}
	}

	// ----------------------------------------------------------
	// Created By DTO
	// ----------------------------------------------------------

	fullName := firstName

	if lastName != nil {
		fullName += " " + *lastName
	}

	invitation.CreatedBy =
		UserSummary{
			ID: createdByID,

			Name: fullName,
		}

	return &invitation, nil
}


func DeleteInvitationByID(
	ctx context.Context,
	id int64,
) error {

	query := `
		DELETE FROM invitations
		WHERE id = $1
	`

	commandTag, err := db.DB.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return shared.ErrInvitationNotFound
	}

	return nil
}


// GetInvitationByIDRepository fetches invitation details
// required for permission validation.
func GetInvitationByIDRepository(
	ctx context.Context,
	invitationID int64,
	userID string,
	roles []string,
) (*InvitationPermissionCheck, error) {

	var invitation InvitationPermissionCheck

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT
			i.id,
			i.invited_by,
			r.code::text,
			i.status::text,
			i.expires_at
		FROM invitations i
		JOIN roles r
			ON r.id = i.role_id
		WHERE i.id = $1
		AND i.invited_by = $2
		AND r.code::text = ANY($3)
		LIMIT 1
		`,
		invitationID,
		userID,
		roles,
	).Scan(
		&invitation.ID,
		&invitation.InvitedBy,
		&invitation.RoleCode,
		&invitation.Status,
		&invitation.ExpiresAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &invitation, nil
}


func ExistsPendingInvitationByEmail(
	ctx context.Context,
	email string,
) (bool, error) {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM invitations
			WHERE email = $1
			AND status = 'pending'
		)
	`

	var exists bool

	err := db.DB.QueryRow(
		ctx,
		query,
		email,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

// AcceptInvitationTx marks an invitation
// as accepted within an existing transaction.
func AcceptInvitationTx(
	ctx context.Context,
	tx pgx.Tx,
	invitationID int64,
	usedBy string,
) error {

	commandTag, err := tx.Exec(
		ctx,
		`
		UPDATE invitations
		SET
			status = 'accepted',
			used_by = $2,
			accepted_at = CURRENT_TIMESTAMP,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		`,
		invitationID,
		usedBy,
	)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return shared.ErrInvitationNotFound
	}

	return nil
}

//------------------------------------------------------------------------------------------------
//BORDER
//------------------------------------------------------------------------------------------------









// GetInvitationByIDRepository fetches invitation details
// for the provided invitation ID.
// func GetInvitationByIDRepository(
// 	ctx context.Context,
// 	invitationID int64,
// 	userID string,
// 	roles []string,
// ) (*InvitationResponse, error) {

// 	var invitation InvitationResponse

// 	err := db.DB.QueryRow(
// 		ctx,
// 		`
// 		SELECT
// 			i.id,
// 			i.email,
// 			r.code,
// 			i.status,
// 			i.expires_at,
// 			i.created_at
// 		FROM invitations i
// 		JOIN roles r
// 			ON r.id = i.role_id
// 		WHERE i.id = $1
// 		AND i.invited_by = $2
// 		AND r.code = ANY($3)
// 		LIMIT 1
// 		`,
// 		invitationID,
// 		userID,
// 		roles,
// 	).Scan(
// 		&invitation.ID,
// 		&invitation.Email,
// 		&invitation.Role,
// 		&invitation.Status,
// 		&invitation.ExpiresAt,
// 		&invitation.CreatedAt,
// 	)

// 	if err != nil {

// 		if errors.Is(err, pgx.ErrNoRows) {
// 			return nil, nil
// 		}

// 		return nil, err
// 	}

// 	return &invitation, nil
// }

// RevokeInvitationRepository revokes an existing invitation.
func RevokeInvitationRepository(
	ctx context.Context,
	invitationID int64,
) error {

	commandTag, err := db.DB.Exec(
		ctx,
		`
		UPDATE invitations
		SET
			status = 'revoked',
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		`,
		invitationID,
	)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return shared.ErrInvitationNotFound
	}

	return nil
}



