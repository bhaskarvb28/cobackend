package states

import (
	"net/http"

	"cobackend/internal/shared"
	"cobackend/internal/utils"

	"strings"


	"fmt"
)

func GetStatesHandler(w http.ResponseWriter, r *http.Request) {

	search := strings.TrimSpace(
		r.URL.Query().Get("search"),
	)	
	order := strings.TrimSpace(
		r.URL.Query().Get("order"),
	)

	queryParams := GetStatesQueryParams {
		Search: search,
		Order: order,
	}

	states, err := GetStatesService(r.Context(), queryParams)
	if err != nil {
		fmt.Println(err)
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