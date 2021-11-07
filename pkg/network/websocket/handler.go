package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func WSHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()
		fmt.Println("Client was successfuly connected")
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Println(string(p))
			fmt.Println(messageType)
			if err := conn.WriteMessage(messageType, p); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
