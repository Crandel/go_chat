package network

import (
	"net/http"

	"github.com/Crandel/go_chat/pkg/adding"
	"github.com/Crandel/go_chat/pkg/auth"
	"github.com/Crandel/go_chat/pkg/chatting"
	hl "github.com/Crandel/go_chat/pkg/network/html"
	"github.com/Crandel/go_chat/pkg/network/rest"
	"github.com/Crandel/go_chat/pkg/network/ws"
	"github.com/Crandel/go_chat/pkg/reading"
	"github.com/gorilla/mux"
)

func InitHandlers(
	aths auth.Service,
	adds adding.Service,
	rdns reading.Service,
	chts chatting.Service,
) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", hl.RootHandler())
	r.HandleFunc("/ws", ws.WSHandler(chts))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", hl.StaticHandler()))
	apiRouter := r.PathPrefix("/api").Subrouter()
	userRouter := apiRouter.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("/login", rest.LoginHandler(aths)).Methods(http.MethodPost)
	userRouter.HandleFunc("/signin", rest.SigninHandler(aths)).Methods(http.MethodPost)
	userRouter.HandleFunc("", rest.ListUsersHandler(rdns)).Methods(http.MethodGet)
	userRouter.HandleFunc("/{user_id}", rest.GetUserHandler(rdns)).Methods(http.MethodGet)
	roomRouter := apiRouter.PathPrefix("/rooms").Subrouter()
	roomRouter.HandleFunc("", rest.ListRoomsHandler(rdns)).Methods(http.MethodGet)
	roomRouter.HandleFunc("", rest.AddRoomHandler(adds)).Methods(http.MethodPost)
	roomRouter.HandleFunc("/{room_id}", rest.GetRoomHandler(rdns)).Methods(http.MethodGet)
	return r
}
