package invitation

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"fmt"

	"github.com/go-chi/chi/v5"

	"cobackend/internal/shared"
	"cobackend/internal/utils"
	"cobackend/internal/middleware"
	"cobackend/internal/validation"
)

// CreateInvitationHandler creates a new invitation.
//
// Authorization:
//
//	- super_admin
//	- state_admin
//	- district_admin
//	- academy_admin
//
// Request Body:
//
//	{
//		"email": "john@example.com",
//		"role_id": 2,
//		"state_id": 1,
//		"district_id": 10,
//		"academy_id": "uuid"
//	}
//
// Responses:
//
//	- 201:
//	  Invitation created successfully.
//
//	- 400:
//	  Invalid request body or validation failure.
//
//	- 401:
//	  Unauthorized.
//
//	- 403:
//	  Forbidden.
//
//	- 500:
//	  Internal server error.
func CreateInvitationHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	defer r.Body.Close()

	var input CreateInvitationInput

	//------------------------------------------------
	// Decode Request Body
	//------------------------------------------------

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

	//------------------------------------------------
	// Normalize Input
	//------------------------------------------------

	input.Email = strings.ToLower(
		strings.TrimSpace(input.Email),
	)

	input.Role = strings.TrimSpace(
		input.Role,
	)

	input.ScopeType = strings.TrimSpace(
		input.ScopeType,
	)

	input.ScopeID = strings.TrimSpace(
		input.ScopeID,
	)

	//------------------------------------------------
	// Validate Email
	//------------------------------------------------

	if input.Email == "" {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "email is required",
			},
		)

		return
	}

	if !validation.IsValidEmail(
		input.Email,
	) {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid email format",
			},
		)

		return
	}

	//------------------------------------------------
	// Validate Role
	//------------------------------------------------

	if input.Role == "" {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "role is required",
			},
		)

		return
	}

	//------------------------------------------------
	// Authenticated User
	//------------------------------------------------

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

	//------------------------------------------------
	// Create Invitation
	//------------------------------------------------

	response, err := CreateInvitationService(
		r.Context(),
		input,
		authUserID,
	)

	if err != nil {

		fmt.Print(err)

		switch {

		//------------------------------------------------
		// Unauthorized / Forbidden
		//------------------------------------------------

		case errors.Is(
			err,
			shared.ErrUnauthorized,
		):

			utils.WriteJSON(
				w,
				http.StatusUnauthorized,
				shared.APIResponse{
					Success: false,
					Message: err.Error(),
				},
			)

			return

		case errors.Is(
			err,
			shared.ErrForbidden,
		),
			errors.Is(
				err,
				shared.ErrRoleNotAuthorized,
			):

			utils.WriteJSON(
				w,
				http.StatusForbidden,
				shared.APIResponse{
					Success: false,
					Message: err.Error(),
				},
			)

			return

		//------------------------------------------------
		// Validation Errors
		//------------------------------------------------

		case errors.Is(
			err,
			shared.ErrInvalidScope,
		):

			utils.WriteJSON(
				w,
				http.StatusBadRequest,
				shared.APIResponse{
					Success: false,
					Message: err.Error(),
				},
			)

			return

		//------------------------------------------------
		// Conflict Errors
		//------------------------------------------------

		case errors.Is(
			err,
			shared.ErrEmailAlreadyExists,
		),
			errors.Is(
				err,
				shared.ErrInvitationAlreadyExists,
			):

			utils.WriteJSON(
				w,
				http.StatusConflict,
				shared.APIResponse{
					Success: false,
					Message: err.Error(),
				},
			)

			return

		//------------------------------------------------
		// Not Found Errors
		//------------------------------------------------

		case errors.Is(
			err,
			shared.ErrStateNotFound,
		),
			errors.Is(
				err,
				shared.ErrDistrictNotFound,
			),
			errors.Is(
				err,
				shared.ErrAcademyNotFound,
			):

			utils.WriteJSON(
				w,
				http.StatusNotFound,
				shared.APIResponse{
					Success: false,
					Message: err.Error(),
				},
			)

			return

		//------------------------------------------------
		// Internal Server Error
		//------------------------------------------------

		default:

			utils.WriteJSON(
				w,
				http.StatusInternalServerError,
				shared.APIResponse{
					Success: false,
					Message: "failed to create invitation",
				},
			)

			return
		}
	}

	//------------------------------------------------
	// Success Response
	//------------------------------------------------

	utils.WriteJSON(
		w,
		http.StatusCreated,
		shared.APIResponse{
			Success: true,
			Message: "invitation created successfully",
			Data:    response,
		},
	)
}

// GetInvitationsHandler returns all invitations.
//
// Authorization:
//
//	- super_admin
//	- state_admin
//	- district_admin
//	- academy_admin
//
// Responses:
//
//	- 200:
//	  Invitations fetched successfully.
//
//	- 401:
//	  Unauthorized.
//
//	- 500:
//	  Internal server error.
func GetInvitationsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Get Authenticated User
	// ----------------------------------------------------------

	userID := r.Context().
		Value(middleware.UserIDKey).
		(string)

	role := r.Context().
		Value(middleware.RoleNameKey).
		(string)

	// ----------------------------------------------------------
	// Parse Query Parameters
	// ----------------------------------------------------------

	query := InvitationsQueryParams{
		Page:   1,
		Limit:  10,
		SortBy: "created_at",
		Order:  "desc",
	}

	q := r.URL.Query()

	// Page
	if pageStr := q.Get("page"); pageStr != "" {

		page, err := strconv.Atoi(pageStr)

		if err == nil && page > 0 {
			query.Page = page
		}
	}

	// Limit
	if limitStr := q.Get("limit"); limitStr != "" {

		limit, err := strconv.Atoi(limitStr)

		if err == nil && limit > 0 && limit <= 100 {
			query.Limit = limit
		}
	}

	// Filters
	query.Search = q.Get("search")
	query.Status = q.Get("status")
	query.Role = q.Get("role")

	// Sorting
	if sortBy := q.Get("sort_by"); sortBy != "" {
		query.SortBy = sortBy
	}

	if order := q.Get("order"); order != "" {

		order = strings.ToLower(order)

		if order == "asc" || order == "desc" {
			query.Order = order
		}
	}

	// ----------------------------------------------------------
	// Get Invitations
	// ----------------------------------------------------------

	invitations, err := GetInvitationsService(
		r.Context(),
		userID,
		role,
		query,
	)

	if err != nil {

		switch err {

		case shared.ErrForbidden:

			utils.WriteJSON(
				w,
				http.StatusForbidden,
				shared.APIResponse{
					Success: false,
					Message: "forbidden",
				},
			)

			return
		}

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: "failed to fetch invitations",
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
			Message: "invitations fetched successfully",
			Data:    invitations,
		},
	)
}

// GetInvitationByIDHandler returns invitation details
// for the provided invitation ID.
//
// Authorization:
//
//	- super_admin
//	- state_admin
//	- district_admin
//	- academy_admin
//
// Path Params:
//
//	- id:
//	  Invitation ID.
//
// Responses:
//
//	- 200:
//	  Invitation fetched successfully.
//
//	- 400:
//	  Invalid invitation ID.
//
//	- 404:
//	  Invitation not found.
//
//	- 500:
//	  Internal server error.
func GetInvitationByIDHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Parse Path Parameter
	// ----------------------------------------------------------

	idStr := chi.URLParam(
		r,
		"id",
	)

	id, err := strconv.ParseInt(
		idStr,
		10,
		64,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid invitation id",
			},
		)

		return
	}

	// ----------------------------------------------------------
	// Get Authenticated User
	// ----------------------------------------------------------

	userID := r.Context().
		Value(middleware.UserIDKey).
		(string)

	role := r.Context().
		Value(middleware.RoleNameKey).
		(string)

	// ----------------------------------------------------------
	// Fetch Invitation
	// ----------------------------------------------------------

	invitation, err := GetInvitationByIDService(
		r.Context(),
		id,
		userID,
		role,
	)

	if err != nil {

		switch err {

		case shared.ErrInvitationNotFound:

			utils.WriteJSON(
				w,
				http.StatusNotFound,
				shared.APIResponse{
					Success: false,
					Message: "invitation not found",
				},
			)

			return

		case shared.ErrForbidden:

			utils.WriteJSON(
				w,
				http.StatusForbidden,
				shared.APIResponse{
					Success: false,
					Message: "forbidden",
				},
			)

			return
		}

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: "failed to fetch invitation",
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
			Message: "invitation fetched successfully",
			Data:    invitation,
		},
	)
}

// RevokeInvitationHandler revokes an existing invitation.
//
// Authorization:
//
//	- super_admin
//	- state_admin
//	- district_admin
//	- academy_admin
func RevokeInvitationHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Parse Path Parameter
	// ----------------------------------------------------------

	idStr := chi.URLParam(
		r,
		"id",
	)

	id, err := strconv.ParseInt(
		idStr,
		10,
		64,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid invitation id",
			},
		)

		return
	}

	// ----------------------------------------------------------
	// Get Authenticated User
	// ----------------------------------------------------------

	userID := r.Context().
		Value(middleware.UserIDKey).
		(string)

	role := r.Context().
		Value(middleware.RoleNameKey).
		(string)

	// ----------------------------------------------------------
	// Revoke Invitation
	// ----------------------------------------------------------

	err = RevokeInvitationService(
		r.Context(),
		id,
		userID,
		role,
	)

	if err != nil {

		switch err {

		case shared.ErrInvitationNotFound:

			utils.WriteJSON(
				w,
				http.StatusNotFound,
				shared.APIResponse{
					Success: false,
					Message: "invitation not found",
				},
			)

			return

		case shared.ErrForbidden:

			utils.WriteJSON(
				w,
				http.StatusForbidden,
				shared.APIResponse{
					Success: false,
					Message: "forbidden",
				},
			)

			return

		case shared.ErrInvitationAlreadyAccepted:

			utils.WriteJSON(
				w,
				http.StatusBadRequest,
				shared.APIResponse{
					Success: false,
					Message: "invitation already accepted",
				},
			)

			return

		case shared.ErrInvitationAlreadyRevoked:

			utils.WriteJSON(
				w,
				http.StatusBadRequest,
				shared.APIResponse{
					Success: false,
					Message: "invitation already revoked",
				},
			)

			return
		}

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: "failed to revoke invitation",
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
			Message: "invitation revoked successfully",
		},
	)
}

// AcceptInvitationHandler accepts an invitation
// using a valid invitation token.
//
// Public Route.
func AcceptInvitationHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	defer r.Body.Close()

	var input AcceptInvitationInput

	// ----------------------------------------------------------
	// Decode Request Body
	// ----------------------------------------------------------

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

	// ----------------------------------------------------------
	// Accept Invitation
	// ----------------------------------------------------------

	response, err := AcceptInvitationService(
		r.Context(),
		input,
	)

	if err != nil {

		fmt.Print("accept invitation service: ",err)

		switch err {

		case shared.ErrTokenRequired,
			shared.ErrFirstNameRequired,
			shared.ErrPasswordRequired,
			shared.ErrWeakPassword,
			shared.ErrInvalidPhoneNumber,
			shared.ErrInvalidInvitationScope:
			utils.WriteJSON(
				w,
				http.StatusBadRequest,
				shared.APIResponse{
					Success: false,
					Message: err.Error(),
				},
			)

			return

		case shared.ErrInvitationNotFound:

			utils.WriteJSON(
				w,
				http.StatusNotFound,
				shared.APIResponse{
					Success: false,
					Message: "invitation not found",
				},
			)

			return

		case shared.ErrInvitationExpired:

			utils.WriteJSON(
				w,
				http.StatusBadRequest,
				shared.APIResponse{
					Success: false,
					Message: "invitation expired",
				},
			)

			return

		case shared.ErrInvitationRevoked:

			utils.WriteJSON(
				w,
				http.StatusBadRequest,
				shared.APIResponse{
					Success: false,
					Message: "invitation revoked",
				},
			)

			return

		case shared.ErrInvitationAlreadyAccepted:

			utils.WriteJSON(
				w,
				http.StatusBadRequest,
				shared.APIResponse{
					Success: false,
					Message: "invitation already accepted",
				},
			)

			return

		case shared.ErrEmailAlreadyExists:

			utils.WriteJSON(
				w,
				http.StatusConflict,
				shared.APIResponse{
					Success: false,
					Message: "email already exists",
				},
			)

			return
		}

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: "failed to accept invitation",
			},
		)

		return
	}

	// ----------------------------------------------------------
	// Success Response
	// ----------------------------------------------------------

	utils.WriteJSON(
		w,
		http.StatusCreated,
		shared.APIResponse{
			Success: true,
			Message: "invitation accepted successfully",
			Data:    response,
		},
	)
}

// GetInvitationByTokenHandler fetches invitation
// details using an invitation token.
//
// Public Route.
func GetInvitationByTokenHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Get Token
	// ----------------------------------------------------------

	token := chi.URLParam(r, "token")

	if token == "" {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invitation token is required",
			},
		)

		return
	}

	// ----------------------------------------------------------
	// Get Invitation
	// ----------------------------------------------------------

	invitation, err := GetInvitationByTokenService(
		r.Context(),
		token,
	)

	if err != nil {

		switch err {

		case shared.ErrInvitationNotFound:

			utils.WriteJSON(
				w,
				http.StatusNotFound,
				shared.APIResponse{
					Success: false,
					Message: "invitation not found",
				},
			)

			return

		case shared.ErrInvitationExpired:

			utils.WriteJSON(
				w,
				http.StatusBadRequest,
				shared.APIResponse{
					Success: false,
					Message: "invitation has expired",
				},
			)

			return

		case shared.ErrInvitationRevoked:

			utils.WriteJSON(
				w,
				http.StatusBadRequest,
				shared.APIResponse{
					Success: false,
					Message: "invitation has been revoked",
				},
			)

			return
		}

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: "failed to fetch invitation",
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
			Message: "invitation fetched successfully",
			Data:    invitation,
		},
	)
}
