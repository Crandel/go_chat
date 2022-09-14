package chatting

import (
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

type Client struct {
	Nick     *string
	conn     *websocket.Conn
	commands chan<- Command
}

func (u *Client) ReadCommands() {
	log.SetPrefix("chatting#user#ReadCommands ")
	defer u.conn.Close()
	for {
		_, p, err := u.conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			u.commands <- Command{
				id:     CmdQuit,
				client: u,
			}
			return
		}
		rawCommand := string(p)
		args := strings.Split(rawCommand, " ")
		cmd := strings.TrimSpace(args[0])
		log.Println("Command: " + cmd)
		var cmdID CommandID
		if !strings.HasPrefix(cmd, "/") {
			cmdID = CmdMsg
		} else {
			switch cmd {
			case CmdPing:
				cmdID = CmdPing
			case CmdJoin:
				cmdID = CmdJoin
			case CmdUsers:
				cmdID = CmdUsers
			case CmdRooms:
				cmdID = CmdRooms
			case CmdQuit:
				cmdID = CmdQuit
			default:
				u.WriteMsg("ERR: Unknown command " + cmd)
				continue
			}
		}
		u.commands <- Command{
			id:     cmdID,
			client: u,
			args:   args,
		}
	}
}

func (u *Client) WriteMsg(message string) {
	messageType := websocket.TextMessage
	if err := u.conn.WriteMessage(messageType, []byte(message)); err != nil {
		log.Println("WriteMsg " + err.Error())
		return
	}
}
