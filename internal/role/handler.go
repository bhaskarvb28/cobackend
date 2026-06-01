package role

import (
	"net/http"

	"cobackend/internal/middleware"
	"cobackend/internal/shared"
	"cobackend/internal/utils"
)

func GetInvitableRoleOptionsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// ----------------------------------------------------------
	// Get Authenticated User Role
	// ----------------------------------------------------------

	roleCode, ok := r.Context().
		Value(middleware.RoleNameKey).
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

	// ----------------------------------------------------------
	// Fetch Invitable Role Options
	// ----------------------------------------------------------

	roles, err := GetInvitableRoleOptionsService(
		r.Context(),
		roleCode,
	)

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: "failed to fetch invitable role options",
			},
		)

		return
	}

	// ----------------------------------------------------------
	// Response
	// ----------------------------------------------------------

	utils.WriteJSON(
		w,
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "invitable role options fetched successfully",
			Data: map[string]interface{}{
				"roles": roles,
			},
		},
	)
}
