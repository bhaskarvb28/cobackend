package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"cobackend/internal/shared"
	"cobackend/internal/utils"
)

// LoginHandler authenticates a user and
// returns authentication tokens.
//
// Responses:
//   - 200:
//     Login successful.
//
//   - 400:
//     Invalid request body.
//
//   - 401:
//     Invalid email or password.
func LoginHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	var input LoginInput

	// ------------------------------------------------------------------
	// Decode Request Body
	// ------------------------------------------------------------------

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&input)

	if err != nil {

		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: shared.ErrInvalidRequestBody.Error(),
		})

		return
	}

	// ------------------------------------------------------------------
	// Normalize Input
	// ------------------------------------------------------------------

	input.Email = strings.TrimSpace(input.Email)
	input.Password = strings.TrimSpace(input.Password)

	// ------------------------------------------------------------------
	// Validate Required Fields
	// ------------------------------------------------------------------

	if input.Email == "" || input.Password == "" {

		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
			Success: false,
			Message: "email and password are required",
		})

		return
	}

	// ------------------------------------------------------------------
	// Execute Authentication
	// ------------------------------------------------------------------

	response, err := LoginService(
		r.Context(),
		input,
	)

	if err != nil {

		fmt.Print(err)

		if errors.Is(err, shared.ErrInvalidCredentials) {

			utils.WriteJSON(w, http.StatusUnauthorized, shared.APIResponse{
				Success: false,
				Message: shared.ErrInvalidCredentials.Error(),
			})

			return
		}

		utils.WriteJSON(w, http.StatusInternalServerError, shared.APIResponse{
			Success: false,
			Message: shared.ErrInternalServerError.Error(),
		})

		return
	}

	// ------------------------------------------------------------------
	// Success Response
	// ------------------------------------------------------------------

	utils.WriteJSON(w, http.StatusOK, shared.APIResponse{
		Success: true,
		Message: "login successful",
		Data:    response,
	})
}

// func AcceptInviteHandler(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) {

// 	var input AcceptInvitationInput

// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()
// 	defer r.Body.Close()

// 	err := decoder.Decode(&input)

// 	if err != nil {
// 		fmt.Print(err)

// 		utils.WriteJSON(
// 			w,
// 			http.StatusBadRequest,
// 			shared.APIResponse{
// 				Success: false,
// 				Message: "invalid request body",
// 			},
// 		)

// 		return
// 	}

// 	err = AcceptInvitationService(
// 		r.Context(),
// 		input,
// 	)

// 	if err != nil {
// 		fmt.Print(err)

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


// 		utils.WriteJSON(
// 			w,
// 			http.StatusInternalServerError,
// 			shared.APIResponse{
// 				Success: false,
// 				Message: "internal server error",
// 			},
// 		)

// 		return
// 	}

// 	utils.WriteJSON(
// 		w,
// 		http.StatusCreated,
// 		shared.APIResponse{
// 			Success: true,
// 			Message: "account setup completed successfully",
// 		},
// 	)
// }