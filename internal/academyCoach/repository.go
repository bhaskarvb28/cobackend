package academyCoach

import (
	"context"
	"github.com/jackc/pgx/v5"

	"cobackend/internal/db"
)

func CreateAcademyCoachTx(
	ctx context.Context,
	tx pgx.Tx,
	input CreateAcademyCoachInput,
) error {

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO academy_coaches (
			user_id,
			academy_id
		)
		VALUES ($1, $2)
		`,
		input.UserID,
		input.AcademyID,
	)

	return err
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

func CheckAcademyCoachExists(
	ctx context.Context,
	academyCoachID string,
) (bool, error) {

	var exists bool

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT EXISTS(
			SELECT 1
			FROM academy_coaches
			WHERE profile_id = $1
		)
		`,
		academyCoachID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func CheckCoachBelongsToAcademy(
	ctx context.Context,
	academyCoachID string,
	academyID int,
) (bool, error) {

	var exists bool

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT EXISTS(
			SELECT 1
			FROM academy_coaches
			WHERE profile_id = $1
			AND academy_id = $2
		)
		`,
		academyCoachID,
		academyID,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}