package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var input chan string
var chat chan string
var done chan interface{}
var interrupt chan os.Signal

func msgHandler(conn *websocket.Conn, rdr bufio.Reader) {
	defer close(input)
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			return
		default:
			line, err := rdr.ReadString('\n')
			if err != nil {
				log.Println("Could not scan the message")
				close(done)
				return
			}
			input <- strings.Trim(line, "\n")
		}
	}
}

func reader(conn *websocket.Conn) {
	defer close(chat)
	for {
		select {
		case <-interrupt:
			close(done)
			return
		case <-done:
			return
		case <-time.After(1 * time.Second):
			_, p, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Err: %s\n", err.Error())
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
	rdr := bufio.NewReader(os.Stdin)

	log.Println("Please provide room name:")
	roomName, err := rdr.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Please provide user name:")
	userName, err := rdr.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT
	socketURL := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(socketURL, nil)
	if err != nil {
		log.Fatal("Could not connect to WebSocker server '"+socketURL+"'.", err)
	}
	defer conn.Close()

	// Join test room
	// TODO: replace this with correct authentication
	err = conn.WriteMessage(websocket.TextMessage,
		[]byte(
			fmt.Sprintf("/join %s %s",
				strings.Trim(roomName, "\n"),
				strings.Trim(userName, "\n"))))

	go msgHandler(conn, *rdr)
	go reader(conn)
	if err != nil {
		log.Println("Error during writing to websocket:", err)
		return
	}
	log.Printf("You are in room %s", roomName)
	for {
		select {
		case m := <-chat:
			log.Println("> ", m)
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
