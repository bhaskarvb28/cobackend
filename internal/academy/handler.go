package academy

import (
	"net/http"
	"encoding/json"

	"cobackend/internal/utils"
	"cobackend/internal/shared"

	"strings"
	"strconv"
	"fmt"
)

func CreateAcademyHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	var input CreateAcademyInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	defer r.Body.Close()

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

	if strings.TrimSpace(input.Name) == "" {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "academy name is required",
			},
		)

		return
	}

	if input.DistrictID <= 0 {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "valid district_id is required",
			},
		)

		return
	}

	if strings.TrimSpace(input.Address) == "" {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "address is required",
			},
		)

		return
	}

	err = CreateAcademyService(
		r.Context(),
		input,
	)

	if err != nil {

		utils.WriteError(
			w,
			err,
			"failed to create academy",
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

func GetAcademiesHandler(
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

	result, err := GetAcademiesService(
		r.Context(),
		query,
	)

	if err != nil {
		fmt.Print(err)

		utils.WriteError(
			w,
			err,
			"failed to fetch academies",
		)

		return
	}

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