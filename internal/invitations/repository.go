package invitations

import (
	"cobackend/internal/db"
	"context"
	"time"
	"github.com/jackc/pgx/v5"
)

func ExistsPendingInvitationByEmail(
	ctx context.Context,
	email string,
) (bool, error) {

	var exists bool

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT EXISTS (
			SELECT 1
			FROM invitations
			WHERE email = $1
			AND status ILIKE 'pending'
		)
		`,
		email,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func CreateInvitationRepository(
	ctx context.Context,
	email string,
	roleID string,
	invitedBy string,
	token string,
	assignedStateID string,
	expiresAt time.Time,
) error {

	_, err := db.DB.Exec(
		ctx,
		`
		INSERT INTO invitations (
			email,
			role_id,
			invited_by,
			token,
			assigned_state_id,
			expires_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)
		`,
		email,
		roleID,
		invitedBy,
		token,
		assignedStateID,
		expiresAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func DeleteInvitationByToken(
	ctx context.Context,
	token string,
) error {
	_, err := db.DB.Exec(
		ctx,
		`
		DELETE FROM invitations
		WHERE token = $1
		`,
		token,
	)
	
	if err != nil {
		return err
	}

	return nil
}

func GetInvitationByToken(
	ctx context.Context,
	token string,
) (Invitation, error) {

	var invitation Invitation

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT
			id,
			email,
			role_id,
			token,
			assigned_state_id,
			status,
			expires_at
		FROM invitations
		WHERE token = $1
		`,
		token,
	).Scan(
		&invitation.ID,
		&invitation.Email,
		&invitation.RoleID,
		&invitation.Token,
		&invitation.AssignedStateID,
		&invitation.Status,
		&invitation.ExpiresAt,
	)

	if err != nil {
		return Invitation{}, err
	}

	return invitation, nil
}

func MarkInvitationUsedTx(
	ctx context.Context,
	tx pgx.Tx,
	invitationID string,
) error {

	_, err := tx.Exec(
		ctx,
		`
		UPDATE invitations
		SET status = 'accepted'
		WHERE id = $1
		`,
		invitationID,
	)

	if err != nil {
		return err
	}

	return nil
}