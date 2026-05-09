package districts

import (
	"net/http"

	"cobackend/internal/shared"
	"cobackend/internal/utils"

	"github.com/go-chi/chi/v5"

	"strconv"
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
	stateIDParam := chi.URLParam(r, "state_id")

	if stateIDParam == "" {
		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: "State ID is required",
		})
		return
	}

	stateID, err := strconv.Atoi(stateIDParam)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: "Invalid state ID",
		})
		return
	}

	districts, err := GetDistrictsByStateIdService(
		r.Context(),
		stateID,
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