package academy

import (
	"net/http"
	"encoding/json"

	"cobackend/internal/utils"
	"cobackend/internal/shared"

	"strings"
)

func CreateAcademyHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	var input CreateAcademyInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	defer r.Body.Close()

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

	if strings.TrimSpace(input.Name) == "" {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "academy name is required",
			},
		)

		return
	}

	if input.DistrictID <= 0 {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "valid district_id is required",
			},
		)

		return
	}

	if strings.TrimSpace(input.Address) == "" {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "address is required",
			},
		)

		return
	}

	err = CreateAcademyService(
		r.Context(),
		input,
	)

	if err != nil {

		utils.WriteError(
			w,
			err,
			"failed to create academy",
		)

		return
	}

	utils.WriteJSON(
		w,
		http.StatusCreated,
		shared.APIResponse{
			Success: true,
			Message: "academy created successfully",
		},
	)
}