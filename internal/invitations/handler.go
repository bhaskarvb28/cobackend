package invitations

import (
	"net/http"
	"cobackend/internal/utils"

	"strings"

	"cobackend/internal/shared"

	"github.com/go-chi/chi/v5"

)

func GetInvitationByTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimSpace(
		chi.URLParam(r, "token"),
	)

	if token == "" {
		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: "Token Required",
		})

		return
	}

	result, err := GetInvitationByTokenService(
		r.Context(),
		token,
	)

	if err != nil {

		utils.WriteError(
			w,
			err,
			"failed to fetch invitation",
		)

		return
	}

	utils.WriteJSON(
		w,
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "invitation fetched successfully",
			Data:    result,
		},
	)
}