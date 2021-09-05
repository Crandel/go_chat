package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Crandel/go_chat/pkg/login"
)

func loginHandler(ls login.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var lu login.User
		if err := json.NewDecoder(r.Body).Decode(&lu); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		token, err := ls.LoginUser(lu)
		if err != nil {
			log.Println("Failed during login", err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(token)
	}
}
