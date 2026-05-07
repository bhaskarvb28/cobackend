package states

import (
	"encoding/json"
	"net/http"
)

func GetStatesHandler(w http.ResponseWriter, r *http.Request) {
	states, err := GetStatesService(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch states", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(states)
}