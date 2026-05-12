package stateadmin

import (
	"encoding/json"
	"errors"
	"net/http"

	// "strings"

	"github.com/go-chi/chi/v5"

	"cobackend/internal/shared"
	"cobackend/internal/utils"

	"cobackend/internal/middleware"
	"fmt"

	"cobackend/internal/validation"

)

func InviteStateAdminHandler(w http.ResponseWriter, r *http.Request) {
	var input InviteStateAdminInput

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

	inviteLink, err := InviteStateAdminService(r.Context(), input, authUserID)
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
				Message: "failed to create state admin invite",
			},
		)

		return
	}

	utils.WriteJSON(w, http.StatusCreated, shared.APIResponse{
		Success: true,
		Message: "state admin invitation created successfully",
		Data: inviteLink,
	})
}

// func GetStateAdminsHandler(w http.ResponseWriter, r *http.Request) {

// 	pageStr := r.URL.Query().Get("page")
// 	limitStr := r.URL.Query().Get("limit")
// 	search := strings.TrimSpace(
// 		r.URL.Query().Get("search"),
// 	)
// 	// stateStr := r.URL.Query().Get("assigned_state")

// 	// default values
// 	page := 1
// 	limit := 10
// 	state := 0

// 	// parse page
// 	if pageStr != "" {

// 		parsedPage, err := strconv.Atoi(pageStr)
// 		if err != nil || parsedPage < 1 {
// 			utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 				Success: false,
// 				Message: "invalid page parameter",
// 			})
// 			return
// 		}

// 		page = parsedPage
// 	}

// 	// parse limit
// 	if limitStr != "" {

// 		parsedLimit, err := strconv.Atoi(limitStr)
// 		if err != nil || parsedLimit < 1 {
// 			utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 				Success: false,
// 				Message: "invalid limit parameter",
// 			})
// 			return
// 		}

// 		limit = parsedLimit
// 	}

// 	// parse assigned_state
// 	// if stateStr != "" {

// 	// 	parsedState, err := strconv.Atoi(stateStr)
// 	// 	if err != nil || parsedState < 1 {
// 	// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 	// 			Success: false,
// 	// 			Message: "invalid assigned_state parameter",
// 	// 		})
// 	// 		return
// 	// 	}

// 	// 	state = parsedState
// 	// }

// 	query := GetStateAdminsQuery{
// 		Page:   page,
// 		Limit:  limit,
// 		Search: search,
// 		AssignedState:  state,
// 	}

// 	stateAdmins, err := GetStateAdminsService(
// 		r.Context(),
// 		query,
// 	)

// 	if err != nil {
// 		utils.WriteJSON(w, http.StatusInternalServerError, shared.APIResponse{
// 			Success: false,
// 			Message: "failed to fetch state admins",
// 		})
// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusOK, shared.APIResponse{
// 		Success: true,
// 		Message: "state admins fetched successfully",
// 		Data:    stateAdmins,
// 	})
// }

func UpdateAssignedStateHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input UpdateStateInput

	id := chi.URLParam(r, "profile_id")

	if id == "" {
		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: "state admin id is required",
		})
		return
	}

	if !validation.IsValidUUID(id) {
		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: "invalid state admin id",
		})
		return
	}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&input)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: "invalid request body",
		})
		return
	}

	if input.State == 0 {
		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: "state id is required",
		})
		return
	}

	err = UpdateStateService(
		r.Context(),
		id,
		input,
	)

	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, shared.APIResponse{
		Success: true,
		Message: "assigned state updated successfully",
	})
}

func DeleteStateAdminHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	profileID := chi.URLParam(r, "profile_id")

	if profileID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: "state admin id is required",
		})
		return
	}

	if !validation.IsValidUUID(profileID) {
		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: shared.ErrInvalidUUID.Error(),
		})
		return
	}

	err := DeleteStateAdminService(
		r.Context(),
		profileID,
	)

	if err != nil {

		statusCode := http.StatusBadRequest

		if errors.Is(err, shared.ErrStateAdminNotFound) {
			statusCode = http.StatusNotFound
		}

		utils.WriteJSON(w, statusCode, shared.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, shared.APIResponse{
		Success: true,
		Message: "state admin deleted successfully",
	})
}