package rest

import (
	"encoding/json"
	"net/http"

	"github.com/Crandel/go_chat/pkg/signin"
)

func signinHandler(sis signin.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var su signin.User
		if err := json.NewDecoder(r.Body).Decode(&su); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		response, err := sis.SigninUser(su)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(response)
	}
}
