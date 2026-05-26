package district

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"cobackend/internal/shared"
	"cobackend/internal/utils"
)

// GetDistrictsHandler returns all districts.
//
// Authorization:
//
//	- Only authenticated users can fetch districts.
//
// Responses:
//
//	- 200:
//	  Districts fetched successfully.
//
//	- 401:
//	  Unauthorized.
//
//	- 500:
//	  Internal server error.
func GetDistrictsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Fetch Districts
	// ----------------------------------------------------------

	districts, err := GetDistrictsService(
		r.Context(),
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: "failed to fetch districts",
			},
		)

		return
	}

	// ----------------------------------------------------------
	// Success Response
	// ----------------------------------------------------------

	utils.WriteJSON(
		w,
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "districts fetched successfully",
			Data:    districts,
		},
	)
}

// GetDistrictsByStateIdHandler returns all districts
// belonging to a particular state.
//
// Authorization:
//
//	- Only authenticated users can fetch districts.
//
// Path Params:
//
//	- state_id:
//	  Unique identifier of the state.
//
// Query Params:
//
//	- search:
//	  Filters districts by name.
//
//	- order:
//	  Sort order for results.
//
// Responses:
//
//	- 200:
//	  Districts fetched successfully.
//
//	- 400:
//	  Invalid request parameters.
//
//	- 401:
//	  Unauthorized.
//
//	- 500:
//	  Internal server error.
func GetDistrictsByStateIdHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Parse Path Parameters
	// ----------------------------------------------------------

	stateIDStr := chi.URLParam(
		r,
		"state_id",
	)

	if stateIDStr == "" {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "state id is required",
			},
		)

		return
	}

	stateID, err := strconv.Atoi(
		stateIDStr,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid state id",
			},
		)

		return
	}

	// ----------------------------------------------------------
	// Parse Query Parameters
	// ----------------------------------------------------------

	search := strings.TrimSpace(
		r.URL.Query().Get("search"),
	)

	order := strings.TrimSpace(
		r.URL.Query().Get("order"),
	)

	queryParams := GetDistrictQueryParams{
		Search: search,
		Order:  order,
	}

	// ----------------------------------------------------------
	// Fetch Districts
	// ----------------------------------------------------------

	districts, err := GetDistrictsByStateIdService(
		r.Context(),
		stateID,
		queryParams,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: "failed to fetch districts",
			},
		)

		return
	}

	// ----------------------------------------------------------
	// Success Response
	// ----------------------------------------------------------

	utils.WriteJSON(
		w,
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "districts fetched successfully",
			Data:    districts,
		},
	)
}