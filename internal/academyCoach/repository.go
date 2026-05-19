package academyCoach

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func CreateAcademyCoachTx(
	ctx context.Context,
	tx pgx.Tx,
	profileID string,
	academyID int,
	dpdpConsent bool,
	coachingCredentialsProof string,
) error {

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO academy_coaches(
			profile_id,
			academy_id,
			dpdp_consent,
			coaching_credentials_proof
		)
		VALUES (
			$1, $2, $3, $4
		)
		`,
		profileID,
		academyID,
		dpdpConsent,
		coachingCredentialsProof,
	)

	if err != nil {
		return err
	}

	return nil
}

func AddAcademyCoachDisciplineTx(
	ctx context.Context,
	tx pgx.Tx,
	profileID string,
	categoryID int32,
) error {


	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO academy_coach_disciplines (
			coach_profile_id,
			category_id
		)
		VALUES ($1, $2)
		`,
		profileID,
		categoryID,
	)

	if err != nil {
		return err
	}

	return nil
}