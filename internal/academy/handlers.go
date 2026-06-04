package academy

import (
	"encoding/json"
	"net/http"

	"cobackend/internal/middleware"
	"cobackend/internal/player"
	"cobackend/internal/shared"
	"cobackend/internal/utils"
	"cobackend/internal/academyCoach"

	"github.com/go-chi/chi/v5"

	"fmt"
	"strconv"
	"strings"
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

	err = CreateAcademyService(
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

func GetAcademyPlayersHandler(
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

	disciplineStr := r.URL.Query().Get(
		"discipline_id",
	)

	status := strings.TrimSpace(
		r.URL.Query().Get("status"),
	)

	coachAssignedStr := strings.TrimSpace(
		r.URL.Query().Get("coach_assigned"),
	)

	sortBy := r.URL.Query().Get("sort_by")
	orderBy := r.URL.Query().Get("order_by")

	// ----------------------------------------------------------
	// Parse Discipline
	// ----------------------------------------------------------

	disciplineID := 0

	if disciplineStr != "" {

		parsed, err := strconv.Atoi(
			disciplineStr,
		)

		if err != nil {

			utils.WriteJSON(
				w,
				http.StatusBadRequest,
				shared.APIResponse{
					Success: false,
					Message: "invalid discipline_id",
				},
			)

			return
		}

		disciplineID = parsed
	}

	// ----------------------------------------------------------
	// Parse Coach Assigned
	// ----------------------------------------------------------

	var coachAssigned *bool

	if coachAssignedStr != "" {

		parsed, err := strconv.ParseBool(
			coachAssignedStr,
		)

		if err != nil {

			utils.WriteJSON(
				w,
				http.StatusBadRequest,
				shared.APIResponse{
					Success: false,
					Message: "invalid coach_assigned",
				},
			)

			return
		}

		coachAssigned = &parsed
	}

	// ----------------------------------------------------------
	// Pagination
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
	// Sorting
	// ----------------------------------------------------------

	if sortBy == "" {
		sortBy = "created_at"
	}

	if orderBy == "" {
		orderBy = "desc"
	}

	_, exists := AllowedAcademyPlayerSortFields[
		sortBy,
	]

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

	// ----------------------------------------------------------
	// Build Query
	// ----------------------------------------------------------

	query := player.GetAcademyPlayersQuery{
		Page:          page,
		Limit:         limit,
		Search:        search,
		DisciplineID:  disciplineID,
		CoachAssigned: coachAssigned,
		Status:        status,
		SortBy:        sortBy,
		OrderBy:       orderBy,
	}

	// ----------------------------------------------------------
	// Get Auth User
	// ----------------------------------------------------------

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		( string )

	// ----------------------------------------------------------
	// Execute Service
	// ----------------------------------------------------------

	result, err := GetAcademyPlayersService(
		r.Context(),
		authUserID,
		query,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: err.Error(),
			},
		)

		fmt.Print(err)

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
			Message: "players fetched successfully",
			Data:    result,
		},
	)
}

func GetAcademyPlayerHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	playerID := chi.URLParam(
		r,
		"playerID",
	)

	if playerID == "" {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "playerID is required",
			},
		)

		return
	}

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		( string )

	result, err := GetAcademyPlayerService(
		r.Context(),
		authUserID,
		playerID,
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
			Message: "player fetched successfully",
			Data:    result,
		},
	)
}

func GetAcademyCoachesHandler(
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

	disciplineStr := r.URL.Query().Get(
		"discipline_id",
	)

	sortBy := r.URL.Query().Get("sort_by")
	orderBy := r.URL.Query().Get("order_by")

	// ----------------------------------------------------------
	// Parse Discipline
	// ----------------------------------------------------------

	disciplineID := 0

	if disciplineStr != "" {

		parsed, err := strconv.Atoi(
			disciplineStr,
		)

		if err != nil {

			utils.WriteJSON(
				w,
				http.StatusBadRequest,
				shared.APIResponse{
					Success: false,
					Message: "invalid discipline_id",
				},
			)

			return
		}

		disciplineID = parsed
	}

	// ----------------------------------------------------------
	// Pagination
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

			parsed, err := strconv.Atoi(
				limitStr,
			)

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
	// Sorting
	// ----------------------------------------------------------

	if sortBy == "" {
		sortBy = "created_at"
	}

	if orderBy == "" {
		orderBy = "desc"
	}

	orderBy = strings.ToUpper(orderBy)

	if orderBy != "ASC" &&
		orderBy != "DESC" {

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

	// ----------------------------------------------------------
	// Build Query
	// ----------------------------------------------------------

	query := academyCoach.GetAcademyCoachesQuery{
		Page:         page,
		Limit:        limit,
		Search:       search,
		DisciplineID: disciplineID,
		SortBy:       sortBy,
		OrderBy:      orderBy,
	}

	// ----------------------------------------------------------
	// Get Auth User
	// ----------------------------------------------------------

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		( string )

	// ----------------------------------------------------------
	// Execute Service
	// ----------------------------------------------------------

	result, err := GetAcademyCoachesService(
		r.Context(),
		authUserID,
		query,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: err.Error(),
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
			Message: "coaches fetched successfully",
			Data:    result,
		},
	)
}

func GetAcademyCoachHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Params
	// ----------------------------------------------------------

	coachID := chi.URLParam(
		r,
		"coachID",
	)

	if coachID == "" {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "coachID is required",
			},
		)

		return
	}

	// ----------------------------------------------------------
	// Get Auth User
	// ----------------------------------------------------------

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		( string )

	// ----------------------------------------------------------
	// Execute Service
	// ----------------------------------------------------------

	result, err := GetAcademyCoachService(
		r.Context(),
		authUserID,
		coachID,
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

	// ----------------------------------------------------------
	// Success Response
	// ----------------------------------------------------------

	utils.WriteJSON(
		w,
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "coach fetched successfully",
			Data:    result,
		},
	)
}

func AssignCoachHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Params
	// ----------------------------------------------------------

	playerID := chi.URLParam(
		r,
		"playerID",
	)

	if playerID == "" {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "playerID is required",
			},
		)

		return
	}

	// ----------------------------------------------------------
	// Parse Body
	// ----------------------------------------------------------

	var payload academyCoach.AssignCoachInput

	err := json.NewDecoder(
		r.Body,
	).Decode(&payload)

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

	if payload.CoachUserID == "" {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "coach_user_id is required",
			},
		)

		return
	}

	// ----------------------------------------------------------
	// Get Auth User
	// ----------------------------------------------------------

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		( string )

	// ----------------------------------------------------------
	// Execute Service
	// ----------------------------------------------------------

	err = AssignCoachService(
		r.Context(),
		authUserID,
		playerID,
		payload.CoachUserID,
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

	// ----------------------------------------------------------
	// Success Response
	// ----------------------------------------------------------

	utils.WriteJSON(
		w,
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "coach assigned successfully",
		},
	)
}

func RemoveCoachHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Params
	// ----------------------------------------------------------

	playerID := chi.URLParam(
		r,
		"playerID",
	)

	if playerID == "" {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "playerID is required",
			},
		)

		return
	}

	// ----------------------------------------------------------
	// Get Auth User
	// ----------------------------------------------------------

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		( string )

	// ----------------------------------------------------------
	// Execute Service
	// ----------------------------------------------------------

	err := RemoveCoachService(
		r.Context(),
		authUserID,
		playerID,
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

	// ----------------------------------------------------------
	// Success Response
	// ----------------------------------------------------------

	utils.WriteJSON(
		w,
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "coach removed successfully",
		},
	)
}

func GetAcademyBuildingHandler(
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

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		( string )

	response, err := GetAcademyBuildingService(
		r.Context(),
		authUserID,
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
			Message: "building fetched successfully",
			Data:    response,
		},
	)
}

func DeleteAcademyBuildingHandler(
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

	authUserID := r.Context().
		Value(
			middleware.UserIDKey,
		).
		(string)

	err = DeleteAcademyBuildingService(
		r.Context(),
		authUserID,
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
			Message: "building deleted successfully",
		},
	)
}

func UpdateAcademyBuildingHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	defer r.Body.Close()

	var input UpdateAcademyBuildingInput

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

	buildingID, err := strconv.ParseInt(
		chi.URLParam(r, "buildingID"),
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

	err = UpdateAcademyBuildingService(
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
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "building updated successfully",
		},
	)
}

func RemoveAcademyBuildingDisciplineHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	buildingID, err := strconv.ParseInt(
		chi.URLParam(r, "buildingID"),
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

	disciplineID, err := strconv.Atoi(
		chi.URLParam(r, "disciplineID"),
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid discipline id",
			},
		)

		return
	}

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		( string )

	err = RemoveAcademyBuildingDisciplineService(
		r.Context(),
		authUserID,
		buildingID,
		disciplineID,
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
			Message: "building discipline removed successfully",
		},
	)
}

func RemoveAcademyBuildingEventHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	buildingID, err := strconv.ParseInt(
		chi.URLParam(r, "buildingID"),
		10,
		64,
	)

	if err != nil {

		return
	}

	eventID, err := strconv.Atoi(
		chi.URLParam(r, "eventID"),
	)

	if err != nil {

		return
	}

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		( string )

	err = RemoveAcademyBuildingEventService(
		r.Context(),
		authUserID,
		buildingID,
		eventID,
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
			Message: "building event removed successfully",
		},
	)
}

func GetAcademyBuildingLanesHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	buildingID, err := strconv.ParseInt(
		chi.URLParam(r, "buildingID"),
		10,
		64,
	)

	if err != nil {

		return
	}

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		( string )

	response, err := GetAcademyBuildingLanesService(
		r.Context(),
		authUserID,
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
			Message: "building lanes fetched successfully",
			Data:    response,
		},
	)
}

func UpdateAcademyBuildingLaneHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	defer r.Body.Close()

	var input UpdateAcademyBuildingLaneInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&input)

	if err != nil {

		return
	}

	laneID, err := strconv.ParseInt(
		chi.URLParam(r, "laneID"),
		10,
		64,
	)

	if err != nil {

		return
	}

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		( string )

	err = UpdateAcademyBuildingLaneService(
		r.Context(),
		authUserID,
		laneID,
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
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "lane updated successfully",
		},
	)
}

func DeleteAcademyBuildingLaneHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	laneID, err := strconv.ParseInt(
		chi.URLParam(r, "laneID"),
		10,
		64,
	)

	if err != nil {

		return
	}

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		( string )

	err = DeleteAcademyBuildingLaneService(
		r.Context(),
		authUserID,
		laneID,
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
			Message: "lane deleted successfully",
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


// GetDistrictAdminAcademiesHandler fetches
// academies belonging to the authenticated
// district admin.
func GetDistrictAdminAcademiesHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Get User ID From Context
	// ----------------------------------------------------------

	userID, ok := r.Context().Value(
		middleware.UserIDKey,
	).(string)

	if !ok || userID == "" {

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
	// Parse Query Params
	// ----------------------------------------------------------

	pageStr := r.URL.Query().Get("page")

	limitStr := r.URL.Query().Get("limit")

	search := strings.TrimSpace(
		r.URL.Query().Get("search"),
	)

	sortBy := r.URL.Query().Get("sort_by")

	orderBy := r.URL.Query().Get("order_by")

	// ----------------------------------------------------------
	// Validate Page
	// ----------------------------------------------------------

	page := 1

	if pageStr != "" {

		parsed, err :=
			strconv.Atoi(pageStr)

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

	// ----------------------------------------------------------
	// Validate Limit
	// ----------------------------------------------------------

	limit := 10

	if limitStr != "" {

		if limitStr == "all" {

			limit = 0

		} else {

			parsed, err :=
				strconv.Atoi(limitStr)

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
	// Validate Sorting
	// ----------------------------------------------------------

	if sortBy == "" {
		sortBy = "name"
	}

	if orderBy == "" {
		orderBy = "asc"
	}

	_, exists :=
		AllowedAcademySortFields[
			sortBy,
		]

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

	orderBy =
		strings.ToUpper(orderBy)

	if orderBy != "ASC" &&
		orderBy != "DESC" {

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

	// ----------------------------------------------------------
	// Build Query
	// ----------------------------------------------------------

	query := GetAcademiesQuery{
		Page: page,

		Limit: limit,

		Search: search,

		SortBy: sortBy,

		OrderBy: orderBy,
	}

	// ----------------------------------------------------------
	// Execute Service
	// ----------------------------------------------------------

	result, err :=
		GetDistrictAdminAcademiesService(
			r.Context(),
			userID,
			query,
		)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: "failed to fetch district academies",
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
			Message: "district academies fetched successfully",
			Data: result,
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

