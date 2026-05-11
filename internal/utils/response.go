package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"cobackend/internal/shared"
)

func WriteJSON(
	w http.ResponseWriter,
	status int,
	data interface{},
) {

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

func WriteError(
	w http.ResponseWriter,
	err error,
	defaultMessage string,
) {

	var apiErr *shared.APIError

	if errors.As(err, &apiErr) {

		WriteJSON(
			w,
			apiErr.StatusCode,
			shared.APIResponse{
				Success: false,
				Message: apiErr.Message,
			},
		)

		return
	}

	WriteJSON(
		w,
		http.StatusInternalServerError,
		shared.APIResponse{
			Success: false,
			Message: defaultMessage,
		},
	)
}