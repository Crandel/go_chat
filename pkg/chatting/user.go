package chatting

import (
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

type User struct {
	Nick     string
	conn     *websocket.Conn
	commands chan<- Command
}

func (u *User) ReadCommands() {
	defer u.conn.Close()
	for {
		_, p, err := u.conn.ReadMessage()
		if err != nil {
			log.Println("chatting#user#ReadCommands " + err.Error())
			u.commands <- Command{
				id:   CMD_QUIT,
				user: u,
			}
			return
		}
		raw_command := string(p)
		args := strings.Split(raw_command, " ")
		cmd := strings.TrimSpace(args[0])
		fmt.Println("chatting#user#ReadCommands Command: " + cmd)
		var cmdId CommandID
		if !strings.HasPrefix(cmd, "/") {
			cmdId = CMD_MSG
		} else {
			switch cmd {
			case CMD_PING:
				cmdId = CMD_PING
			case CMD_JOIN:
				cmdId = CMD_JOIN
			case CMD_USERS:
				cmdId = CMD_USERS
			case CMD_ROOMS:
				cmdId = CMD_ROOMS
			case CMD_QUIT:
				cmdId = CMD_QUIT
			default:
				u.WriteMsg("ERR: Unknown command " + cmd)
				continue
			}
		}
		u.commands <- Command{
			id:   cmdId,
			user: u,
			args: args,
		}
	}
}

func (u *User) WriteMsg(message string) {
	messageType := websocket.TextMessage
	if err := u.conn.WriteMessage(messageType, []byte(message)); err != nil {
		log.Println("chatting#user#WriteMsg " + err.Error())
		return
	}
}
