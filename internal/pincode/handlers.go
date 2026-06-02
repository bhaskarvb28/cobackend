package pincode

import (
	"cobackend/internal/shared"
	"cobackend/internal/utils"
	"net/http"
)

func GetPincodesHandler(w http.ResponseWriter, r *http.Request) {

	pincodes, err := GetPincodesService(r.Context())

	if err != nil {
		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: "failed to fetch pincodes",
			},
		)
	}

	utils.WriteJSON(
		w,
		http.StatusOK,
		shared.APIResponse {
			Success: true,
			Message: "fetched pincodes successfully",
			Data: pincodes,
		},
	)
}