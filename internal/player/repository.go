package player

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func CreatePlayerTx(
	ctx context.Context,
	tx pgx.Tx,
	input CreatePlayerInput,
) error {

	_, err := tx.Exec(
		ctx,
		`
		INSERT INTO players (
			user_id,
			academy_id,
			registered_by
		)
		VALUES (
			$1,
			$2,
			$3
		)
		`,
		input.UserID,
		input.AcademyID,
		input.RegisteredBy,
	)

	return err
}