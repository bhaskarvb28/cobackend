package districts

import (
	"net/http"
	"strconv"
	"strings"

	"cobackend/internal/shared"
	"cobackend/internal/utils"

	"github.com/go-chi/chi/v5"
)

func GetDistrictsHandler(w http.ResponseWriter, r *http.Request) {
	districts, err := GetDistrictsService(r.Context())
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, shared.APIResponse{
			Success: false,
			Message: "Failed to fetch districts",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, shared.APIResponse{
		Success: true,
		Message: "Districts Fetched Successfully",
		Data: districts,
	})
}

func GetDistrictsByStateIdHandler(w http.ResponseWriter, r *http.Request) {
	stateIDStr := chi.URLParam(r, "state_id")

	search := strings.TrimSpace(
		r.URL.Query().Get("search"),
	)
	order := strings.TrimSpace(
		r.URL.Query().Get("order"),
	)

	if stateIDStr == "" {
		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: "State ID is required",
		})
		return
	}

	stateID, err := strconv.Atoi(stateIDStr)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: "Invalid state ID",
		})
		return
	}

	queryParams := GetDistrictQueryParams {
		Search: search,
		Order: order,
	}

	districts, err := GetDistrictsByStateIdService(
		r.Context(),
		stateID,
		queryParams,
	)

	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, shared.APIResponse{
			Success: false,
			Message: "Failed to fetch districts",
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, shared.APIResponse{
		Success: true,
		Message: "Districts fetched successfully",
		Data:    districts,
	})
}