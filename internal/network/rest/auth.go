package rest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Crandel/go_chat/internal/auth"
)

func LoginHandler(athS auth.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var lu auth.LoginUser
		if err := json.NewDecoder(r.Body).Decode(&lu); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		response, err := athS.LoginUser(lu)
		if err != nil {
			log.Println("Error during login", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		r.WithContext(context.WithValue(r.Context(), "nick", lu.Nick))
		r.WithContext(context.WithValue(r.Context(), "token", response.Token))
		json.NewEncoder(w).Encode(response)
	}
}

func SigninHandler(athS auth.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var su auth.SigninUser
		if err := json.NewDecoder(r.Body).Decode(&su); err != nil {
			log.Println("Error while json decoding", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		response, err := athS.SigninUser(su)
		if err != nil {
			log.Println("Error during signing", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		r.WithContext(context.WithValue(r.Context(), "nick", su.Nick))
		r.WithContext(context.WithValue(r.Context(), "token", response.Token))

		json.NewEncoder(w).Encode(response)
	}
}
