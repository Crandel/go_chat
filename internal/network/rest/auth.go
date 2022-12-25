package rest

import (
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
		ctx := r.Context()
		ctxUser := ctx.Value(auth.AuthKey)
		if ctxUser != nil {
			authUser := &auth.CtxUser{
				Nick:  lu.Nick,
				Token: response.Token,
			}
			ctxUser = authUser
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Println("Error during login", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
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

		ctx := r.Context()
		ctxUser := ctx.Value(auth.AuthKey)
		if ctxUser != nil {
			authUser := &auth.CtxUser{
				Nick:  su.Nick,
				Token: response.Token,
			}
			ctxUser = authUser
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Println("Error during signing", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
