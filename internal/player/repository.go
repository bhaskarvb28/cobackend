package player

import (
	"context"
	"github.com/jackc/pgx/v5"

	"cobackend/internal/db"
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

func GetAvailableShootingEventsRepository(
	ctx context.Context,
	userID string,
	disciplineID int16,
) ([]ShootingEventResponse, error) {

	query := `
		SELECT DISTINCT
			se.id,
			se.code,
			se.display_name,
			sd.meters

		FROM shooting_events se

		INNER JOIN academy_building_events abe
			ON abe.shooting_event_id = se.id

		INNER JOIN academy_buildings ab
			ON ab.id = abe.academy_building_id

		INNER JOIN players p
			ON p.academy_id = ab.academy_id

		LEFT JOIN shooting_distances sd
			ON sd.id = se.distance_id

		WHERE
			p.user_id = $1
			AND se.discipline_id = $2

		ORDER BY
			se.display_name ASC
	`

	rows, err := db.DB.Query(
		ctx,
		query,
		userID,
		disciplineID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []ShootingEventResponse{}

	for rows.Next() {

		var event ShootingEventResponse

		err := rows.Scan(
			&event.ID,
			&event.Code,
			&event.DisplayName,
			&event.Distance,
		)

		if err != nil {
			return nil, err
		}

		events = append(
			events,
			event,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func GetCompatibleBuildingsRepository(
	ctx context.Context,
	userID string,
	shootingEventID int16,
) ([]CompatibleBuildingResponse, error) {

	query := `
		SELECT DISTINCT
			ab.id,
			ab.building_name,
			ab.created_at,
			ab.updated_at

		FROM academy_buildings ab

		INNER JOIN academy_building_events abe
			ON abe.academy_building_id = ab.id

		INNER JOIN players p
			ON p.academy_id = ab.academy_id

		WHERE
			p.user_id = $1
			AND abe.shooting_event_id = $2

		ORDER BY
			ab.building_name ASC
	`

	rows, err := db.DB.Query(
		ctx,
		query,
		userID,
		shootingEventID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	buildings := []CompatibleBuildingResponse{}

	for rows.Next() {

		var building CompatibleBuildingResponse

		err := rows.Scan(
			&building.ID,
			&building.BuildingName,
			&building.CreatedAt,
			&building.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		buildings = append(
			buildings,
			building,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return buildings, nil
}

