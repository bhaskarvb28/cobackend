package session

import (
	"context"

	"cobackend/internal/shared"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// repository.go

func StartPracticeSessionRepository(
	ctx context.Context,
	tx pgx.Tx,
	userID string,
	input StartPracticeSessionInput,
) (*PracticeSessionResponse, error) {

	query := `
		INSERT INTO practice_sessions (
			player_user_id,
			academy_building_lane_id,
			shooting_event_id
		)
		VALUES (
			$1,
			$2,
			$3
		)
		RETURNING
			id,
			player_user_id,
			academy_building_lane_id,
			shooting_event_id,
			status,
			total_score,
			total_shot_count,
			started_at,
			ended_at,
			created_at,
			updated_at
	`

	var session PracticeSessionResponse

	err := tx.QueryRow(
		ctx,
		query,
		userID,
		input.AcademyBuildingLaneID,
		input.ShootingEventID,
	).Scan(
		&session.ID,
		&session.PlayerUserID,
		&session.AcademyBuildingLaneID,
		&session.ShootingEventID,
		&session.Status,
		&session.TotalScore,
		&session.TotalShotCount,
		&session.StartedAt,
		&session.EndedAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {

		if pgError, ok := err.(*pgconn.PgError); ok {

			switch pgError.ConstraintName {

			case "unique_active_lane_session":
				return nil, shared.ErrLaneAlreadyOccupied

			case "unique_active_player_session":
				return nil, shared.ErrPlayerAlreadyHasActiveSession
			}
		}

		return nil, err
	}

	return &session, nil
}