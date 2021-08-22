package rest

import (
	"github.com/Crandel/go_chat/pkg/login"
	"github.com/Crandel/go_chat/pkg/signin"
	"github.com/gorilla/mux"
)

func InitHandlers(ls login.Service, sis signin.Service) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api", rootHandler()).Methods("GET")
	router.HandleFunc("/api/login", loginHandler(ls)).Methods("POST")
	router.HandleFunc("/api/signin", signinHandler(sis)).Methods("POST")
	return router
}
