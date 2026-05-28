package player

import (
	"cobackend/internal/middleware"
	"cobackend/internal/shared"
	"cobackend/internal/utils"
	"strconv"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
	"net/http"
)

// func InvitePlayerHandler(w http.ResponseWriter, r *http.Request) {
// 	var input InvitePlayerInput

// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()

// 	defer r.Body.Close()

// 	err := decoder.Decode(&input)
// 	if err != nil {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "invalid request body",
// 		})
// 		return
// 	}

// 	authUserID, ok := r.Context().
// 		Value(middleware.UserIDKey).(string)

// 	if !ok {
// 		utils.WriteJSON(
// 			w,
// 			http.StatusUnauthorized,
// 			shared.APIResponse{
// 				Success: false,
// 				Message: "unauthorized",
// 			},
// 		)
// 		return
// 	}

// 	inviteLink, err := InvitePlayerService(r.Context(), input, authUserID)

// 	if err != nil {

// 		var apiErr *shared.APIError

// 		if errors.As(err, &apiErr) {

// 			utils.WriteJSON(
// 				w,
// 				apiErr.StatusCode,
// 				shared.APIResponse{
// 					Success: false,
// 					Message: apiErr.Message,
// 				},
// 			)

// 			return
// 		}

// 		fmt.Print(err)

// 		utils.WriteJSON(
// 			w,
// 			http.StatusInternalServerError,
// 			shared.APIResponse{
// 				Success: false,
// 				Message: "failed to create player invite",
// 			},
// 		)

// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusCreated, shared.APIResponse{
// 		Success: true,
// 		Message: "player invitation created successfully",
// 		Data:    inviteLink,
// 	})
// }

func GetAvailableShootingEventsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	disciplineIDParam := r.URL.Query().Get(
		"discipline_id",
	)

	disciplineID, err := strconv.ParseInt(
		disciplineIDParam,
		10,
		16,
	)

	if err != nil || disciplineID <= 0 {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid discipline id",
			},
		)

		return
	}

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		(string)

	response, err := GetAvailableShootingEventsService(
		r.Context(),
		authUserID,
		int16(disciplineID),
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
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "shooting events fetched successfully",
			Data:    response,
		},
	)
}

func GetCompatibleBuildingsHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	shootingEventIDParam := r.URL.Query().Get(
		"shooting_event_id",
	)

	shootingEventID, err := strconv.ParseInt(
		shootingEventIDParam,
		10,
		16,
	)

	if err != nil || shootingEventID <= 0 {

		utils.WriteJSON(
			w,
			http.StatusBadRequest,
			shared.APIResponse{
				Success: false,
				Message: "invalid shooting event id",
			},
		)

		return
	}

	authUserID := r.Context().
		Value(middleware.UserIDKey).
		(string)

	response, err := GetCompatibleBuildingsService(
		r.Context(),
		authUserID,
		int16(shootingEventID),
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
		http.StatusOK,
		shared.APIResponse{
			Success: true,
			Message: "compatible buildings fetched successfully",
			Data:    response,
		},
	)
}

