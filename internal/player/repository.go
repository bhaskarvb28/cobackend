package player

import (
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"

	"cobackend/internal/db"
	"cobackend/internal/profile"
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

func GetAcademyPlayersRepository(
	ctx context.Context,
	academyID string,
	query GetAcademyPlayersQuery,
) (PaginatedPlayers, error) {

	players := []PlayerListItemResponse{}

	// ----------------------------------------------------------
	// Pagination
	// ----------------------------------------------------------

	offset := 0

	if query.Limit > 0 {

		offset = (query.Page - 1) * query.Limit
	}

	// ----------------------------------------------------------
	// Base Query
	// ----------------------------------------------------------

	baseQuery := `
		FROM players p

		INNER JOIN users u
			ON u.id = p.user_id

		LEFT JOIN player_personal_info ppi
			ON ppi.player_user_id = p.user_id

		LEFT JOIN player_disciplines pd
			ON pd.player_user_id = p.user_id
			AND pd.is_primary = true

		LEFT JOIN disciplines d
			ON d.id = pd.discipline_id

		LEFT JOIN academy_coaches ac
			ON ac.user_id = p.current_coach_user_id

		LEFT JOIN users coach_user
			ON coach_user.id = ac.user_id

		WHERE p.academy_id = $1
	`

	args := []interface{}{
		academyID,
	}

	argPos := 2

	// ----------------------------------------------------------
	// Search Filter
	// ----------------------------------------------------------

	if query.Search != "" {

		baseQuery += `
			AND (
				u.first_name ILIKE $` + strconv.Itoa(argPos) + `
				OR u.last_name ILIKE $` + strconv.Itoa(argPos) + `
				OR u.email ILIKE $` + strconv.Itoa(argPos) + `
			)
		`

		args = append(
			args,
			"%"+query.Search+"%",
		)

		argPos++
	}

	// ----------------------------------------------------------
	// Discipline Filter
	// ----------------------------------------------------------

	if query.DisciplineID > 0 {

		baseQuery += `
			AND pd.discipline_id = $` + strconv.Itoa(argPos)

		args = append(
			args,
			query.DisciplineID,
		)

		argPos++
	}

	// ----------------------------------------------------------
	// Coach Assignment Filter
	// ----------------------------------------------------------

	if query.CoachAssigned != nil {

		if *query.CoachAssigned {

			baseQuery += `
				AND p.current_coach_user_id IS NOT NULL
			`

		} else {

			baseQuery += `
				AND p.current_coach_user_id IS NULL
			`
		}
	}

	// ----------------------------------------------------------
	// Status Filter
	// ----------------------------------------------------------

	if query.Status != "" {

		baseQuery += `
			AND p.status = $` + strconv.Itoa(argPos)

		args = append(
			args,
			query.Status,
		)

		argPos++
	}

	// ----------------------------------------------------------
	// Count Query
	// ----------------------------------------------------------

	countQuery := `
		SELECT COUNT(DISTINCT p.user_id)
	` + baseQuery

	var total int

	err := db.DB.QueryRow(
		ctx,
		countQuery,
		args...,
	).Scan(&total)

	if err != nil {

		return PaginatedPlayers{}, err
	}

	// ----------------------------------------------------------
	// Main Query
	// ----------------------------------------------------------

	mainQuery := `
		SELECT
			p.user_id,
			u.first_name,
			u.last_name,
			u.email,
			p.status,
			p.profile_completed,
			COALESCE(ppi.gender::text, ''),
			d.id,
			d.code,
			d.display_name,
			coach_user.id,
			coach_user.first_name,
			coach_user.last_name,
			p.joined_at
	` + baseQuery + `
		ORDER BY p.` + query.SortBy + ` ` + query.OrderBy

	if query.Limit > 0 {

		mainQuery += `
			LIMIT $` + strconv.Itoa(argPos) + `
			OFFSET $` + strconv.Itoa(argPos+1)

		args = append(
			args,
			query.Limit,
			offset,
		)
	}

	rows, err := db.DB.Query(
		ctx,
		mainQuery,
		args...,
	)

	if err != nil {

		return PaginatedPlayers{}, err
	}

	defer rows.Close()

	for rows.Next() {

		var playerItem PlayerListItemResponse

		var disciplineID *int16
		var disciplineCode *string
		var disciplineName *string

		var coachID *string
		var coachFirstName *string
		var coachLastName *string

		err := rows.Scan(
			&playerItem.ID,
			&playerItem.FirstName,
			&playerItem.LastName,
			&playerItem.Email,
			&playerItem.Status,
			&playerItem.ProfileCompleted,
			&playerItem.Gender,
			&disciplineID,
			&disciplineCode,
			&disciplineName,
			&coachID,
			&coachFirstName,
			&coachLastName,
			&playerItem.JoinedAt,
		)

		if err != nil {

			return PaginatedPlayers{}, err
		}

		if disciplineID != nil {

			playerItem.PrimaryDiscipline =
				&profile.Discipline{
					ID:          *disciplineID,
					Code:        *disciplineCode,
					DisplayName: *disciplineName,
				}
		}

		if coachID != nil {

			playerItem.CurrentCoach =
				&profile.CoachSummary{
					UserID: *coachID,
					FullName: strings.TrimSpace(
						*coachFirstName + " " + *coachLastName,
					),
				}
		}

		players = append(
			players,
			playerItem,
		)
	}

	if err := rows.Err(); err != nil {

		return PaginatedPlayers{}, err
	}

	// ----------------------------------------------------------
	// Pagination Metadata
	// ----------------------------------------------------------

	totalPages := 1

	if query.Limit > 0 {

		totalPages = int(
			math.Ceil(
				float64(total) /
					float64(query.Limit),
			),
		)
	}

	return PaginatedPlayers{
		Items:       players,
		Page:        query.Page,
		Limit:       query.Limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     query.Page < totalPages,
		HasPrevious: query.Page > 1,
	}, nil
}

func GetAcademyPlayerRepository(
	ctx context.Context,
	academyID string,
	playerID string,
) (profile.PlayerProfileResponse, error) {

	var playerProfile profile.PlayerProfileResponse

	query := `
		SELECT
			p.profile_completed,
			p.dpdp_consent,
			p.status,
			p.joined_at,

			a.id,
			a.name,
			a.address,
			a.pincode_id,
			pc.code,
			d.name,
			s.name,

			ac.user_id,
			coach_user.first_name,
			coach_user.last_name,
			ac.coach_code

		FROM players p

		INNER JOIN academies a
			ON a.id = p.academy_id

		INNER JOIN pincodes pc
			ON pc.id = a.pincode_id

		INNER JOIN districts d
			ON d.id = pc.district_id

		INNER JOIN states s
			ON s.id = d.state_id

		LEFT JOIN academy_coaches ac
			ON ac.user_id = p.current_coach_user_id

		LEFT JOIN users coach_user
			ON coach_user.id = ac.user_id

		WHERE
			p.user_id = $1
			AND p.academy_id = $2
	`

	var coachID *string
	var coachFirstName *string
	var coachLastName *string
	var coachCode *string

	err := db.DB.QueryRow(
		ctx,
		query,
		playerID,
		academyID,
	).Scan(
		&playerProfile.ProfileCompleted,
		&playerProfile.DPDPConsent,
		&playerProfile.Status,
		&playerProfile.JoinedAt,

		&playerProfile.Academy.ID,
		&playerProfile.Academy.Name,
		&playerProfile.Academy.Address,
		&playerProfile.Academy.PincodeID,
		&playerProfile.Academy.Pincode,
		&playerProfile.Academy.District,
		&playerProfile.Academy.State,

		&coachID,
		&coachFirstName,
		&coachLastName,
		&coachCode,
	)

	if err != nil {

		return profile.PlayerProfileResponse{}, err
	}

	if coachID != nil {

		playerProfile.CurrentCoach =
			&profile.CoachSummary{
				UserID: *coachID,
				FullName: strings.TrimSpace(
					*coachFirstName + " " + *coachLastName,
				),
				CoachCode: *coachCode,
			}
	}

	personalInfo, err :=
		profile.GetPlayerPersonalInfoByUserID(
			ctx,
			playerID,
		)

	if err == nil {

		playerProfile.PersonalInfo = &personalInfo
	}

	sportsProfile, err :=
		profile.GetPlayerSportsProfileByUserID(
			ctx,
			playerID,
		)

	if err == nil {

		playerProfile.SportsProfile = &sportsProfile
	}

	disciplines, err :=
		profile.GetPlayerDisciplinesByUserID(
			ctx,
			playerID,
		)

	if err == nil {

		playerProfile.Disciplines = disciplines
	}

	passport, err :=
		profile.GetPlayerPassportByUserID(
			ctx,
			playerID,
		)

	if err == nil {

		playerProfile.Passport = &passport
	}

	guardians, err :=
		profile.GetPlayerGuardiansByUserID(
			ctx,
			playerID,
		)

	if err == nil {

		playerProfile.Guardians = guardians
	}

	return playerProfile, nil
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

