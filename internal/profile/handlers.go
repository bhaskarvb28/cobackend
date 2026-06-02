package profile

import (
	"net/http"

	"cobackend/internal/middleware"
	"cobackend/internal/shared"
	"cobackend/internal/utils"

	"fmt"
)

func GetProfileHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Get Authenticated User ID
	// ----------------------------------------------------------

	authUserID, ok := r.Context().
		Value(middleware.UserIDKey).
		(string)

	if !ok {

		utils.WriteJSON(
			w,
			http.StatusUnauthorized,
			shared.APIResponse{
				Success: false,
				Message: "unauthorized",
			},
		)

		return
	}

	// ----------------------------------------------------------
	// Fetch Profile
	// ----------------------------------------------------------

	profile, err := GetProfileService(
		r.Context(),
		authUserID,
	)

	if err != nil {

		fmt.Print(err)

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: "failed to fetch profile",
			},
		)

		return
	}

	// ----------------------------------------------------------
	// Response
	// ----------------------------------------------------------

	utils.WriteJSON(
		w,
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "profile fetched successfully",
			Data:    profile,
		},
	)
}

// ----------------------------------------------------------------------------------------------------------

func CompleteProfileHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Get Authenticated User ID
	// ----------------------------------------------------------

	authUserID, ok := r.Context().
		Value(middleware.UserIDKey).
		(string)

	if !ok {

		utils.WriteJSON(
			w,
			http.StatusUnauthorized,
			shared.APIResponse{
				Success: false,
				Message: "unauthorized",
			},
		)

		return
	}

	// ----------------------------------------------------------
	// Complete Profile
	// ----------------------------------------------------------

	profile, err := CompleteProfileService(
		r.Context(),
		authUserID,
		r,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: "failed to complete profile",
			},
		)

		return
	}

	// ----------------------------------------------------------
	// Response
	// ----------------------------------------------------------

	utils.WriteJSON(
		w,
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "profile completed successfully",
			Data:    profile,
		},
	)
}