package session

import (
	"net/http"
	"encoding/json"
	"cobackend/internal/utils"
	"cobackend/internal/shared"

	"cobackend/internal/middleware"
)
// handler.go

func StartPracticeSessionHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	defer r.Body.Close()

	var input StartPracticeSessionInput

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

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		(string)

	response, err := StartPracticeSessionService(
		r.Context(),
		authUserID,
		input,
	)

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
			Message: "practice session started successfully",
			Data:    response,
		},
	)
}