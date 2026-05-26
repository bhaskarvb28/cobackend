package academyAdmin

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"strings"
// 	"io"

// 	"cobackend/internal/middleware"
// 	"cobackend/internal/shared"
// 	"cobackend/internal/utils"
// 	"cobackend/internal/validation"

// 	"github.com/go-chi/chi/v5"
// )

// func InviteAcademyAdminHandler(w http.ResponseWriter, r *http.Request) {
// 	var input InviteAcademyAdminInput

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

// 	inviteLink, err := InviteAcademyAdminService(r.Context(), input, authUserID)

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

// func UpdateAcademyAdminHandler(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) {

// 	var input UpdateAcademyAdminInput

// 	id := chi.URLParam(r, "id")

// 	if id == "" {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "academy admin id is required",
// 		})
// 		return
// 	}

// 	if !validation.IsValidUUID(id) {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "invalid academy admin id",
// 		})
// 		return
// 	}

// 	defer r.Body.Close()

// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()

// 	err := decoder.Decode(&input)
// 	if err != nil {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "invalid request body",
// 		})
// 		return
// 	}

// 	if decoder.Decode(&struct{}{}) != io.EOF {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "request body must contain only one JSON object",
// 		})
// 		return
// 	}

// 	if input.AcademyID != nil && *input.AcademyID <= 0 {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "invalid academy id",
// 		})
// 		return
// 	}

// 	err = UpdateAcademyAdminService(
// 		r.Context(),
// 		id,
// 		input,
// 	)

// 	if err != nil {
// 		utils.WriteJSON(w, http.StatusInternalServerError, shared.APIResponse{
// 			Success: false,
// 			Message: "failed to update academy admin",
// 		})
// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusOK, shared.APIResponse{
// 		Success: true,
// 		Message: "academy admin updated successfully",
// 	})
// }

// func GetAcademyAdminsHandler(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) {

// 	pageStr := r.URL.Query().Get("page")
// 	limitStr := r.URL.Query().Get("limit")
// 	search := strings.TrimSpace(
// 		r.URL.Query().Get("search"),
// 	)
// 	stateStr := r.URL.Query().Get("state_id")
// 	districtStr := r.URL.Query().Get("district_id")
// 	academyStr := r.URL.Query().Get("academy_id")
// 	sortBy := r.URL.Query().Get("sort_by")
// 	orderBy := r.URL.Query().Get("order_by")

// 	page := 1
// 	if pageStr != "" {

// 		parsed, err := strconv.Atoi(pageStr)
// 		if err != nil || parsed < 1 {

// 			utils.WriteJSON(
// 				w,
// 				http.StatusBadRequest,
// 				shared.APIResponse{
// 					Success: false,
// 					Message: "invalid page",
// 				},
// 			)

// 			return
// 		}

// 		page = parsed
// 	}

// 	limit := 10
// 	if limitStr != "" {

// 		if limitStr == "all" {

// 			limit = 0

// 		} else {

// 			parsed, err := strconv.Atoi(limitStr)
// 			if err != nil || parsed < 1 {

// 				utils.WriteJSON(
// 					w,
// 					http.StatusBadRequest,
// 					shared.APIResponse{
// 						Success: false,
// 						Message: "invalid limit",
// 					},
// 				)

// 				return
// 			}

// 			limit = parsed
// 		}
// 	}

// 	stateID := 0
// 	if stateStr != "" {

// 		parsed, err := strconv.Atoi(stateStr)
// 		if err != nil {

// 			utils.WriteJSON(
// 				w,
// 				http.StatusBadRequest,
// 				shared.APIResponse{
// 					Success: false,
// 					Message: "invalid state_id",
// 				},
// 			)

// 			return
// 		}

// 		stateID = parsed
// 	}

// 	districtID := 0
// 	if districtStr != "" {

// 		parsed, err := strconv.Atoi(districtStr)
// 		if err != nil {

// 			utils.WriteJSON(
// 				w,
// 				http.StatusBadRequest,
// 				shared.APIResponse{
// 					Success: false,
// 					Message: "invalid district_id",
// 				},
// 			)

// 			return
// 		}

// 		districtID = parsed
// 	}

// 	academyID := 0
// 	if academyStr != "" {

// 		parsed, err := strconv.Atoi(academyStr)
// 		if err != nil {

// 			utils.WriteJSON(
// 				w,
// 				http.StatusBadRequest,
// 				shared.APIResponse{
// 					Success: false,
// 					Message: "invalid academy_id",
// 				},
// 			)

// 			return
// 		}

// 		academyID = parsed
// 	}

// 	if sortBy == "" {
// 		sortBy = "first_name"
// 	}

// 	if orderBy == "" {
// 		orderBy = "asc"
// 	}

// 	_, exists := AllowedAcademyAdminSortFields[sortBy]

// 	if !exists {

// 		utils.WriteJSON(
// 			w,
// 			http.StatusBadRequest,
// 			shared.APIResponse{
// 				Success: false,
// 				Message: "invalid sort_by field",
// 			},
// 		)

// 		return
// 	}

// 	orderBy = strings.ToUpper(orderBy)

// 	if orderBy != "ASC" && orderBy != "DESC" {

// 		utils.WriteJSON(
// 			w,
// 			http.StatusBadRequest,
// 			shared.APIResponse{
// 				Success: false,
// 				Message: "invalid order_by value",
// 			},
// 		)

// 		return
// 	}

// 	query := GetAcademyAdminsQuery{
// 		Page:       page,
// 		Limit:      limit,
// 		Search:     search,
// 		StateID:    stateID,
// 		DistrictID: districtID,
// 		AcademyID:  academyID,
// 		SortBy:     sortBy,
// 		OrderBy:    orderBy,
// 	}

// 	result, err := GetAcademyAdminsService(
// 		r.Context(),
// 		query,
// 	)

// 	if err != nil {
// 		fmt.Print(err)

// 		utils.WriteError(
// 			w,
// 			err,
// 			"failed to fetch academy admins",
// 		)

// 		return
// 	}

// 	utils.WriteJSON(
// 		w,
// 		http.StatusOK,
// 		shared.APIResponse{
// 			Success: true,
// 			Message: "academy admins fetched successfully",
// 			Data:    result,
// 		},
// 	)
// }

// func GetAcademyAdminByIDHandler(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) {

// 	profileID := chi.URLParam(r, "profile_id")

// 	if profileID == "" {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "profile_id is required",
// 		})
// 		return
// 	}

// 	if !validation.IsValidUUID(profileID) {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "invalid profile_id",
// 		})
// 		return
// 	}

// 	admin, err := GetAcademyAdminByIDService(
// 		r.Context(),
// 		profileID,
// 	)

// 	if err != nil {
// 		utils.WriteError(
// 			w,
// 			err,
// 			"failed to fetch academy admin",
// 		)
// 		return
// 	}

// 	utils.WriteJSON(
// 		w,
// 		http.StatusOK,
// 		shared.APIResponse{
// 			Success: true,
// 			Message: "academy admin fetched successfully",
// 			Data:    admin,
// 		},
// 	)
// }

// func DeleteAcademyAdminHandler(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) {

// 	id := chi.URLParam(r, "id")

// 	if id == "" {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "academy admin id is required",
// 		})
// 		return
// 	}

// 	if !validation.IsValidUUID(id) {
// 		utils.WriteJSON(w, http.StatusBadRequest, shared.APIResponse{
// 			Success: false,
// 			Message: "invalid academy admin id",
// 		})
// 		return
// 	}

// 	authUserID, ok := r.Context().Value(
// 		middleware.UserIDKey,
// 	).(string)

// 	if !ok || authUserID == "" {
// 		utils.WriteJSON(w, http.StatusUnauthorized, shared.APIResponse{
// 			Success: false,
// 			Message: "unauthorized",
// 		})
// 		return
// 	}

// 	err := DeleteAcademyAdminService(
// 		r.Context(),
// 		authUserID,
// 		id,
// 	)

// 	if err != nil {

// 		if errors.Is(err, shared.ErrAcademyAdminNotFound) {

// 			utils.WriteJSON(w, http.StatusNotFound, shared.APIResponse{
// 				Success: false,
// 				Message: "academy admin not found",
// 			})

// 			return
// 		}

// 		if errors.Is(err, shared.ErrForbidden) {

// 			utils.WriteJSON(w, http.StatusForbidden, shared.APIResponse{
// 				Success: false,
// 				Message: "you are not allowed to delete this academy admin",
// 			})

// 			return
// 		}

// 		utils.WriteJSON(w, http.StatusInternalServerError, shared.APIResponse{
// 			Success: false,
// 			Message: "failed to delete academy admin",
// 		})

// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusOK, shared.APIResponse{
// 		Success: true,
// 		Message: "academy admin deleted successfully",
// 	})
// }