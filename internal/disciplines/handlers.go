package disciplines

import (
	"net/http"

	"cobackend/internal/shared"
	"cobackend/internal/utils"
)

func GetDisciplinesHandler(w http.ResponseWriter, r *http.Request) {
	disciplines, err := GetDisciplinesService(r.Context());

	if err != nil {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: err.Error(),
			},
		)

		return
	}

	utils.WriteJSON(
		w,
		http.StatusCreated,
		shared.APIResponse{
			Success: true,
			Message: "disciplines fetched successfully",
			Data:    disciplines,
		},
	)

}