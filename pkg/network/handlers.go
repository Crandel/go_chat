package network

import (
	"net/http"

	"github.com/Crandel/go_chat/pkg/adding"
	"github.com/Crandel/go_chat/pkg/auth"
	"github.com/Crandel/go_chat/pkg/network/rest"
	"github.com/Crandel/go_chat/pkg/reading"
	"github.com/gorilla/mux"
)

func InitHandlers(
	aths auth.Service,
	adds adding.Service,
	rdns reading.Service,
) *mux.Router {
	r := mux.NewRouter()
	api_router := r.PathPrefix("/api").Subrouter()
	user_router := api_router.PathPrefix("/users").Subrouter()
	user_router.HandleFunc("/login", rest.LoginHandler(aths)).Methods(http.MethodPost)
	user_router.HandleFunc("/signin", rest.SigninHandler(aths)).Methods(http.MethodPost)
	user_router.HandleFunc("", rest.ListUsersHandler(rdns)).Methods(http.MethodGet)
	user_router.HandleFunc("/{user_id}", rest.GetUserHandler(rdns)).Methods(http.MethodGet)
	room_router := api_router.PathPrefix("/rooms").Subrouter()
	room_router.HandleFunc("", rest.ListRoomsHandler(rdns)).Methods(http.MethodGet)
	room_router.HandleFunc("", rest.AddRoomHandler(adds)).Methods(http.MethodPost)
	room_router.HandleFunc("/{room_id}", rest.GetRoomHandler(rdns)).Methods(http.MethodGet)
	return r
}
