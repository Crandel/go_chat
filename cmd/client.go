package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var done chan interface{}
var interrupt chan os.Signal

func msgHandler(conn *websocket.Conn) {
	defer close(done)
}

func main() {
	done = make(chan interface{})
	interrupt = make(chan os.Signal)

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT
	socketUrl := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketUrl, nil)
	if err != nil {
		log.Fatal("Could not connect to WebSocker server '"+socketUrl+"'.", err)
	}
	defer conn.Close()
	go msgHandler(conn)

	for {
		select {
		case <-time.After(time.Duration(1) * time.Millisecond * 1000):
			err := conn.WriteMessage(websocket.TextMessage, []byte("/ping"))
			if err != nil {
				log.Println("Error during writing to websocket:", err)
				return
			}
		case <-interrupt:
			log.Println("Closing all pending connections due to SIGINT signal")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error during closing websocket:", err)
				return
			}
			select {
			case <-done:
				log.Println("Receiver Channel Closed! Exiting....")
			case <-time.After(time.Duration(1) * time.Second):
				log.Println("Timeout in closing receiving channel. Exiting....")
			}
			return
		}
	}
}
