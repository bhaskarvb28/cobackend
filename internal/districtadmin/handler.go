package districtadmin

import (
	"encoding/json"
	"net/http"

	"cobackend/internal/middleware"
	"cobackend/internal/utils"
	"cobackend/internal/shared"
	"errors"

	"fmt"

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

// func CreateDistrictAdminHandler(w http.ResponseWriter, r *http.Request) {
// 	var input CreateDistrictAdminInput

// 	err := json.NewDecoder(r.Body).Decode(&input)
// 	if err != nil {
// 		http.Error(w, "Invalid request body", http.StatusBadRequest)
// 		return
// 	}

// 	err = CreateDistrictAdminService(r.Context(), input)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message": "District admin created successfully",
// 	})
// }

// func GetDistrictAdminsHandler(w http.ResponseWriter, r *http.Request) {
// 	pageStr := r.URL.Query().Get("page")
// 	limitStr := r.URL.Query().Get("limit")
// 	search := r.URL.Query().Get("search")
// 	stateStr := r.URL.Query().Get("state_id")
// 	districtStr := r.URL.Query().Get("district_id")
// 	status := r.URL.Query().Get("status")

// 	stateID := 0
// 	if stateStr != "" {
// 		if parsed, err := strconv.Atoi(stateStr); err == nil {
// 			stateID = parsed
// 		}
// 	}

// 	districtID := 0
// 	if districtStr != "" {
// 		if parsed, err := strconv.Atoi(districtStr); err == nil {
// 			districtID = parsed
// 		}
// 	}

// 	page, err := strconv.Atoi(pageStr)
// 	if err != nil || page < 1 {
// 		page = 1
// 	}

// 	limit, err := strconv.Atoi(limitStr)
// 	if err != nil || limit < 1 {
// 		limit = 10
// 	}

// 	query := GetDistrictAdminsQuery{
// 		Page:       page,
// 		Limit:      limit,
// 		Search:     search,
// 		StateID:    stateID,
// 		DistrictID: districtID,
// 		Status:     status,
// 	}

// 	admins, err := GetDistrictAdminsService(r.Context(), query)
// 	if err != nil {
// 		http.Error(w, "Failed to fetch district admins", http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(admins)
// }

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
