package districtCoach

import (
	"net/http"
	"encoding/json"

	"cobackend/internal/utils"
	"cobackend/internal/shared"

	"cobackend/internal/middleware"

	"errors"

	"fmt"
	"strconv"
	"strings"
)

func InviteDistrictCoachHandler(w http.ResponseWriter, r *http.Request) {
	var input InviteDistrictCoachInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	defer r.Body.Close()

	err := decoder.Decode(&input)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: "invalid request body",
		})
		return
	}

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

	inviteLink, err := InviteDistrictCoachService(r.Context(), input, authUserID)
	
	if err != nil {

		var apiErr *shared.APIError

		if errors.As(err, &apiErr) {

			utils.WriteJSON(
				w,
				apiErr.StatusCode,
				shared.APIResponse{
					Success: false,
					Message: apiErr.Message,
				},
			)

			return
		}

		fmt.Print(err)

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: "failed to create district coach invite",
			},
		)

		return
	}

	utils.WriteJSON(w, http.StatusCreated, shared.APIResponse{
		Success: true,
		Message: "district coach invitation created successfully",
		Data: inviteLink,
	})
}

func GetDistrictCoachesHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

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

	if sortBy == "" {
		sortBy = "first_name"
	}

	if orderBy == "" {
		orderBy = "asc"
	}

	_, exists := AllowedDistrictCoachSortFields[sortBy]

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

	query := GetDistrictCoachesQuery{
		Page:       page,
		Limit:      limit,
		Search:     search,
		StateID:    stateID,
		DistrictID: districtID,
		SortBy:     sortBy,
		OrderBy:    orderBy,
	}

	result, err := GetDistrictCoachesService(
		r.Context(),
		query,
	)

	if err != nil {

		utils.WriteError(
			w,
			err,
			"failed to fetch district coaches",
		)

		return
	}

	utils.WriteJSON(
		w,
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "district coaches fetched successfully",
			Data:    result,
		},
	)
}