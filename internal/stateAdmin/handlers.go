package stateAdmin

// import (
// 	"encoding/json"
// 	"errors"
// 	"net/http"

// 	"strconv"

// 	"github.com/go-chi/chi/v5"

// 	"cobackend/internal/shared"
// 	"cobackend/internal/utils"

// 	"cobackend/internal/middleware"
// 	"fmt"

// 	"cobackend/internal/validation"

// 	"strings"

// )

// func InviteStateAdminHandler(w http.ResponseWriter, r *http.Request) {
// 	var input InviteStateAdminInput

// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()

// 	defer r.Body.Close()

// 	err := decoder.Decode(&input)
// 	if err != nil {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "invalid request body",
// 		})
// 		return
// 	}

// 	authUserID, ok := r.Context().
// 		Value(middleware.UserIDKey).
// 		(string)

// 	if !ok {
// 		utils.WriteJSON(
// 			w,
// 			http.StatusUnauthorized,
// 			shared.APIResponse{
// 				Success: false,
// 				Message: "unauthorized",
// 			},
// 		)
// 		return
// 	}

// 	inviteLink, err := InviteStateAdminService(r.Context(), input, authUserID)
// 	if err != nil {

// 		var apiErr *shared.APIError

// 		if errors.As(err, &apiErr) {

// 			utils.WriteJSON(
// 				w,
// 				apiErr.StatusCode,
// 				shared.APIResponse{
// 					Success: false,
// 					Message: apiErr.Message,
// 				},
// 			)

// 			return
// 		}

// 		fmt.Print(err)

// 		utils.WriteJSON(
// 			w,
// 			http.StatusInternalServerError,
// 			shared.APIResponse{
// 				Success: false,
// 				Message: "failed to create state admin invite",
// 			},
// 		)

// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusCreated, shared.APIResponse{
// 		Success: true,
// 		Message: "state admin invitation created successfully",
// 		Data: inviteLink,
// 	})
// }

// func GetStateAdminsHandler(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) {

// 	pageStr := r.URL.Query().Get("page")
// 	limitStr := r.URL.Query().Get("limit")
// 	search := strings.TrimSpace(
// 		r.URL.Query().Get("search"),
// 	)
// 	stateStr := r.URL.Query().Get("state_id")
// 	sortBy := r.URL.Query().Get("sort_by")
// 	orderBy := r.URL.Query().Get("order_by")

// 	page := 1
// 	if pageStr != "" {

// 		parsed, err := strconv.Atoi(pageStr)
// 		if err != nil || parsed < 1 {

// 			utils.WriteJSON(
// 				w,
// 				http.StatusBadRequest,
// 				shared.APIResponse{
// 					Success: false,
// 					Message: "invalid page",
// 				},
// 			)

// 			return
// 		}

// 		page = parsed
// 	}

// 	limit := 10
// 	if limitStr != "" {

// 		if limitStr == "all" {

// 			limit = 0

// 		} else {

// 			parsed, err := strconv.Atoi(limitStr)
// 			if err != nil || parsed < 1 {

// 				utils.WriteJSON(
// 					w,
// 					http.StatusBadRequest,
// 					shared.APIResponse{
// 						Success: false,
// 						Message: "invalid limit",
// 					},
// 				)

// 				return
// 			}

// 			limit = parsed
// 		}
// 	}

// 	stateID := 0
// 	if stateStr != "" {

// 		parsed, err := strconv.Atoi(stateStr)
// 		if err != nil {

// 			utils.WriteJSON(
// 				w,
// 				http.StatusBadRequest,
// 				shared.APIResponse{
// 					Success: false,
// 					Message: "invalid state_id",
// 				},
// 			)

// 			return
// 		}

// 		stateID = parsed
// 	}

// 	if sortBy == "" {
// 		sortBy = "first_name"
// 	}

// 	if orderBy == "" {
// 		orderBy = "asc"
// 	}

// 	validSortFields := map[string]bool{
// 		"first_name": true,
// 		"last_name":  true,
// 		"email":      true,
// 		"created_at": true,
// 	}

// 	_, exists := validSortFields[sortBy]

// 	if !exists {

// 		utils.WriteJSON(
// 			w,
// 			http.StatusBadRequest,
// 			shared.APIResponse{
// 				Success: false,
// 				Message: "invalid sort_by field",
// 			},
// 		)

// 		return
// 	}

// 	orderBy = strings.ToUpper(orderBy)

// 	if orderBy != "ASC" && orderBy != "DESC" {

// 		utils.WriteJSON(
// 			w,
// 			http.StatusBadRequest,
// 			shared.APIResponse{
// 				Success: false,
// 				Message: "invalid order_by value",
// 			},
// 		)

// 		return
// 	}

// 	query := GetStateAdminsQuery{
// 		Page:          page,
// 		Limit:         limit,
// 		Search:        search,
// 		StateID: 	   stateID,
// 		SortBy:        sortBy,
// 		OrderBy:       orderBy,
// 	}

// 	result, err := GetStateAdminsService(
// 		r.Context(),
// 		query,
// 	)

// 	if err != nil {
// 		fmt.Print(err)

// 		utils.WriteError(
// 			w,
// 			err,
// 			"failed to fetch state admins",
// 		)

// 		return
// 	}

// 	utils.WriteJSON(
// 		w,
// 		http.StatusOK,
// 		shared.APIResponse{
// 			Success: true,
// 			Message: "state admins fetched successfully",
// 			Data:    result,
// 		},
// 	)
// }

// func UpdateAssignedStateHandler(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) {
// 	var input UpdateStateInput

// 	id := chi.URLParam(r, "profile_id")

// 	if id == "" {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "state admin id is required",
// 		})
// 		return
// 	}

// 	if !validation.IsValidUUID(id) {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "invalid state admin id",
// 		})
// 		return
// 	}

// 	defer r.Body.Close()

// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()

// 	err := decoder.Decode(&input)
// 	if err != nil {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "invalid request body",
// 		})
// 		return
// 	}

// 	if input.StateID == 0 {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "state id is required",
// 		})
// 		return
// 	}

// 	err = UpdateStateService(
// 		r.Context(),
// 		id,
// 		input,
// 	)

// 	if err != nil {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: err.Error(),
// 		})
// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusOK, shared.APIResponse{
// 		Success: true,
// 		Message: "assigned state updated successfully",
// 	})
// }

// func DeleteStateAdminHandler(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) {

// 	profileID := chi.URLParam(r, "profile_id")

// 	if profileID == "" {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "state admin id is required",
// 		})
// 		return
// 	}

// 	if !validation.IsValidUUID(profileID) {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: shared.ErrInvalidUUID.Error(),
// 		})
// 		return
// 	}

// 	err := DeleteStateAdminService(
// 		r.Context(),
// 		profileID,
// 	)

// 	if err != nil {

// 		statusCode := http.StatusBadRequest

// 		if errors.Is(err, shared.ErrStateAdminNotFound) {
// 			statusCode = http.StatusNotFound
// 		}

// 		utils.WriteJSON(w, statusCode, shared.APIResponse{
// 			Success: false,
// 			Message: err.Error(),
// 		})
// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusOK, shared.APIResponse{
// 		Success: true,
// 		Message: "state admin deleted successfully",
// 	})
// }