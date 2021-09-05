package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Crandel/go_chat/pkg/signin"
)

func signinHandler(sis signin.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var su signin.User
		if err := json.NewDecoder(r.Body).Decode(&su); err != nil {
			log.Println("Error while json decoding", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		response, err := sis.SigninUser(su)
		if err != nil {
			log.Println("Error after signing", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(response)
	}
}
