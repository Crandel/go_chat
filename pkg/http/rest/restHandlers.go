package rest

import (
	"net/http"

	"github.com/Crandel/go_chat/pkg/adding"
	"github.com/Crandel/go_chat/pkg/login"
	"github.com/Crandel/go_chat/pkg/reading"
	"github.com/Crandel/go_chat/pkg/signin"
	"github.com/gorilla/mux"
)

type userHandler struct{}

func InitHandlers(
	ls login.Service,
	sis signin.Service,
	as adding.Service,
	rs reading.Service,
) *mux.Router {
	r := mux.NewRouter()
	api_router := r.PathPrefix("/api").Subrouter()
	api_router.HandleFunc("/login", loginHandler(ls)).Methods(http.MethodPost)
	api_router.HandleFunc("/signin", signinHandler(sis)).Methods(http.MethodPost)
	user_router := api_router.PathPrefix("/users").Subrouter()
	user_router.HandleFunc("", listUsersHandler(rs)).Methods(http.MethodGet)
	user_router.HandleFunc("/{user_id}", getUserHandler(rs)).Methods(http.MethodGet)
	room_router := api_router.PathPrefix("/rooms").Subrouter()
	room_router.HandleFunc("", listRoomsHandler(rs)).Methods(http.MethodGet)
	room_router.HandleFunc("", addRoomHandler(as)).Methods(http.MethodPost)
	room_router.HandleFunc("/{room_id}", getRoomHandler(rs)).Methods(http.MethodGet)
	return r
}
