package academy

import (
	"net/http"
	"encoding/json"

	"cobackend/internal/utils"
	"cobackend/internal/shared"
	"cobackend/internal/middleware"

	"github.com/go-chi/chi/v5"

	"strings"
	"strconv"
	"fmt"
)

func CreateAcademyHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	defer r.Body.Close()

	var input CreateAcademyInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&input)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid request body",
			},
		)

		return
	}

	authUserID := r.Context().Value(middleware.UserIDKey).(string)

	academy, err := CreateAcademyService(
		r.Context(),
		authUserID,
		input,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: err.Error(),
			},
		)

		return
	}

	utils.WriteJSON(
		w,
		http.StatusCreated,
		shared.APIResponse{
			Success: true,
			Message: "academy created successfully",
			Data:    academy,
		},
	)
}

// GetAcademiesHandler fetches a list of academies
// based on query parameters.
//
// Authorization:
//
//	- district_admin
//
// Query Params:
//
//	- page (optional, default: 1): Page number to fetch (must be >= 1).
//	- limit (optional, default: 10): Items per page (must be >= 1, or string "all").
//	- search (optional): Filters results by name or address using case-insensitive substring match.
//	- state_id (optional): Filters results by state ID.
//	- district_id (optional): Filters results by district ID.
//	- sort_by (optional, default: "name"): Fields allowed: "id", "name", "state_id", "district_id", "created_at", "updated_at".
//	- order_by (optional, default: "asc"): Allowed values: "asc" or "desc".
//
// Responses:
//
//	- 200:
//	  Academies fetched successfully. Returns a paginated object containing:
//	    - items: array of Academy objects
//	    - page, limit, total, total_pages, has_next, has_previous metadata fields.
//
//	- 400:
//	  Invalid query parameters. Possible error messages:
//	    - "invalid state_id" (state_id is not a valid integer)
//	    - "invalid district_id" (district_id is not a valid integer)
//	    - "invalid page" (page value is less than 1 or not a valid integer)
//	    - "invalid limit" (limit value is less than 1 or not a valid integer/all)
//	    - "invalid sort_by field" (sort_by column is not in allowed sort list)
//	    - "invalid order_by value" (order_by direction is not asc/desc)
//
//	- 401:
//	  Unauthorized (missing or invalid JWT token).
//
//	- 403:
//	  Forbidden (user does not possess the district_admin role).
//
//	- 500:
//	  Internal server error (unexpected query execution failure).
func GetAcademiesHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Parse Query Parameters
	// ----------------------------------------------------------

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := strings.TrimSpace(
		r.URL.Query().Get("search"),
	)
	stateStr := r.URL.Query().Get("state_id")
	districtStr := r.URL.Query().Get("district_id")
	sortBy := r.URL.Query().Get("sort_by")
	orderBy := r.URL.Query().Get("order_by")

	stateID := 0
	if stateStr != "" {

		parsed, err := strconv.Atoi(stateStr)
		if err != nil {

			utils.WriteJSON(
				w,
				http.StatusBadRequest,
				shared.APIResponse{
					Success: false,
					Message: "invalid state_id",
				},
			)

			return
		}

		stateID = parsed
	}

	districtID := 0
	if districtStr != "" {

		parsed, err := strconv.Atoi(districtStr)
		if err != nil {

			utils.WriteJSON(
				w,
				http.StatusBadRequest,
				shared.APIResponse{
					Success: false,
					Message: "invalid district_id",
				},
			)

			return
		}

		districtID = parsed
	}

	// ----------------------------------------------------------
	// Validate Pagination Limits
	// ----------------------------------------------------------

	page := 1
	if pageStr != "" {

		parsed, err := strconv.Atoi(pageStr)
		if err != nil || parsed < 1 {

			utils.WriteJSON(
				w,
				http.StatusBadRequest,
				shared.APIResponse{
					Success: false,
					Message: "invalid page",
				},
			)

			return
		}

		page = parsed
	}

	limit := 10
	if limitStr != "" {

		if limitStr == "all" {

			limit = 0

		} else {

			parsed, err := strconv.Atoi(limitStr)
			if err != nil || parsed < 1 {

				utils.WriteJSON(
					w,
					http.StatusBadRequest,
					shared.APIResponse{
						Success: false,
						Message: "invalid limit",
					},
				)

				return
			}

			limit = parsed
		}
	}

	// ----------------------------------------------------------
	// Validate Sorting Fields
	// ----------------------------------------------------------

	if sortBy == "" {
		sortBy = "name"
	}

	if orderBy == "" {
		orderBy = "asc"
	}

	_, exists := AllowedAcademySortFields[sortBy]

	if !exists {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid sort_by field",
			},
		)

		return
	}

	orderBy = strings.ToUpper(orderBy)

	if orderBy != "ASC" && orderBy != "DESC" {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid order_by value",
			},
		)

		return
	}

	query := GetAcademiesQuery{
		Page:       page,
		Limit:      limit,
		Search:     search,
		StateID:    stateID,
		DistrictID: districtID,
		SortBy:     sortBy,
		OrderBy:    orderBy,
	}

	// ----------------------------------------------------------
	// Execute Service
	// ----------------------------------------------------------

	result, err := GetAcademiesService(
		r.Context(),
		query,
	)

	if err != nil {
		fmt.Print(err)

		// utils.WriteError(
		// 	w,
		// 	err,
		// 	"failed to fetch academies",
		// )

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
			Message: "academies fetched successfully",
			Data:    result,
		},
	)
}

func CreateAcademyBuildingHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	defer r.Body.Close()

	var input CreateAcademyBuildingInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&input)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid request body",
			},
		)

		return
	}

	authUserID := r.Context().Value(middleware.UserIDKey).(string)

	academyBuilding, err := CreateAcademyBuildingService(
		r.Context(),
		authUserID,
		input,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: err.Error(),
			},
		)

		return
	}

	utils.WriteJSON(
		w,
		http.StatusCreated,
		shared.APIResponse{
			Success: true,
			Message: "academy building created successfully",
			Data:    academyBuilding,
		},
	)
}

func AddAcademyBuildingDisciplineHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	defer r.Body.Close()

	var input AddAcademyBuildingDisciplineInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&input)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid request body",
			},
		)

		return
	}

	buildingIDParam := chi.URLParam(
		r,
		"buildingID",
	)

	buildingID, err := strconv.ParseInt(
		buildingIDParam,
		10,
		64,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid building id",
			},
		)

		return
	}

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		( string )

	response, err := AddAcademyBuildingDisciplineService(
		r.Context(),
		authUserID,
		buildingID,
		input,
	)

	if err != nil {

		fmt.Print(err)

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "Internal Server Error",
			},
		)
		return
	}

	utils.WriteJSON(
		w,
		http.StatusCreated,
		shared.APIResponse{
			Success: true,
			Message: "building discipline added successfully",
			Data:    response,
		},
	)
}

// ============================================================================
// handler.go
// ============================================================================

func AddAcademyBuildingEventHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	defer r.Body.Close()

	var input AddAcademyBuildingEventInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&input)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid request body",
			},
		)

		return
	}

	buildingIDParam := chi.URLParam(
		r,
		"buildingID",
	)

	buildingID, err := strconv.ParseInt(
		buildingIDParam,
		10,
		64,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid building id",
			},
		)

		return
	}

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		( string )

	response, err := AddAcademyBuildingEventService(
		r.Context(),
		authUserID,
		buildingID,
		input,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: err.Error(),
			},
		)

		return
	}

	utils.WriteJSON(
		w,
		http.StatusCreated,
		shared.APIResponse{
			Success: true,
			Message: "building event added successfully",
			Data:    response,
		},
	)
}

func AddAcademyBuildingLaneHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	defer r.Body.Close()

	var input AddAcademyBuildingLaneInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&input)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid request body",
			},
		)

		return
	}

	buildingIDParam := chi.URLParam(
		r,
		"buildingID",
	)

	buildingID, err := strconv.ParseInt(
		buildingIDParam,
		10,
		64,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid building id",
			},
		)

		return
	}

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		(string)

	response, err := AddAcademyBuildingLaneService(
		r.Context(),
		authUserID,
		buildingID,
		input,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: err.Error(),
			},
		)

		return
	}

	utils.WriteJSON(
		w,
		http.StatusCreated,
		shared.APIResponse{
			Success: true,
			Message: "building lane added successfully",
			Data:    response,
		},
	)
}

// ============================================================================
// handler.go
// ============================================================================

func GetAcademyBuildingsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		( string )

	buildings, err := GetAcademyBuildingsService(
		r.Context(),
		authUserID,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: err.Error(),
			},
		)

		return
	}

	utils.WriteJSON(
		w,
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "academy buildings fetched successfully",
			Data:    buildings,
		},
	)
}

func GetAvailableLanesHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	buildingIDParam := chi.URLParam(
		r,
		"buildingID",
	)

	buildingID, err := strconv.ParseInt(
		buildingIDParam,
		10,
		64,
	)

	if err != nil || buildingID <= 0 {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid building id",
			},
		)

		return
	}

	response, err := GetAvailableLanesService(
		r.Context(),
		buildingID,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: err.Error(),
			},
		)

		return
	}

	utils.WriteJSON(
		w,
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "available lanes fetched successfully",
			Data:    response,
		},
	)
}

