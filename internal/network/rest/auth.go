package rest

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Crandel/go_chat/internal/auth"
	"gitlab.com/greyxor/slogor"
)

func LoginHandler(athS auth.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var lu auth.LoginUser
		if err := json.NewDecoder(r.Body).Decode(&lu); err != nil {
			slog.Error("Error during decoding", slogor.Err(err))
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		slog.Debug("Inside login ", slog.String("user", lu.String()))
		response, err := athS.LoginUser(lu)
		if err != nil {
			slog.Error("Error during login", slogor.Err(err))
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		ctx := r.Context()
		ctxUser := ctx.Value(auth.AuthKey)

		if ctxUser != nil {
			ctxUserFull := ctxUser.(*auth.AuthUser)
			ctxUserFull.Nick = lu.Nick
			ctxUserFull.Token = response.Token
			slog.Debug("CTX User", slog.String("user", ctxUserFull.String()))
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Error("Error during login", slogor.Err(err))
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
	}
}

func SigninHandler(athS auth.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var su auth.SigninUser
		if err := json.NewDecoder(r.Body).Decode(&su); err != nil {
			slog.Error("Error while json decoding", slogor.Err(err))
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		response, err := athS.SigninUser(su)
		if err != nil {
			slog.Error("Error during signing", slogor.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		ctxUser := ctx.Value(auth.AuthKey)
		if ctxUser != nil {
			ctxUserFull := ctxUser.(*auth.AuthUser)
			ctxUserFull.Nick = su.Nick
			ctxUserFull.Token = response.Token
			slog.Debug("CTX User", slog.String("user", ctxUserFull.String()))
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Error("Error during signing", slogor.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
