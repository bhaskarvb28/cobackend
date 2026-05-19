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
	StateID *int,
	DistrictID *int,
	AcademyID *int,
	DisciplinesSpecialized []int,
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
			state_id,
			district_id,
			academy_id,
			disciplines_specialized,
			expires_at
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9
		)
		`,
		email,
		roleID,
		invitedBy,
		token,
		StateID,
		DistrictID,
		AcademyID,
		DisciplinesSpecialized,
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
			i.id,
			i.email,
			i.role_id,
			r.name,
			i.token,
			i.disciplines_specialized,
			i.state_id,
			i.district_id,
			i.academy_id,
			i.status,
			i.expires_at
		FROM invitations i
		JOIN roles r
			ON r.id = i.role_id
		WHERE i.token = $1
		`,
		token,
	).Scan(
		&invitation.ID,
		&invitation.Email,
		&invitation.RoleID,
		&invitation.RoleName,
		&invitation.Token,
		&invitation.DisciplinesSpecialized,
		&invitation.StateID,
		&invitation.DistrictID,
		&invitation.AcademyID,
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