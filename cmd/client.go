package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Crandel/go_chat/internal/auth"
	ch "github.com/Crandel/go_chat/internal/chatting"
	lg "github.com/Crandel/go_chat/internal/logging"
	"github.com/gorilla/websocket"
	"gitlab.com/greyxor/slogor"
)

type CommandID string

const (
	CmdJoin  CommandID = "/join"
	CmdPing  CommandID = "/ping"
	CmdQuit  CommandID = "/quit"
	CmdRooms CommandID = "/rooms"
	CmdUsers CommandID = "/users"
)

const host = "localhost:8080"
const apiHost = host + "/api"

var input chan ch.ChatMessage
var chat chan ch.ChatMessage
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
				slog.Debug("Could not scan the message")
				close(done)
				return
			}

			line := strings.Trim(rawLine, "\n")
			args := strings.Split(line, " ")
			cmd := strings.TrimSpace(args[0])
			comId, err := convertToChatCommandID(cmd)
			var message ch.ChatMessage
			slog.Debug("args before error", slog.String("attrs", strings.Join(args, " ")))
			if err != nil {
				comId = ch.CmdMsg
			} else {
				args = args[1:]
			}
			slog.Debug("args after error", slog.String("attrs", strings.Join(args, " ")))
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
				slog.Warn("err", slogor.Err(err))
				close(done)
				return
			}
			chat <- message
		}
	}
}

func main() {
	input = make(chan ch.ChatMessage)
	chat = make(chan ch.ChatMessage)
	done = make(chan interface{})
	interrupt = make(chan os.Signal)

	var newUser bool
	show := os.Getenv("SHOW")
	debug := os.Getenv("DEBUG")
	logLevel := slog.LevelInfo
	if debug != "" {
		logLevel = slog.LevelDebug
	}
	showSourse := show != ""
	lg.InitLogger(logLevel, showSourse)

	intLog := slog.NewLogLogger(slog.NewTextHandler(os.Stdout, nil), logLevel)

	if len(os.Args) > 1 {
		newUser = true
	}
	rdr := bufio.NewReader(os.Stdin)

	fmt.Println("Please provide user name:")
	userName, err := rdr.ReadString('\n')
	if err != nil {
		intLog.Fatal(err)
	}
	userName = strings.Trim(userName, "\n")

	fmt.Println("Please provide password:")
	password, err := rdr.ReadString('\n')
	if err != nil {
		intLog.Fatal("Error after password", err)
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
		intLog.Fatal("Error in Post ", err)
	}
	slog.Debug(resp.Request.URL.String())

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		intLog.Fatal(err)
		return
	}
	var auth auth.Response
	err = json.Unmarshal(body, &auth)
	if err != nil {
		intLog.Fatal("Response Body ", err)
		return
	}

	fmt.Println("Please provide room name:")
	roomName, err := rdr.ReadString('\n')
	if err != nil {
		intLog.Fatal(err)
	}
	roomName = strings.Trim(roomName, "\n")

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT
	socketURL := "ws://" + apiHost + "/ws"
	token := "Basic " + auth.Token
	conn, _, err := websocket.DefaultDialer.Dial(socketURL, http.Header{"Authorization": []string{token}})
	if err != nil {
		intLog.Fatal("Could not connect to WebSocker server '"+socketURL+"'.", err)
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
		slog.Warn("Error during writing to websocket:", err)
		return
	}
	fmt.Printf("You are in room '%s'\n", roomName)
	for {
		select {
		case <-done:
			return
		case m := <-chat:
			var message string
			if m.User != nil {
				message = "[" + *m.User + "]-> "
			}
			if len(m.Args) > 0 {
				message = message + strings.Join(m.Args, " ")
			}
			fmt.Println("# ", message)
		case i := <-input:
			err := conn.WriteJSON(&i)
			if err != nil {
				slog.Warn("Error during writing to websocket:", err)
				return
			}
		case <-interrupt:
			slog.Warn("Closing all pending connections due to SIGINT signal")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				slog.Warn("Error during closing websocket:", err)
				return
			}
			select {
			case <-done:
				slog.Warn("Receiver Channel Closed! Exiting....")
			case <-time.After(1 * time.Second):
				slog.Warn("Timeout in closing receiving channel. Exiting....")
			}
			return
		}
	}
}
