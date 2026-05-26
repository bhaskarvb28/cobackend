package academyCoach

// import (
// 	"encoding/json"
	
// 	"cobackend/internal/utils"
// 	"cobackend/internal/middleware"
// 	"cobackend/internal/shared"

// 	"net/http"
// 	"errors"
// 	"fmt"
// )

// func InviteAcademyCoachHandler(w http.ResponseWriter, r *http.Request) {
// 	var input InviteAcademyCoachInput

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
// 		Value(middleware.UserIDKey).
// 		(string)

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

// 	inviteLink, err := InviteAcademyCoachService(r.Context(), input, authUserID)

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
// 				Message: "failed to create academy admin invite",
// 			},
// 		)

// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusCreated, shared.APIResponse{
// 		Success: true,
// 		Message: "academy admin invitation created successfully",
// 		Data: inviteLink,
// 	})

// }