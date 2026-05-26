package state

import (
	"net/http"
	"strings"

	"cobackend/internal/shared"
	"cobackend/internal/utils"
)

// GetStatesHandler returns all available states.
//
// Authorization:
//
//	- Only authenticated users can fetch states.
//
// Query Params:
//
//	- search:
//	  Filters states by name.
//
//	- order:
//	  Sort order for results.
//
// Responses:
//
//	- 200:
//	  States fetched successfully.
//
//	- 401:
//	  Unauthorized.
//
//	- 500:
//	  Internal server error.
func GetStatesHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Parse Query Parameters
	// ----------------------------------------------------------

	search := strings.TrimSpace(
		r.URL.Query().Get("search"),
	)

	order := strings.TrimSpace(
		r.URL.Query().Get("order"),
	)

	queryParams := GetStatesQueryParams{
		Search: search,
		Order:  order,
	}

	// ----------------------------------------------------------
	// Fetch States
	// ----------------------------------------------------------

	states, err := GetStatesService(
		r.Context(),
		queryParams,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: "failed to fetch states",
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
			Message: "states fetched successfully",
			Data:    states,
		},
	)
}