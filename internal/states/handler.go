package states

import (
	"net/http"

	"cobackend/internal/shared"
	"cobackend/internal/utils"
)

func GetStatesHandler(w http.ResponseWriter, r *http.Request) {
	states, err := GetStatesService(r.Context())
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, shared.APIResponse{
			Success: false,
			Message: "Failed to fetch states",
		})
		return
	}


	utils.WriteJSON(w, http.StatusOK, shared.APIResponse{
		Success: true,
		Message: "States fetched successfully",
		Data: states,
	})
}