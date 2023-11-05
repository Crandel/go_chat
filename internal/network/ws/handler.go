package ws

import (
	"log"
	"net/http"

	"github.com/Crandel/go_chat/internal/auth"
	cht "github.com/Crandel/go_chat/internal/chatting"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WSHandler(chts cht.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctxUser := ctx.Value(auth.AuthKey)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		if ctxUser != nil {
			ctxUserFull := ctxUser.(*auth.AuthUser)
			go chts.NewClient(conn, ctxUserFull.Nick)
		}

	}
}
