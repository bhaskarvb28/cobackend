package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"cobackend/internal/utils"
	"cobackend/internal/shared"

	"errors"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var input LoginInput

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&input)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Trim whitespace
	input.Email = strings.TrimSpace(input.Email)
	input.Password = strings.TrimSpace(input.Password)

	// Validate required fields
	if input.Email == "" || input.Password == "" {
		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: "Email and password are required",
		})
		return
	}

	// Login service
	response, err := Login(r.Context(), input)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, shared.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Success response
	utils.WriteJSON(w, http.StatusOK, shared.APIResponse{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

func AcceptInviteHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	var input AcceptInvitationInput

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

	err = AcceptInvitationService(
		r.Context(),
		input,
	)

	if err != nil {

		var apiErr *shared.APIError

		if errors.As(err, &apiErr) {

			utils.WriteJSON(
				w,
				apiErr.StatusCode,
				shared.APIResponse{
					Success: false,
					Message: apiErr.Message,
				},
			)

			return
		}

		utils.WriteJSON(
			w,
			http.StatusInternalServerError,
			shared.APIResponse{
				Success: false,
				Message: "internal server error",
			},
		)

		return
	}

	utils.WriteJSON(
		w,
		http.StatusCreated,
		shared.APIResponse{
			Success: true,
			Message: "account setup completed successfully",
		},
	)
}