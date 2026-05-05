package auth

import (
	"encoding/json"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var input RegisterInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err := Register(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Write([]byte("user created"))
}