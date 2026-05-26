package invitation

import (
	"context"
	"errors"
	"time"
	"strconv"

	"github.com/jackc/pgx/v5"

	"cobackend/internal/db"
	"cobackend/internal/shared"
)

// CreateInvitationRepository creates a new invitation
// and returns the created invitation response.
func CreateInvitationRepository(
	ctx context.Context,
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
			email,
			role_id,
			invited_by,
			token_hash,
			scope_type,
			scope_id,
			expires_at
		)
		VALUES (
			$1, $2, $3,
			$4, $5, $6, $7
		)
		RETURNING id
	`

	var invitationID int64

	err := db.DB.QueryRow(
		ctx,
		query,
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

// GetInvitationsRepository fetches all invitations
// visible to the authenticated user.
func GetInvitationsRepository(
	ctx context.Context,
	userID string,
	roles []string,
	query InvitationsQueryParams,
) (*PaginatedInvitations, error) {

	// ----------------------------------------------------------
	// Pagination
	// ----------------------------------------------------------

	offset := (query.Page - 1) * query.Limit

	// ----------------------------------------------------------
	// Count Total Records
	// ----------------------------------------------------------

	var total int

	countQuery := `
		SELECT COUNT(*)
		FROM invitations i
		JOIN roles r
			ON r.id = i.role_id
		WHERE i.invited_by = $1
		AND r.code = ANY($2)
	`

	countArgs := []interface{}{
		userID,
		roles,
	}

	// Search Filter
	if query.Search != "" {

		countQuery += `
			AND (
				i.email ILIKE '%' || $3 || '%'
			)
		`

		countArgs = append(
			countArgs,
			query.Search,
		)
	}

	// Status Filter
	if query.Status != "" {

		countQuery += `
			AND i.status = $` + strconv.Itoa(len(countArgs)+1)

		countArgs = append(
			countArgs,
			query.Status,
		)
	}

	// Role Filter
	if query.Role != "" {

		countQuery += `
			AND r.code = $` + strconv.Itoa(len(countArgs)+1)

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
	// Fetch Invitations
	// ----------------------------------------------------------

	dataQuery := `
		SELECT
			i.id,
			i.email,
			r.code,
			i.scope_type,
			i.scope_id,
			i.status,
			i.expires_at,
			i.created_at
		FROM invitations i
		JOIN roles r
			ON r.id = i.role_id
		WHERE i.invited_by = $1
		AND r.code = ANY($2)
	`

	dataArgs := []interface{}{
		userID,
		roles,
	}

	// Search Filter
	if query.Search != "" {

		dataQuery += `
			AND (
				i.email ILIKE '%' || $3 || '%'
			)
		`

		dataArgs = append(
			dataArgs,
			query.Search,
		)
	}

	// Status Filter
	if query.Status != "" {

		dataQuery += `
			AND i.status = $` + strconv.Itoa(len(dataArgs)+1)

		dataArgs = append(
			dataArgs,
			query.Status,
		)
	}

	// Role Filter
	if query.Role != "" {

		dataQuery += `
			AND r.code = $` + strconv.Itoa(len(dataArgs)+1)

		dataArgs = append(
			dataArgs,
			query.Role,
		)
	}

	// Sorting
	dataQuery += `
		ORDER BY ` + query.SortBy + ` ` + query.Order

	// Pagination
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

		err := rows.Scan(
			&invitation.ID,
			&invitation.Email,
			&invitation.Role,
			&invitation.ScopeType,
			&invitation.ScopeID,
			&invitation.Status,
			&invitation.ExpiresAt,
			&invitation.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		invitations = append(
			invitations,
			invitation,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// ----------------------------------------------------------
	// Pagination Metadata
	// ----------------------------------------------------------

	totalPages := 0

	if total > 0 {
		totalPages = (total + query.Limit - 1) / query.Limit
	}

	response := &PaginatedInvitations{
		Items:       invitations,
		Page:        query.Page,
		Limit:       query.Limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     query.Page < totalPages,
		HasPrevious: query.Page > 1,
	}

	return response, nil
}

// GetInvitationByIDRepository fetches invitation details
// for the provided invitation ID.
func GetInvitationByIDRepository(
	ctx context.Context,
	invitationID int64,
	userID string,
	roles []string,
) (*InvitationResponse, error) {

	var invitation InvitationResponse

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT
			i.id,
			i.email,
			r.code,
			i.status,
			i.expires_at,
			i.created_at
		FROM invitations i
		JOIN roles r
			ON r.id = i.role_id
		WHERE i.id = $1
		AND i.invited_by = $2
		AND r.code = ANY($3)
		LIMIT 1
		`,
		invitationID,
		userID,
		roles,
	).Scan(
		&invitation.ID,
		&invitation.Email,
		&invitation.Role,
		&invitation.Status,
		&invitation.ExpiresAt,
		&invitation.CreatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &invitation, nil
}

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

// GetInvitationByTokenRepository fetches invitation
// details using a hashed invitation token.
func GetInvitationByTokenRepository(
	ctx context.Context,
	hashedToken string,
) (*InvitationResponse, error) {

	var invitation InvitationResponse

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT
			i.id,
			i.email,
			r.code,
			i.invited_by,
			i.scope_type,
			i.scope_id,
			i.status,
			i.expires_at,
			i.created_at
		FROM invitations i
		JOIN roles r
			ON r.id = i.role_id
		WHERE i.token_hash = $1
		LIMIT 1
		`,
		hashedToken,
	).Scan(
		&invitation.ID,
		&invitation.Email,
		&invitation.Role,
		&invitation.InvitedBy,
		&invitation.ScopeType,
		&invitation.ScopeID,
		&invitation.Status,
		&invitation.ExpiresAt,
		&invitation.CreatedAt,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &invitation, nil
}