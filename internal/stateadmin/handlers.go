package stateadmin

import (
	"encoding/json"
	"net/http"

	// "strconv"

	// "github.com/go-chi/chi/v5"

	"cobackend/internal/shared"
	"cobackend/internal/utils"

	"cobackend/internal/middleware"


	// "errors"

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
		utils.WriteJSON(w, http.StatusInternalServerError, shared.APIResponse{
			Success: false,
			Message: "failed to create state admin invite",
		})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, shared.APIResponse{
		Success: true,
		Message: "state admin invitation created successfully",
		Data: inviteLink,
	})
}

// func CreateStateAdminHandler(w http.ResponseWriter, r *http.Request) {
// 	var input CreateStateAdminInput

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

// 	// trim whitespace
// 	input.FirstName = strings.TrimSpace(input.FirstName)
// 	input.LastName = strings.TrimSpace(input.LastName)
// 	input.Email = strings.TrimSpace(input.Email)
// 	input.ContactNumber = strings.TrimSpace(input.ContactNumber)

// 	// basic request validation
// 	if input.FirstName == "" ||
// 		input.LastName == "" ||
// 		input.Email == "" ||
// 		input.Password == "" ||
// 		input.ContactNumber == "" ||
// 		input.AssignedState == "" {

// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "all fields are required",
// 		})
// 		return
// 	}

// 	// validate email format
// 	if !validation.IsValidEmail(input.Email) {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: shared.ErrInvalidEmailFormat.Error(),
// 		})
// 		return
// 	}

// 	// validate phone number
// 	if !validation.IsValidIndianPhone(input.ContactNumber) {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: shared.ErrInvalidPhoneNumber.Error(),
// 		})
// 		return
// 	}

// 	// validate password strength
// 	if !validation.IsStrongPassword(input.Password) {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: shared.ErrWeakPassword.Error(),
// 		})
// 		return
// 	}

// 	err = CreateStateAdminService(r.Context(), input)
// 	if err != nil {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: err.Error(),
// 		})
// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusCreated, shared.APIResponse{
// 		Success: true,
// 		Message: "state admin created successfully",
// 	})
// }


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

// func UpdateAssignedStateHandler(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) {

// 	id := chi.URLParam(r, "id")

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

// 	var input UpdateAssignedStateInput

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

// 	// validate assigned state
// 	if input.AssignedState == 0 {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "assigned state is required",
// 		})
// 		return
// 	}

// 	err = UpdateAssignedStateService(
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

// 	id := chi.URLParam(r, "id")

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
// 			Message: shared.ErrInvalidUUID.Error(),
// 		})
// 		return
// 	}

// 	err := DeleteStateAdminService(
// 		r.Context(),
// 		id,
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