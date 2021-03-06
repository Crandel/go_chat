package ws

import (
	"fmt"
	"log"
	"net/http"
	"time"

	cht "github.com/Crandel/go_chat/pkg/chatting"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WSHandler(chts cht.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		go chts.NewClient(conn, fmt.Sprintf("Anonymous%d", time.Now().UnixNano()))
	}
}
