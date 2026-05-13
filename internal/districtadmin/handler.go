package districtadmin

import (
	"encoding/json"
	"net/http"

	"cobackend/internal/middleware"
	"cobackend/internal/utils"
	"cobackend/internal/shared"
	"errors"

	"fmt"

	"strconv"
	"strings"

	// "github.com/go-chi/chi/v5"
	// "github.com/google/uuid"
)

func InviteDistrictAdminHandler(w http.ResponseWriter, r *http.Request) {
	var input InviteDistrictAdminInput

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

	inviteLink, err := InviteDistrictAdminService(r.Context(), input, authUserID)
	
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
				Message: "failed to create district admin invite",
			},
		)

		return
	}

	utils.WriteJSON(w, http.StatusCreated, shared.APIResponse{
		Success: true,
		Message: "district admin invitation created successfully",
		Data: inviteLink,
	})
}

func GetDistrictAdminsHandler(
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

	_, exists := AllowedDistrictAdminSortFields[sortBy]

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

	query := GetDistrictAdminsQuery{
		Page:       page,
		Limit:      limit,
		Search:     search,
		StateID:    stateID,
		DistrictID: districtID,
		SortBy: 	sortBy,
		OrderBy: 	orderBy,
	}

	result, err := GetDistrictAdminsService(
		r.Context(),
		query,
	)

	if err != nil {
		fmt.Print(err)

		utils.WriteError(
			w,
			err,
			"failed to fetch district admins",
		)

		return
	}

	utils.WriteJSON(
		w,
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "district admins fetched successfully",
			Data:    result,
		},
	)
}

// func UpdateDistrictAdminHandler(w http.ResponseWriter, r *http.Request) {
// 	id := chi.URLParam(r, "id")

// 	_, err := uuid.Parse(id)
// 	if err != nil {
// 		http.Error(w, "Invalid district admin id", http.StatusBadRequest)
// 		return
// 	}

// 	var input UpdateDistrictAdminInput
// 	err = json.NewDecoder(r.Body).Decode(&input)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	err = UpdateDistrictAdminService(r.Context(), id, input)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message": "District admin updated successfully",
// 	})
// }

// func DeleteDistrictAdminHandler(w http.ResponseWriter, r *http.Request) {
// 	id := chi.URLParam(r, "id")

// 	_, err := uuid.Parse(id)
// 	if err != nil {
// 		http.Error(w, "Invalid district admin id", http.StatusBadRequest)
// 		return
// 	}

// 	err = DeleteDistrictAdminService(r.Context(), id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message": "District admin deleted successfully",
// 	})
// }
