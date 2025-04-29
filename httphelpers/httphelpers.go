package httphelpers

import (
	"encoding/json"
	"net/http"
	"time"
)

func RespondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, map[string]string{"error": message})
}

func ParseDateParam(param string) (time.Time, error) {
	return time.Parse("2006-01-02", param)
}
