package rest

import (
	"net/http"

	"github.com/Crandel/go_chat/pkg/adding"
	"github.com/Crandel/go_chat/pkg/auth"
	"github.com/Crandel/go_chat/pkg/reading"
	"github.com/gorilla/mux"
)

type userHandler struct{}

func InitHandlers(
	aths auth.Service,
	adds adding.Service,
	rdns reading.Service,
) *mux.Router {
	r := mux.NewRouter()
	api_router := r.PathPrefix("/api").Subrouter()
	user_router := api_router.PathPrefix("/users").Subrouter()
	user_router.HandleFunc("/login", loginHandler(aths)).Methods(http.MethodPost)
	user_router.HandleFunc("/signin", signinHandler(aths)).Methods(http.MethodPost)
	user_router.HandleFunc("", listUsersHandler(rdns)).Methods(http.MethodGet)
	user_router.HandleFunc("/{user_id}", getUserHandler(rdns)).Methods(http.MethodGet)
	room_router := api_router.PathPrefix("/rooms").Subrouter()
	room_router.HandleFunc("", listRoomsHandler(rdns)).Methods(http.MethodGet)
	room_router.HandleFunc("", addRoomHandler(adds)).Methods(http.MethodPost)
	room_router.HandleFunc("/{room_id}", getRoomHandler(rdns)).Methods(http.MethodGet)
	return r
}
