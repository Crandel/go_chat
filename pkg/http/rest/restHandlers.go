package rest

import "github.com/gorilla/mux"

func InitHandlers() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api", rootHandler()).Methods("GET")
	return router
}
