package profile

import (
	"context"
	"fmt"
	"encoding/json"
	"net/http"

	"cobackend/internal/auth"
)

func GetProfileService(
	ctx context.Context,
	userID string,
) (ProfileResponse, error) {

	// ----------------------------------------------------------
	// Get Base User
	// ----------------------------------------------------------

	user, err := auth.GetUserByID(
		ctx,
		userID,
	)

	if err != nil {
		return ProfileResponse{}, err
	}

	// ----------------------------------------------------------
	// Build Base Response
	// ----------------------------------------------------------

	response := ProfileResponse{
		User: user,
	}

	// ----------------------------------------------------------
	// Load Role Profile
	// ----------------------------------------------------------

	switch user.Role.Code {

	case "state_admin":

		stateAdminProfile, err := GetStateAdminProfileByUserID(
			ctx,
			userID,
		)

		if err != nil {
			return ProfileResponse{}, err
		}

		response.Profile = stateAdminProfile

	case "district_admin":

		districtAdminProfile, err := GetDistrictAdminProfileByUserID(
			ctx,
			userID,
		)

		if err != nil {
			return ProfileResponse{}, err
		}

		response.Profile = districtAdminProfile

	case "district_coach":

		districtCoachProfile, err := GetDistrictCoachProfileByUserID(
			ctx,
			userID,
		)

		if err != nil {
			
			return ProfileResponse{}, err
		}

		response.Profile = districtCoachProfile

	case "academy_admin":

		academyAdminProfile, err := GetAcademyAdminProfileByUserID(
			ctx,
			userID,
		)

		if err != nil {
			return ProfileResponse{}, err
		}

		response.Profile = academyAdminProfile

	case "academy_coach":

		academyCoachProfile, err := GetAcademyCoachProfileByUserID(
			ctx,
			userID,
		)

		if err != nil {
			return ProfileResponse{}, err
		}

		response.Profile = academyCoachProfile

	case "player":

		playerProfile, err := GetPlayerProfileByUserID(
			ctx,
			userID,
		)

		if err != nil {
			
			return ProfileResponse{}, err
		}

		response.Profile = playerProfile

	default:
		return ProfileResponse{}, fmt.Errorf("unsupported role")
	}

	return response, nil
}

func CompleteProfileService(
	ctx context.Context,
	userID string,
	r *http.Request,
) (ProfileResponse, error) {

	// ----------------------------------------------------------
	// Get Authenticated User
	// ----------------------------------------------------------

	user, err := auth.GetUserByID(
		ctx,
		userID,
	)

	if err != nil {
		return ProfileResponse{}, err
	}

	// ----------------------------------------------------------
	// Complete Role Profile
	// ----------------------------------------------------------

	switch user.Role.Code {

	case "state_admin":

		var input CompleteStateAdminProfileInput

		err := json.NewDecoder(r.Body).
			Decode(&input)

		if err != nil {
			return ProfileResponse{}, err
		}

		err = CompleteStateAdminProfile(
			ctx,
			userID,
			input,
		)

		if err != nil {
			return ProfileResponse{}, err
		}

	case "district_admin":

		var input CompleteDistrictAdminProfileInput

		err := json.NewDecoder(r.Body).
			Decode(&input)

		if err != nil {
			return ProfileResponse{}, err
		}

		err = CompleteDistrictAdminProfile(
			ctx,
			userID,
			input,
		)

		if err != nil {
			return ProfileResponse{}, err
		}

	case "district_coach":

		var input CompleteDistrictCoachProfileInput

		err := json.NewDecoder(r.Body).
			Decode(&input)

		if err != nil {
			return ProfileResponse{}, err
		}

		err = CompleteDistrictCoachProfile(
			ctx,
			userID,
			input,
		)

		if err != nil {
			
			return ProfileResponse{}, err
		}

	case "academy_admin":

		var input CompleteAcademyAdminProfileInput

		err := json.NewDecoder(r.Body).
			Decode(&input)

		if err != nil {
			return ProfileResponse{}, err
		}

		err = CompleteAcademyAdminProfile(
			ctx,
			userID,
			input,
		)

		if err != nil {
			return ProfileResponse{}, err
		}

	case "academy_coach":

		var input CompleteAcademyCoachProfileInput

		err := json.NewDecoder(r.Body).
			Decode(&input)

		if err != nil {
			return ProfileResponse{}, err
		}

		err = CompleteAcademyCoachProfile(
			ctx,
			userID,
			input,
		)

		if err != nil {
			return ProfileResponse{}, err
		}

	case "player":

		var input CompletePlayerProfileInput

		err := json.NewDecoder(r.Body).
			Decode(&input)

		if err != nil {
			return ProfileResponse{}, err
		}

		err = CompletePlayerProfile(
			ctx,
			userID,
			input,
		)

		if err != nil {
			fmt.Print(err)
			return ProfileResponse{}, err
		}

	default:
		return ProfileResponse{},
			fmt.Errorf("unsupported role")
	}

	// ----------------------------------------------------------
	// Return Updated Profile
	// ----------------------------------------------------------

	return GetProfileService(
		ctx,
		userID,
	)
}