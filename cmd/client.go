package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Crandel/go_chat/internal/auth"
	ch "github.com/Crandel/go_chat/internal/chatting"
	lg "github.com/Crandel/go_chat/internal/logging"
	"github.com/gorilla/websocket"
)

type CommandID string

const (
	CmdJoin  CommandID = "/join"
	CmdPing  CommandID = "/ping"
	CmdQuit  CommandID = "/quit"
	CmdRooms CommandID = "/rooms"
	CmdUsers CommandID = "/users"
)

var log = lg.InitLogger()

const host = "localhost:8080"
const apiHost = host + "/api"

var input chan ch.ChatMessage
var chat chan string
var done chan interface{}
var interrupt chan os.Signal

func convertToChatCommandID(input string) (ch.CommandID, error) {
	commandID := CommandID(input)
	switch commandID {
	case CmdJoin, CmdPing, CmdQuit, CmdRooms, CmdUsers:
		chatInput := strings.TrimPrefix(input, "/")
		chatCommandID, err := ch.ConvertToCommandID(chatInput)
		if err != nil {
			return "", err
		}
		return chatCommandID, nil
	default:
		return "", fmt.Errorf("invalid CommandID: %s", input)
	}
}

func msgHandler(conn *websocket.Conn, rdr bufio.Reader) {
	defer close(input)
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			return
		default:
			rawLine, err := rdr.ReadString('\n')
			if err != nil {
				log.Log(lg.Debug, "Could not scan the message")
				close(done)
				return
			}

			line := strings.Trim(rawLine, "\n")
			args := strings.Split(line, " ")
			cmd := strings.TrimSpace(args[0])
			comId, err := convertToChatCommandID(cmd)
			var message ch.ChatMessage
			if err != nil {
				comId = ch.CmdMsg
			} else {
				args = args[1:]
			}
			message.CommandId = comId
			message.Args = args
			input <- message
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
			var message ch.ChatMessage
			err := conn.ReadJSON(&message)
			if err != nil {
				log.Logf(lg.Warning, "err: %s", err.Error())
				close(done)
				return
			}
			if len(message.Args) > 0 {
				chat <- string(message.Args[0])
			}
		}
	}
}

func main() {
	input = make(chan ch.ChatMessage)
	chat = make(chan string)
	done = make(chan interface{})
	interrupt = make(chan os.Signal)

	var newUser bool
	debug := os.Getenv("DEBUG")
	log.PrintDebug = debug == "1"

	if len(os.Args) > 1 {
		newUser = true
	}
	rdr := bufio.NewReader(os.Stdin)

	log.Log(lg.NoLogging, "Please provide user name:")
	userName, err := rdr.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	userName = strings.Trim(userName, "\n")

	log.Log(lg.NoLogging, "Please provide password:")
	password, err := rdr.ReadString('\n')
	if err != nil {
		log.Fatal("Error after password", err)
	}
	password = strings.Trim(password, "\n")
	postBody, _ := json.Marshal(map[string]string{
		"nick":     userName,
		"password": password,
	})
	responseBody := bytes.NewBuffer(postBody)
	loginUrl := "login"
	if newUser {
		loginUrl = "signin"
	}
	resp, err := http.Post("http://"+apiHost+"/auth/"+loginUrl, "application/json", responseBody)

	if err != nil {
		log.Fatal("Error in Post ", err)
	}
	log.Log(lg.Debug, resp.Request.URL)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	var auth auth.Response
	err = json.Unmarshal(body, &auth)
	if err != nil {
		log.Fatal("Response Body ", err)
		return
	}

	log.Log(lg.NoLogging, "Please provide room name:")
	roomName, err := rdr.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	roomName = strings.Trim(roomName, "\n")

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT
	socketURL := "ws://" + apiHost + "/ws"
	token := "Basic " + auth.Token
	conn, _, err := websocket.DefaultDialer.Dial(socketURL, http.Header{"Authorization": []string{token}})
	if err != nil {
		log.Fatal("Could not connect to WebSocker server '"+socketURL+"'.", err)
	}
	defer conn.Close()

	// Join test room
	joinMsg := ch.ChatMessage{
		CommandId: ch.CmdJoin,
		Args: []string{
			roomName,
		},
	}
	err = conn.WriteJSON(&joinMsg)

	go msgHandler(conn, *rdr)
	go reader(conn)
	if err != nil {
		log.Log(lg.Warning, "Error during writing to websocket:", err)
		return
	}
	log.Logf(lg.NoLogging, "You are in room '%s'\n", roomName)
	for {
		select {
		case <-done:
			return
		case m := <-chat:
			log.Log(lg.NoLogging, "# ", m)
		case i := <-input:
			err := conn.WriteJSON(i)
			if err != nil {
				log.Log(lg.Warning, "Error during writing to websocket:", err)
				return
			}
		case <-interrupt:
			log.Log(lg.Warning, "Closing all pending connections due to SIGINT signal")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Log(lg.Warning, "Error during closing websocket:", err)
				return
			}
			select {
			case <-done:
				log.Log(lg.Warning, "Receiver Channel Closed! Exiting....")
			case <-time.After(1 * time.Second):
				log.Log(lg.Warning, "Timeout in closing receiving channel. Exiting....")
			}
			return
		}
	}
}
