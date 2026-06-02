package profile

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cobackend/internal/db"
	"cobackend/internal/shared"

	"github.com/jackc/pgx/v5"
)

func GetProfileCompletedStatus(
	ctx context.Context,
	userID string,
	roleName string,
) (bool, error) {

	var query string

	switch roleName {

	case "super_admin":
		return true, nil

	case "state_admin":

		query = `
			SELECT profile_completed
			FROM state_admins
			WHERE user_id = $1
		`

	case "district_admin":

		query = `
			SELECT profile_completed
			FROM district_admins
			WHERE user_id = $1
		`

	case "district_coach":

		query = `
			SELECT profile_completed
			FROM district_coaches
			WHERE user_id = $1
		`

	case "academy_admin":

		query = `
			SELECT profile_completed
			FROM academy_admins
			WHERE user_id = $1
		`

	case "academy_coach":

		query = `
			SELECT profile_completed
			FROM academy_coaches
			WHERE user_id = $1
		`

	case "player":

		query = `
			SELECT profile_completed
			FROM players
			WHERE user_id = $1
		`

	default:
		return false, shared.ErrInvalidRole
	}

	var profileCompleted bool

	err := db.DB.QueryRow(
		ctx,
		query,
		userID,
	).Scan(&profileCompleted)

	if err != nil {

		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {

			// --------------------------------------
			// No profile yet means incomplete
			// --------------------------------------

			return false, nil
		}

		fmt.Println(
			"GetProfileCompletedStatus:",
			err,
		)

		return false, err
	}

	return profileCompleted, nil
}

// func CheckEmailExists(
// 	ctx context.Context,
// 	email string,
// ) (bool, error) {

// 	var exists bool

// 	err := db.DB.QueryRow(
// 		ctx,
// 		`
// 		SELECT EXISTS (
// 			SELECT 1
// 			FROM profiles
// 			WHERE email = $1
// 		)
// 		`,
// 		email,
// 	).Scan(&exists)

// 	if err != nil {
// 		return false, err
// 	}

// 	return exists, nil
// }

// func CreateProfileTx(
// 	ctx context.Context,
// 	tx pgx.Tx,
// 	input CreateProfileInput,
// ) (string, error) {

// 	var profileID string

// 	err := tx.QueryRow(
// 		ctx,
// 		`
// 		INSERT INTO profiles (
// 			first_name,
// 			last_name,
// 			email,
// 			password_hash,
// 			contact_number,
// 			role_id
// 		)
// 		VALUES (
// 			$1,
// 			$2,
// 			$3,
// 			$4,
// 			$5,
// 			$6
// 		)
// 		RETURNING id
// 		`,
// 		input.FirstName,
// 		input.LastName,
// 		input.Email,
// 		input.PasswordHash,
// 		input.ContactNumber,
// 		input.RoleID,
// 	).Scan(&profileID)

// 	if err != nil {
// 		return "", err
// 	}

// 	return profileID, nil
// }

// func GetProfileByID(
// 	ctx context.Context,
// 	profileID string,
// ) (ProfileResponse, error) {

// 	var profile ProfileResponse

// 	err := db.DB.QueryRow(
// 		ctx,
// 		`
// 		SELECT
// 			p.id,
// 			p.first_name,
// 			p.last_name,
// 			p.email,
// 			p.contact_number,
// 			p.role_id,
// 			r.name,
// 			p.created_at
// 		FROM profiles p
// 		JOIN roles r
// 			ON r.id = p.role_id
// 		WHERE p.id = $1
// 		`,
// 		profileID,
// 	).Scan(
// 		&profile.ID,
// 		&profile.FirstName,
// 		&profile.LastName,
// 		&profile.Email,
// 		&profile.ContactNumber,
// 		&profile.RoleID,
// 		&profile.RoleName,
// 		&profile.CreatedAt,
// 	)

// 	if err != nil {
// 		return ProfileResponse{}, err
// 	}

// 	return profile, nil
// }


// Get Profile
// State Admin
func GetStateAdminProfileByUserID(
	ctx context.Context,
	userID string,
) (StateAdminProfileResponse, error) {

	query := `
		SELECT
			sa.profile_completed,
			sa.dpdp_consent,
			s.id,
			s.name

		FROM state_admins sa

		INNER JOIN states s
			ON s.id = sa.state_id

		WHERE sa.user_id = $1
	`

	var profile StateAdminProfileResponse

	err := db.DB.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&profile.ProfileCompleted,
		&profile.DPDPConsent,
		&profile.StateID,
		&profile.StateName,
	)

	if err != nil {

		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return StateAdminProfileResponse{},
				shared.ErrProfileNotFound
		}

		return StateAdminProfileResponse{}, err
	}

	return profile, nil
}

// District Admin
func GetDistrictAdminProfileByUserID(
	ctx context.Context,
	userID string,
) (DistrictAdminProfileResponse, error) {

	query := `
		SELECT
			da.profile_completed,
			da.dpdp_consent,

			d.id,
			d.name,

			s.id,
			s.name

		FROM district_admins da

		INNER JOIN districts d
			ON d.id = da.district_id

		INNER JOIN states s
			ON s.id = d.state_id

		WHERE da.user_id = $1
	`

	var profile DistrictAdminProfileResponse

	err := db.DB.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&profile.ProfileCompleted,
		&profile.DPDPConsent,

		&profile.DistrictID,
		&profile.DistrictName,

		&profile.StateID,
		&profile.StateName,
	)

	if err != nil {

		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return DistrictAdminProfileResponse{},
				shared.ErrProfileNotFound
		}

		return DistrictAdminProfileResponse{}, err
	}

	return profile, nil
}

// District Coach
func GetDistrictCoachProfileByUserID(
	ctx context.Context,
	userID string,
) (DistrictCoachProfileResponse, error) {

	query := `
		SELECT
			dc.profile_completed,
			dc.dpdp_consent,
			dc.coach_code,
			dc.coaching_certificate_proof,

			d.id,
			d.name,

			s.id,
			s.name

		FROM district_coaches dc

		INNER JOIN districts d
			ON d.id = dc.district_id

		INNER JOIN states s
			ON s.id = d.state_id

		WHERE dc.user_id = $1
	`

	var profile DistrictCoachProfileResponse

	err := db.DB.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&profile.ProfileCompleted,
		&profile.DPDPConsent,
		&profile.CoachCode,
		&profile.CoachingCertificateProof,

		&profile.DistrictID,
		&profile.DistrictName,

		&profile.StateID,
		&profile.StateName,
	)

	if err != nil {

		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return DistrictCoachProfileResponse{},
				shared.ErrProfileNotFound
		}

		return DistrictCoachProfileResponse{}, err
	}

	// ----------------------------------------------------------
	// Fetch Disciplines
	// ----------------------------------------------------------

	disciplines, err := GetDistrictCoachDisciplinesByUserID(
		ctx,
		userID,
	)

	if err != nil {
		return DistrictCoachProfileResponse{}, err
	}

	profile.Disciplines = disciplines

	return profile, nil
}

// District Coach Disciplines
func GetDistrictCoachDisciplinesByUserID(
	ctx context.Context,
	userID string,
) ([]Discipline, error) {

	query := `
		SELECT
			d.id,
			d.code,
			d.display_name

		FROM district_coach_disciplines dcd

		INNER JOIN disciplines d
			ON d.id = dcd.discipline_id

		WHERE dcd.coach_user_id = $1

		ORDER BY d.display_name
	`

	rows, err := db.DB.Query(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	disciplines := []Discipline{}

	for rows.Next() {

		var discipline Discipline

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

// Academy Admin
// Academy Admin
func GetAcademyAdminProfileByUserID(
	ctx context.Context,
	userID string,
) (AcademyAdminProfileResponse, error) {

	query := `
		SELECT
			aa.profile_completed,
			aa.dpdp_consent,
			aa.gstin,
			aa.registration_proof,

			a.id,
			a.name,
			a.address,

			p.id,
			p.code,

			d.id,
			d.name,

			s.id,
			s.name

		FROM academy_admins aa

		INNER JOIN academies a
			ON a.id = aa.academy_id

		INNER JOIN pincodes p
			ON p.id = a.pincode_id

		INNER JOIN districts d
			ON d.id = p.district_id

		INNER JOIN states s
			ON s.id = d.state_id

		WHERE aa.user_id = $1
	`

	var profile AcademyAdminProfileResponse

	err := db.DB.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&profile.ProfileCompleted,
		&profile.DPDPConsent,
		&profile.GSTIN,
		&profile.RegistrationProof,

		&profile.AcademyID,
		&profile.AcademyName,
		&profile.AcademyAddress,

		&profile.PincodeID,
		&profile.Pincode,

		&profile.DistrictID,
		&profile.DistrictName,

		&profile.StateID,
		&profile.StateName,
	)

	if err != nil {

		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return AcademyAdminProfileResponse{},
				shared.ErrProfileNotFound
		}

		return AcademyAdminProfileResponse{}, err
	}

	return profile, nil
}

// Academy Coach 
func GetAcademyCoachProfileByUserID(
	ctx context.Context,
	userID string,
) (AcademyCoachProfileResponse, error) {

	query := `
		SELECT
			ac.profile_completed,
			ac.dpdp_consent,
			ac.coach_code,
			ac.coaching_certificate_proof,

			a.id,
			a.name,
			a.address,

			d.id,
			d.name,

			s.id,
			s.name

		FROM academy_coaches ac

		INNER JOIN academies a
			ON a.id = ac.academy_id

		INNER JOIN districts d
			ON d.id = a.district_id

		INNER JOIN states s
			ON s.id = d.state_id

		WHERE ac.user_id = $1
	`

	var profile AcademyCoachProfileResponse

	err := db.DB.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&profile.ProfileCompleted,
		&profile.DPDPConsent,
		&profile.CoachCode,
		&profile.CoachingCertificateProof,

		&profile.AcademyID,
		&profile.AcademyName,
		&profile.AcademyAddress,

		&profile.DistrictID,
		&profile.DistrictName,

		&profile.StateID,
		&profile.StateName,
	)

	if err != nil {

		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return AcademyCoachProfileResponse{},
				shared.ErrProfileNotFound
		}

		return AcademyCoachProfileResponse{}, err
	}

	// ----------------------------------------------------------
	// Fetch Disciplines
	// ----------------------------------------------------------

	disciplines, err := GetAcademyCoachDisciplinesByUserID(
		ctx,
		userID,
	)

	if err != nil {
		return AcademyCoachProfileResponse{}, err
	}

	profile.Disciplines = disciplines

	return profile, nil
}

// Academy Coach Disciplines
func GetAcademyCoachDisciplinesByUserID(
	ctx context.Context,
	userID string,
) ([]Discipline, error) {

	query := `
		SELECT
			d.id,
			d.code,
			d.display_name

		FROM academy_coach_disciplines acd

		INNER JOIN disciplines d
			ON d.id = acd.discipline_id

		WHERE acd.coach_user_id = $1

		ORDER BY d.display_name
	`

	rows, err := db.DB.Query(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	disciplines := []Discipline{}

	for rows.Next() {

		var discipline Discipline

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

// Player
func GetPlayerProfileByUserID(
	ctx context.Context,
	userID string,
) (PlayerProfileResponse, error) {

	query := `
		SELECT
			p.profile_completed,
			p.dpdp_consent,
			p.status,
			p.joined_at,
			p.current_coach_user_id,

			a.id,
			a.name,
			a.address,

			d.name,
			s.name

		FROM players p

		INNER JOIN academies a
			ON a.id = p.academy_id

		INNER JOIN districts d
			ON d.id = a.district_id

		INNER JOIN states s
			ON s.id = d.state_id

		WHERE p.user_id = $1
	`

	var profile PlayerProfileResponse

	profile.Academy = AcademySummary{}

	err := db.DB.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&profile.ProfileCompleted,
		&profile.DPDPConsent,
		&profile.Status,
		&profile.JoinedAt,
		&profile.CurrentCoachUserID,

		&profile.Academy.ID,
		&profile.Academy.Name,
		&profile.Academy.Address,

		&profile.Academy.District,
		&profile.Academy.State,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return PlayerProfileResponse{},
				shared.ErrProfileNotFound
		}

		return PlayerProfileResponse{}, err
	}

	// ----------------------------------------------------------
	// Personal Info
	// ----------------------------------------------------------

	personalInfo, err := GetPlayerPersonalInfoByUserID(
		ctx,
		userID,
	)

	if err == nil {
		profile.PersonalInfo = &personalInfo
	}

	// ----------------------------------------------------------
	// Sports Profile
	// ----------------------------------------------------------

	sportsProfile, err := GetPlayerSportsProfileByUserID(
		ctx,
		userID,
	)

	if err == nil {
		profile.SportsProfile = &sportsProfile
	}

	// ----------------------------------------------------------
	// Disciplines
	// ----------------------------------------------------------

	disciplines, err := GetPlayerDisciplinesByUserID(
		ctx,
		userID,
	)

	if err != nil {
		return PlayerProfileResponse{}, err
	}

	profile.Disciplines = disciplines

	// ----------------------------------------------------------
	// Passport
	// ----------------------------------------------------------

	passport, err := GetPlayerPassportByUserID(
		ctx,
		userID,
	)

	if err == nil {
		profile.Passport = &passport
	}

	// ----------------------------------------------------------
	// Guardians
	// ----------------------------------------------------------

	guardians, err := GetPlayerGuardiansByUserID(
		ctx,
		userID,
	)

	if err != nil {
		return PlayerProfileResponse{}, err
	}

	if guardians == nil {
		guardians = []PlayerGuardian{}
	}

	profile.Guardians = guardians

	return profile, nil
}

func GetPlayerPersonalInfoByUserID(
	ctx context.Context,
	userID string,
) (PlayerPersonalInfo, error) {

	query := `
		SELECT
			ppi.date_of_birth,
			ppi.gender,
			ppi.nationality,
			ppi.place_of_birth,
			ppi.city,
			ppi.residential_address,

			p.id,
			p.code,
			d.name,
			s.name,

			ppi.education,
			ppi.institution_name,
			ppi.occupation,
			ppi.temporary_sport_id

		FROM player_personal_info ppi

		LEFT JOIN pincodes p
			ON p.id = ppi.pincode_id

		LEFT JOIN districts d
			ON d.id = p.district_id

		LEFT JOIN states s
			ON s.id = d.state_id

		WHERE ppi.player_user_id = $1
	`

	var info PlayerPersonalInfo

	var (
		pincodeID       *int32
		pincodeCode     *string
		pincodeDistrict *string
		pincodeState    *string
	)

	err := db.DB.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&info.DateOfBirth,
		&info.Gender,
		&info.Nationality,
		&info.PlaceOfBirth,
		&info.City,
		&info.ResidentialAddress,

		&pincodeID,
		&pincodeCode,
		&pincodeDistrict,
		&pincodeState,

		&info.Education,
		&info.InstitutionName,
		&info.Occupation,
		&info.TemporarySportID,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return PlayerPersonalInfo{},
				shared.ErrProfileNotFound
		}

		return PlayerPersonalInfo{}, err
	}

	// ----------------------------------------------------------
	// Build Pincode Response
	// ----------------------------------------------------------

	if pincodeID != nil {

		info.Pincode = &PincodeInfo{
			ID: *pincodeID,
		}

		if pincodeCode != nil {
			info.Pincode.Code = *pincodeCode
		}

		if pincodeDistrict != nil {
			info.Pincode.District = *pincodeDistrict
		}

		if pincodeState != nil {
			info.Pincode.State = *pincodeState
		}
	}

	return info, nil
}

func GetPlayerSportsProfileByUserID(
	ctx context.Context,
	userID string,
) (PlayerSportsProfile, error) {

	query := `
		SELECT
			psp.unit_of_representation,
			psp.dominant_hand,
			psp.height_cm,
			psp.weight_kg,
			psp.shoe_size,
			psp.tracksuit_size

		FROM player_sports_profile psp

		WHERE psp.player_user_id = $1
	`

	var profile PlayerSportsProfile

	err := db.DB.QueryRow(
		ctx,
		query,
		userID,
	).Scan(

		&profile.UnitOfRepresentation,
		&profile.DominantHand,
		&profile.HeightCM,
		&profile.WeightKG,
		&profile.ShoeSize,
		&profile.TracksuitSize,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return PlayerSportsProfile{},
				shared.ErrProfileNotFound
		}

		return PlayerSportsProfile{}, err
	}

	return profile, nil
}

func GetPlayerDisciplinesByUserID(
	ctx context.Context,
	userID string,
) ([]Discipline, error) {

	query := `
		SELECT
			d.id,
			d.code,
			d.display_name,
			pd.is_primary

		FROM player_disciplines pd

		INNER JOIN disciplines d
			ON d.id = pd.discipline_id

		WHERE pd.player_user_id = $1

		ORDER BY
			pd.is_primary DESC,
			d.display_name ASC
	`

	rows, err := db.DB.Query(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	disciplines := []Discipline{}

	for rows.Next() {

		var discipline Discipline

		err := rows.Scan(
			&discipline.ID,
			&discipline.Code,
			&discipline.DisplayName,
			&discipline.IsPrimary,
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

func GetPlayerPassportByUserID(
	ctx context.Context,
	userID string,
) (PlayerPassport, error) {

	query := `
		SELECT
			passport_number,
			passport_issue_date,
			passport_expiry_date,
			passport_issuing_authority,
			passport_place_of_issue

		FROM player_passports

		WHERE player_user_id = $1
	`

	var passport PlayerPassport

	err := db.DB.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&passport.PassportNumber,
		&passport.PassportIssueDate,
		&passport.PassportExpiryDate,
		&passport.PassportIssuingAuthority,
		&passport.PassportPlaceOfIssue,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return PlayerPassport{},
				shared.ErrProfileNotFound
		}

		return PlayerPassport{}, err
	}

	return passport, nil
}

func GetPlayerGuardiansByUserID(
	ctx context.Context,
	userID string,
) ([]PlayerGuardian, error) {

	query := `
		SELECT
			id,
			full_name,
			relationship,
			contact_number,
			alternative_contact,
			parental_consent,
			is_primary

		FROM player_guardians

		WHERE player_user_id = $1

		ORDER BY
			is_primary DESC,
			created_at ASC
	`

	rows, err := db.DB.Query(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var guardians []PlayerGuardian

	for rows.Next() {

		var guardian PlayerGuardian

		err := rows.Scan(
			&guardian.ID,
			&guardian.FullName,
			&guardian.Relationship,
			&guardian.ContactNumber,
			&guardian.AlternativeContact,
			&guardian.ParentalConsent,
			&guardian.IsPrimary,
		)

		if err != nil {
			return nil, err
		}

		guardians = append(
			guardians,
			guardian,
		)
	}

	return guardians, nil
}



// Complete Profile

func CompleteStateAdminProfile(
	ctx context.Context,
	userID string,
	input CompleteStateAdminProfileInput,
) error {

	_, err := db.DB.Exec(
		ctx,
		`
		UPDATE state_admins
		SET
			dpdp_consent = $1,
			profile_completed = true
		WHERE user_id = $2
		`,
		input.DPDPConsent,
		userID,
	)

	return err
}

func CompleteDistrictAdminProfile(
	ctx context.Context,
	userID string,
	input CompleteDistrictAdminProfileInput,
) error {

	_, err := db.DB.Exec(
		ctx,
		`
		UPDATE district_admins
		SET
			dpdp_consent = $1,
			profile_completed = true,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $2
		`,
		input.DPDPConsent,
		userID,
	)

	return err
}

func CompleteDistrictCoachProfile(
	ctx context.Context,
	userID string,
	input CompleteDistrictCoachProfileInput,
) error {

	tx, err := db.DB.Begin(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	// ----------------------------------------------------------
	// Determine Profile Completion
	// ----------------------------------------------------------

	profileCompleted :=
		input.DPDPConsent &&
			strings.TrimSpace(input.CoachCode) != "" &&
			strings.TrimSpace(input.CoachingCertificateProof) != "" &&
			len(input.DisciplineIDs) > 0

	// ----------------------------------------------------------
	// Update District Coach
	// ----------------------------------------------------------

	_, err = tx.Exec(
		ctx,
		`
		UPDATE district_coaches
		SET
			dpdp_consent = $1,
			coach_code = $2,
			coaching_certificate_proof = $3,
			profile_completed = $4
		WHERE user_id = $5
		`,
		input.DPDPConsent,
		input.CoachCode,
		input.CoachingCertificateProof,
		profileCompleted,
		userID,
	)

	if err != nil {
		return err
	}

	// ----------------------------------------------------------
	// Clear Existing Disciplines
	// ----------------------------------------------------------

	_, err = tx.Exec(
		ctx,
		`
		DELETE FROM district_coach_disciplines
		WHERE coach_user_id = $1
		`,
		userID,
	)

	if err != nil {
		return err
	}

	// ----------------------------------------------------------
	// Insert Disciplines
	// ----------------------------------------------------------

	for _, disciplineID := range input.DisciplineIDs {
		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO district_coach_disciplines (
				coach_user_id,
				discipline_id
			)
			VALUES ($1, $2)
			`,
			userID,
			disciplineID,
		)

		if err != nil {
			return err
		}
	}

	// ----------------------------------------------------------
	// Commit Transaction
	// ----------------------------------------------------------

	err = tx.Commit(ctx)

	if err != nil {
		return err
	}

	return nil
}

func CompleteAcademyAdminProfile(
	ctx context.Context,
	userID string,
	input CompleteAcademyAdminProfileInput,
) error {

	profileCompleted :=
		input.DPDPConsent &&
			strings.TrimSpace(input.GSTIN) != "" &&
			strings.TrimSpace(input.RegistrationProof) != ""

	_, err := db.DB.Exec(
		ctx,
		`
		UPDATE academy_admins
		SET
			dpdp_consent = $1,
			gstin = $2,
			registration_proof = $3,
			profile_completed = $4
		WHERE user_id = $5
		`,
		input.DPDPConsent,
		input.GSTIN,
		input.RegistrationProof,
		profileCompleted,
		userID,
	)

	return err
}

func CompleteAcademyCoachProfile(
	ctx context.Context,
	userID string,
	input CompleteAcademyCoachProfileInput,
) error {

	tx, err := db.DB.Begin(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	// ----------------------------------------------------------
	// Determine Profile Completion
	// ----------------------------------------------------------

	profileCompleted :=
		input.DPDPConsent &&
			strings.TrimSpace(input.CoachCode) != "" &&
			strings.TrimSpace(input.CoachingCertificateProof) != "" &&
			len(input.DisciplineIDs) > 0

	// ----------------------------------------------------------
	// Update Academy Coach
	// ----------------------------------------------------------

	_, err = tx.Exec(
		ctx,
		`
		UPDATE academy_coaches
		SET
			dpdp_consent = $1,
			coach_code = $2,
			coaching_certificate_proof = $3,
			profile_completed = $4
		WHERE user_id = $5
		`,
		input.DPDPConsent,
		input.CoachCode,
		input.CoachingCertificateProof,
		profileCompleted,
		userID,
	)

	if err != nil {
		return err
	}

	// ----------------------------------------------------------
	// Clear Existing Disciplines
	// ----------------------------------------------------------

	_, err = tx.Exec(
		ctx,
		`
		DELETE FROM academy_coach_disciplines
		WHERE coach_user_id = $1
		`,
		userID,
	)

	if err != nil {
		return err
	}

	// ----------------------------------------------------------
	// Insert Disciplines
	// ----------------------------------------------------------

	for _, disciplineID := range input.DisciplineIDs {

		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO academy_coach_disciplines (
				coach_user_id,
				discipline_id
			)
			VALUES ($1, $2)
			`,
			userID,
			disciplineID,
		)

		if err != nil {
			return err
		}
	}

	// ----------------------------------------------------------
	// Commit Transaction
	// ----------------------------------------------------------

	err = tx.Commit(ctx)

	if err != nil {
		return err
	}

	return nil
}

func CompletePlayerProfile(
	ctx context.Context,
	userID string,
	input CompletePlayerProfileInput,
) error {

	tx, err := db.DB.Begin(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	// ----------------------------------------------------------
	// Determine Profile Completion
	// ----------------------------------------------------------

	profileCompleted :=
	input.DPDPConsent &&
		!input.PersonalInfo.DateOfBirth.IsZero() &&
		strings.TrimSpace(input.PersonalInfo.Gender) != "" &&
		strings.TrimSpace(input.PersonalInfo.Nationality) != "" &&
		len(input.Disciplines) > 0 &&
		len(input.Guardians) > 0

	// ----------------------------------------------------------
	// Update Players
	// ----------------------------------------------------------

	_, err = tx.Exec(
		ctx,
		`
		UPDATE players
		SET
			dpdp_consent = $1,
			profile_completed = $2
		WHERE user_id = $3
		`,
		input.DPDPConsent,
		profileCompleted,
		userID,
	)

	if err != nil {
		return err
	}

	// ----------------------------------------------------------
	// Upsert Personal Info
	// ----------------------------------------------------------

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO player_personal_info (
			player_user_id,
			date_of_birth,
			gender,
			nationality,
			place_of_birth,
			city,
			residential_address,
			pincode_id,
			education,
			institution_name,
			occupation,
			temporary_sport_id
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12
		)

		ON CONFLICT (player_user_id)
		DO UPDATE SET
			date_of_birth = EXCLUDED.date_of_birth,
			gender = EXCLUDED.gender,
			nationality = EXCLUDED.nationality,
			place_of_birth = EXCLUDED.place_of_birth,
			city = EXCLUDED.city,
			residential_address = EXCLUDED.residential_address,
			pincode_id = EXCLUDED.pincode_id,
			education = EXCLUDED.education,
			institution_name = EXCLUDED.institution_name,
			occupation = EXCLUDED.occupation,
			temporary_sport_id = EXCLUDED.temporary_sport_id
		`,
		userID,
		input.PersonalInfo.DateOfBirth,
		input.PersonalInfo.Gender,
		input.PersonalInfo.Nationality,
		input.PersonalInfo.PlaceOfBirth,
		input.PersonalInfo.City,
		input.PersonalInfo.ResidentialAddress,
		input.PersonalInfo.PincodeID,
		input.PersonalInfo.Education,
		input.PersonalInfo.InstitutionName,
		input.PersonalInfo.Occupation,
		input.PersonalInfo.TemporarySportID,
	)

	if err != nil {
		return err
	}

	// ----------------------------------------------------------
	// Upsert Sports Profile
	// ----------------------------------------------------------

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO player_sports_profile (
			player_user_id,
			unit_of_representation,
			dominant_hand,
			height_cm,
			weight_kg,
			shoe_size,
			tracksuit_size
		)
		VALUES (
			$1, $2, $3, $4,
			$5, $6, $7
		)

		ON CONFLICT (player_user_id)
		DO UPDATE SET
			unit_of_representation = EXCLUDED.unit_of_representation,
			dominant_hand = EXCLUDED.dominant_hand,
			height_cm = EXCLUDED.height_cm,
			weight_kg = EXCLUDED.weight_kg,
			shoe_size = EXCLUDED.shoe_size,
			tracksuit_size = EXCLUDED.tracksuit_size
		`,
		userID,
		input.SportsProfile.UnitOfRepresentation,
		input.SportsProfile.DominantHand,
		input.SportsProfile.HeightCM,
		input.SportsProfile.WeightKG,
		input.SportsProfile.ShoeSize,
		input.SportsProfile.TracksuitSize,
	)

	if err != nil {
		return err
	}

	// ----------------------------------------------------------
	// Replace Disciplines
	// ----------------------------------------------------------

	_, err = tx.Exec(
		ctx,
		`
		DELETE FROM player_disciplines
		WHERE player_user_id = $1
		`,
		userID,
	)

	if err != nil {
		return err
	}

	for _, discipline := range input.Disciplines {

		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO player_disciplines (
				player_user_id,
				discipline_id,
				is_primary
			)
			VALUES ($1, $2, $3)
			`,
			userID,
			discipline.DisciplineID,
			discipline.IsPrimary,
		)

		if err != nil {
			return err
		}
	}

	// ----------------------------------------------------------
	// Replace Guardians
	// ----------------------------------------------------------

	_, err = tx.Exec(
		ctx,
		`
		DELETE FROM player_guardians
		WHERE player_user_id = $1
		`,
		userID,
	)

	if err != nil {
		return err
	}

	for _, guardian := range input.Guardians {

		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO player_guardians (
				player_user_id,
				full_name,
				relationship,
				contact_number,
				alternative_contact,
				parental_consent,
				is_primary
			)
			VALUES (
				$1, $2, $3, $4,
				$5, $6, $7
			)
			`,
			userID,
			guardian.FullName,
			guardian.Relationship,
			guardian.ContactNumber,
			guardian.AlternativeContact,
			guardian.ParentalConsent,
			guardian.IsPrimary,
		)

		if err != nil {
			return err
		}
	}

	// ----------------------------------------------------------
	// Upsert Passport (Optional)
	// ----------------------------------------------------------

	if input.Passport != nil {

		_, err = tx.Exec(
			ctx,
			`
			INSERT INTO player_passports (
				player_user_id,
				passport_number,
				passport_issue_date,
				passport_expiry_date,
				passport_issuing_authority,
				passport_place_of_issue
			)
			VALUES (
				$1, $2, $3,
				$4, $5, $6
			)

			ON CONFLICT (player_user_id)
			DO UPDATE SET
				passport_number = EXCLUDED.passport_number,
				passport_issue_date = EXCLUDED.passport_issue_date,
				passport_expiry_date = EXCLUDED.passport_expiry_date,
				passport_issuing_authority = EXCLUDED.passport_issuing_authority,
				passport_place_of_issue = EXCLUDED.passport_place_of_issue
			`,
			userID,
			input.Passport.PassportNumber,
			input.Passport.PassportIssueDate,
			input.Passport.PassportExpiryDate,
			input.Passport.PassportIssuingAuthority,
			input.Passport.PassportPlaceOfIssue,
		)

		if err != nil {
			return err
		}
	}

	// ----------------------------------------------------------
	// Commit
	// ----------------------------------------------------------

	err = tx.Commit(ctx)

	if err != nil {
		return err
	}

	return nil
}