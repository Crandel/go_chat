package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var input chan string
var chat chan string
var done chan interface{}
var interrupt chan os.Signal

func msgHandler(conn *websocket.Conn) {
	defer close(done)
	defer close(input)
	for {
		select {
		case <-interrupt:
			return
		default:
			i := bufio.NewScanner(os.Stdin)
			i.Scan()
			input <- i.Text()
		}
	}
}

func reader(conn *websocket.Conn) {
	defer close(chat)
	for {
		select {
		case <-interrupt:
			return
		case <-time.After(1 * time.Second):
			_, p, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("Err: %s\n", err.Error())
				close(done)
				return
			}
			chat <- string(p)
		}
	}
}

func main() {
	input = make(chan string)
	chat = make(chan string)
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
	go reader(conn)
	// Join test room
	err = conn.WriteMessage(websocket.TextMessage, []byte("/join test"))
	if err != nil {
		log.Println("Error during writing to websocket:", err)
		return
	}
	for {
		select {
		case m := <-chat:
			fmt.Printf("> %s\n", m)
		case i := <-input:
			err := conn.WriteMessage(websocket.TextMessage, []byte(i))
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
			case <-time.After(1 * time.Second):
				log.Println("Timeout in closing receiving channel. Exiting....")
			}
			return
		}
	}
}
