package academyAdmin

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CreateAcademyAdminTx(
	ctx context.Context,
	tx pgx.Tx,
	profileID string,
	academyID *int,
	gstin string,
	registrationProof string,
	dpdpConsent bool,
) error {

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO academY_admins (
			profile_id,
			academy_id,
			gstin,
			registration_proof,
			dpdp_consent
		)
		VALUES ($1, $2, $3, $4, $5)
		`,
		profileID,
		academyID,
		gstin,
		registrationProof,
		dpdpConsent,
	)

	if err != nil {
		return err
	}

	return nil
}