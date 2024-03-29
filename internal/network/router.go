package network

import (
	"net/http"

	"github.com/Crandel/go_chat/internal/adding"
	"github.com/Crandel/go_chat/internal/auth"
	"github.com/Crandel/go_chat/internal/chatting"
	hl "github.com/Crandel/go_chat/internal/network/html"
	"github.com/Crandel/go_chat/internal/network/rest"
	"github.com/Crandel/go_chat/internal/network/ws"
	"github.com/Crandel/go_chat/internal/reading"
	"github.com/gorilla/mux"
)

func NewRouter(
	aths auth.Service,
	adds adding.Service,
	rdns reading.Service,
	chts chatting.Service,
) *mux.Router {
	authMiddleware := NewAuthMiddleware(aths)
	r := mux.NewRouter()
	r.HandleFunc("/", hl.RootHandler())
	r.HandleFunc("/health", rest.HealthHandler())
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", hl.StaticHandler()))
	apiRouter := r.PathPrefix("/api").Subrouter()
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/login", rest.LoginHandler(aths)).Methods(http.MethodPost)
	authRouter.HandleFunc("/signin", rest.SigninHandler(aths)).Methods(http.MethodPost)
	authRouter.Use(authMiddleware.Populate)
	userRouter := apiRouter.PathPrefix("/users").Subrouter()
	userRouter.HandleFunc("", rest.ListUsersHandler(rdns)).Methods(http.MethodGet)
	userRouter.HandleFunc("/{user_id}", rest.GetUserHandler(rdns)).Methods(http.MethodGet)
	userRouter.Use(authMiddleware.Middleware)
	chatRouter := apiRouter.PathPrefix("/ws").Subrouter()
	chatRouter.HandleFunc("", ws.WSHandler(chts))
	chatRouter.Use(authMiddleware.Middleware)
	roomRouter := apiRouter.PathPrefix("/rooms").Subrouter()
	roomRouter.HandleFunc("", rest.ListRoomsHandler(rdns)).Methods(http.MethodGet)
	roomRouter.HandleFunc("", rest.AddRoomHandler(adds)).Methods(http.MethodPost)
	roomRouter.HandleFunc("/{room_id}", rest.GetRoomHandler(rdns)).Methods(http.MethodGet)
	roomRouter.Use(authMiddleware.Middleware)
	return r
}
