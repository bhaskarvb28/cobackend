package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"cobackend/internal/utils"
	"cobackend/internal/shared"
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