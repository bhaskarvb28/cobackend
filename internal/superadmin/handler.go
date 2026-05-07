package superadmin

import (
	"encoding/json"
	"net/http"

	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func CreateStateAdminHandler(w http.ResponseWriter, r *http.Request) {
	var input CreateStateAdminInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	err = CreateStateAdminService(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "State admin created successfully",
	})

}

func UpdateAssignedStateHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := chi.URLParam(r, "id")

	var input UpdateAssignedStateInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = UpdateAssignedStateService(
		r.Context(),
		id,
		input,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Assigned state updated successfully",
	})
}

func GetStateAdminsHandler(w http.ResponseWriter, r *http.Request) {

	pageStr :=  r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	search := r.URL.Query().Get("search")
	stateStr := r.URL.Query().Get("state")

	state := 0

	if stateStr != "" {

		parsedState, err := strconv.Atoi(stateStr)
		if err == nil {
			state = parsedState
		}
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	query := GetStateAdminsQuery{
		Page:  page,
		Limit: limit,
		Search: search,
		State: state,
	}

	stateAdmins, err := GetStateAdminsService(r.Context(), query)
	if err != nil {
		http.Error(w, "Failed to fetch state admins", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stateAdmins)
}

func DeleteStateAdminHandler(w http.ResponseWriter, r *http.Request) {
	
	id := chi.URLParam(r, "id")

	_, err := uuid.Parse(id)
	if err != nil {
		http.Error(w, "Invalid state admin id", http.StatusBadRequest)
		return
	}

	err = DeleteStateAdminService(
		r.Context(),
		id,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "State admin deleted successfully",
	})

}