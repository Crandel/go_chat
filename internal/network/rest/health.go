package rest

import (
	"encoding/json"
	"net/http"
)

func HealthHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"status": "OK",
		}
		_ = json.NewEncoder(w).Encode(response)
	}
}
