package districts

import (
	"net/http"
	"encoding/json"
)

func GetDistrictsHandler(w http.ResponseWriter, r *http.Request) {
	districts, err := GetDistrictsService(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch districts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(districts)

}
