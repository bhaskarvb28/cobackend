package academyCoach

import (
	"context"
	"math"
	"strconv"

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

func GetAcademyCoaches(
	ctx context.Context,
	academyID string,
	query GetAcademyCoachesQuery,
) (
	PaginatedAcademyCoachesResponse,
	error,
) {

	offset := 0

	if query.Limit > 0 {

		offset =
			(query.Page - 1) *
				query.Limit
	}

	baseQuery := `
	SELECT
		u.id,
		u.first_name,
		u.last_name,

		COALESCE(
			ac.coach_code,
			''
		),

		COUNT(p.user_id) AS assigned_players_count

	FROM academy_coaches ac

	INNER JOIN users u
		ON u.id = ac.user_id

	LEFT JOIN players p
		ON p.current_coach_user_id = ac.user_id

	WHERE ac.academy_id = $1
	`

	countQuery := `
	SELECT COUNT(*)

	FROM academy_coaches ac

	INNER JOIN users u
		ON u.id = ac.user_id

	WHERE ac.academy_id = $1
	`

	args := []interface{}{
		academyID,
	}

	argPos := 2

	// ----------------------------------------------------------
	// Search
	// ----------------------------------------------------------

	if query.Search != "" {

		searchCondition := `
		AND (
			u.first_name ILIKE $` +
			strconv.Itoa(argPos) + `
			OR u.last_name ILIKE $` +
			strconv.Itoa(argPos) + `
			OR u.email ILIKE $` +
			strconv.Itoa(argPos) + `
		)
		`

		baseQuery += searchCondition
		countQuery += searchCondition

		args = append(
			args,
			"%"+query.Search+"%",
		)

		argPos++
	}

	baseQuery += `
	GROUP BY
		u.id,
		u.first_name,
		u.last_name,
		ac.coach_code
	`

	baseQuery += `
	ORDER BY u.first_name ASC
	`

	if query.Limit > 0 {

		baseQuery += `
		LIMIT $` +
			strconv.Itoa(argPos) +
			` OFFSET $` +
			strconv.Itoa(argPos+1)

		args = append(
			args,
			query.Limit,
			offset,
		)
	}

	var total int

	err := db.DB.QueryRow(
		ctx,
		countQuery,
		args[:argPos-1]...,
	).Scan(&total)

	if err != nil {

		return PaginatedAcademyCoachesResponse{},
			err
	}

	rows, err := db.DB.Query(
		ctx,
		baseQuery,
		args...,
	)

	if err != nil {

		return PaginatedAcademyCoachesResponse{},
			err
	}

	defer rows.Close()

	items := []CoachSummary{}

	for rows.Next() {

		var item = CoachSummary{
			Disciplines:
				[]DisciplineSummary{},
		}

		var firstName string
		var lastName string

		err := rows.Scan(
			&item.UserID,
			&firstName,
			&lastName,
			&item.CoachCode,
			&item.AssignedPlayersCount,
		)

		if err != nil {

			return PaginatedAcademyCoachesResponse{},
				err
		}

		item.FullName =
			firstName + " " + lastName

		// ------------------------------------------------------
		// Coach Disciplines
		// ------------------------------------------------------

		disciplines, err :=
			GetCoachDisciplines(
				ctx,
				item.UserID,
			)

		if err != nil {

			return PaginatedAcademyCoachesResponse{},
				err
		}

		item.Disciplines =
			disciplines

		items = append(
			items,
			item,
		)
	}

	totalPages := 1

	if query.Limit > 0 {

		totalPages = int(
			math.Ceil(
				float64(total) /
					float64(query.Limit),
			),
		)
	}

	return PaginatedAcademyCoachesResponse{
		Items:       items,
		Page:        query.Page,
		Limit:       query.Limit,
		Total:       total,
		TotalPages:  totalPages,
		HasNext:     query.Page < totalPages,
		HasPrevious: query.Page > 1,
	}, nil
}

func ValidateAcademyCoach(
	ctx context.Context,
	academyID string,
	coachUserID string,
) (
	bool,
	error,
) {

	var exists bool

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT EXISTS(
			SELECT 1
			FROM academy_coaches
			WHERE academy_id = $1
			AND user_id = $2
		)
		`,
		academyID,
		coachUserID,
	).Scan(&exists)

	return exists, err
}

func GetAcademyCoach(
	ctx context.Context,
	academyID string,
	coachUserID string,
) (
	AcademyCoachProfileResponse,
	error,
) {

	var response AcademyCoachProfileResponse

	err := db.DB.QueryRow(
		ctx,
		`
		SELECT
			u.id,
			u.first_name,
			u.last_name,
			u.email,

			COALESCE(
				ac.coach_code,
				''
			),

			ac.created_at,

			COUNT(p.user_id) AS assigned_players_count

		FROM academy_coaches ac

		INNER JOIN users u
			ON u.id = ac.user_id

		LEFT JOIN players p
			ON p.current_coach_user_id = ac.user_id

		WHERE ac.academy_id = $1
		AND ac.user_id = $2

		GROUP BY
			u.id,
			u.first_name,
			u.last_name,
			u.email,
			ac.coach_code,
			ac.created_at
		`,
		academyID,
		coachUserID,
	).Scan(
		&response.UserID,
		&response.FirstName,
		&response.LastName,
		&response.Email,
		&response.CoachCode,
		&response.JoinedAt,
		&response.AssignedPlayersCount,
	)

	if err != nil {

		return AcademyCoachProfileResponse{},
			err
	}

	response.FullName =
		response.FirstName +
			" " +
			response.LastName

	disciplines, err :=
		GetCoachDisciplines(
			ctx,
			coachUserID,
		)

	if err != nil {

		return AcademyCoachProfileResponse{},
			err
	}

	response.Disciplines =
		disciplines

	return response, nil
}

func GetCoachDisciplines(
	ctx context.Context,
	coachUserID string,
) (
	[]DisciplineSummary,
	error,
) {

	rows, err := db.DB.Query(
		ctx,
		`
		SELECT
			d.id,
			d.code,
			d.display_name

		FROM academy_coach_disciplines acd

		INNER JOIN disciplines d
			ON d.id = acd.discipline_id

		WHERE acd.coach_user_id = $1

		ORDER BY d.display_name ASC
		`,
		coachUserID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	disciplines :=
		[]DisciplineSummary{}

	for rows.Next() {

		var discipline DisciplineSummary

		err := rows.Scan(
			&discipline.ID,
			&discipline.Code,
			&discipline.DisplayName,
		)

		if err != nil {
			return nil, err
		}

		disciplines = append(
			disciplines,
			discipline,
		)
	}

	return disciplines, nil
}