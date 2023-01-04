package rest

import (
	"encoding/json"
	"net/http"

	lg "github.com/Crandel/go_chat/internal/logging"

	"github.com/Crandel/go_chat/internal/auth"
)

var log = lg.InitLogger()

func LoginHandler(athS auth.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var lu auth.LoginUser
		if err := json.NewDecoder(r.Body).Decode(&lu); err != nil {
			log.Log(lg.Warning, "Error during decoding", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		log.Log(lg.Debug, "Inside login ", lu)
		response, err := athS.LoginUser(lu)
		if err != nil {
			log.Log(lg.Warning, "Error during login", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		ctx := r.Context()
		ctxUser := ctx.Value(auth.AuthKey)

		if ctxUser != nil {
			ctxUserFull := ctxUser.(*auth.AuthUser)
			ctxUserFull.Nick = lu.Nick
			ctxUserFull.Token = response.Token
			log.Logf(lg.Debug, "CTX User full: %s\n", ctxUserFull)
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Log(lg.Debug, "Error during login", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	}
}

func SigninHandler(athS auth.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var su auth.SigninUser
		if err := json.NewDecoder(r.Body).Decode(&su); err != nil {
			log.Log(lg.Debug, "Error while json decoding", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		response, err := athS.SigninUser(su)
		if err != nil {
			log.Log(lg.Debug, "Error during signing", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		ctxUser := ctx.Value(auth.AuthKey)
		if ctxUser != nil {
			ctxUserFull := ctxUser.(*auth.AuthUser)
			ctxUserFull.Nick = su.Nick
			ctxUserFull.Token = response.Token
			log.Logf(lg.Debug, "CTX User full: %s\n", ctxUserFull)
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Log(lg.Debug, "Error during signing", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
